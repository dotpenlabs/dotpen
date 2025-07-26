<script lang="ts">
	import { onMount, tick } from 'svelte';
	import { fade } from 'svelte/transition';

	import { pb } from '$/lib';
	import { page } from '$app/state';
	import { toast } from 'svelte-sonner';
	import { getItem as idbGetItem, setItem as idbSetItem } from '$/lib/idb';

	import type { PageData } from './$types';
	import type { RecordModel } from 'pocketbase';
	import type { Bookmark } from '$/lib/types';

	import Link from '$/lib/components/item/link.svelte';
	import Loading from '$/lib/components/loading.svelte';
	import { browser } from '$app/environment';

	let { data: pageData }: { data: PageData } = $props();
	let url = $state('');
	let colId = $state('');

	let Macy: any;
	let masonry: any | undefined;

	let global = $state('show') as 'load' | 'show' | 'error';
	let bookmarks: Bookmark[] = $state([]);
	let token = $state('');

	async function main(): Promise<void> {
		if (url === page.url.pathname) {
			console.info('[main] url has not changed, aborting');
			return;
		}

		url = page.url.pathname;

		console.info(`[main] started process (${page.url.pathname})`);
		colId = pageData.collection;

		if (colId === 'inbox') {
			console.info('[head] This is a system inbox, additional processing is required.');
			const inboxId = await idbGetItem('inbox:id');
			if (inboxId) {
				colId = inboxId;
			} else {
				try {
					console.info('[head] Trying to get inbox from server...');
					colId = await pb
						.collection('collections')
						.getFirstListItem("name = 'system_inbox'")
						.then((e: RecordModel) => {
							return e.id;
						});

					await idbSetItem('inbox:id', colId);
				} catch (e) {
					console.info('[head] No inbox found, creating one...');
					toast.info('Welcome to Dotpen!', {
						description:
							"Thanks for joining us at Dotpen, we're excited to give you the best experience possible."
					});

					colId = await pb
						.collection('collections')
						.create({
							name: 'system_inbox',
							user: pb.authStore.record.id
						})
						.then((e: RecordModel) => {
							return e.id;
						});

					await idbSetItem('inbox:id', colId);
				}
			}
		}

		let bmcache = (await JSON.parse((await idbGetItem(colId + ':cache')) || '[]')) as Bookmark[];
		const bmheartbeat = bmcache[bmcache.length - 1]?.updated;

		console.info('[head] connected to ' + colId);
		console.info('[head] last heartbeat was sent ' + (bmheartbeat ? 'at ' + bmheartbeat : 'never'));
		console.info('[head] bookmark cache is ' + (bmcache ? 'available' : 'unavailable'));

		const getBase64 = async (bm: Bookmark, path: string): Promise<string | null> => {
			if (!path) return null;
			try {
				const r = await fetch(
					`${pb.baseURL}api/files/bookmarks/${bm.id}/${path}?token=${pb.authStore.token}`
				);
				const blob = await r.blob();
				return await new Promise<string>((resolve, reject) => {
					const reader = new FileReader();
					reader.onloadend = () => resolve(reader.result as string);
					reader.onerror = reject;
					reader.readAsDataURL(blob);
				});
			} catch (e) {
				console.warn('[bookmark:sync] Failed to convert image to base64', path, e);
				return null;
			}
		};

		// Zorg dat ALLE bookmarks in de cache base64-afbeeldingen krijgen
		const ensureLocalImages = async (data: Bookmark[]): Promise<Bookmark[]> => {
			await Promise.all(
				data.map(async (bm) => {
					if (bm.favicon && !bm._favicon_base64) {
						bm._favicon_base64 = await getBase64(bm, bm.favicon);
					}
					if (bm.cover && !bm._cover_base64) {
						bm._cover_base64 = await getBase64(bm, bm.cover);
					}
				})
			);
			return data;
		};

		if (!bmcache || bmcache.length === 0) {
			window.SetHydrating('collection', true);
			console.info('[head:download] No cache available, downloading...');
			await idbSetItem(colId + ':cache', JSON.stringify([]));

			const remote = (await pb.collection('bookmarks').getFullList({
				filter: `collection = "${colId}" && deleted = false`,
				sort: '-created'
			})) as Bookmark[];

			await ensureLocalImages(remote);
			console.info('[head:bookmarks] Downloaded ' + remote.length + ' bookmarks');

			await idbSetItem(colId + ':cache', JSON.stringify(remote));
			bmcache = remote;
			console.info('[head:download] Downloaded ' + remote.length + ' bookmarks');
			window.SetHydrating('collection', false);
		} else {
			(async () => {
				window.SetHydrating('collection', true);
				console.info('[head:bookmarks] Downloading latest bookmarks');

				const local = await ensureLocalImages(bmcache);
				const remote = (await pb.collection('bookmarks').getFullList({
					filter: `collection = "${colId}" && updated > '${local.at(-1)?.updated || new Date(0).toISOString()}'`,
					sort: '-created'
				})) as Bookmark[];

				console.info('[head:bookmarks] Downloaded ' + remote.length + ' bookmarks');

				await ensureLocalImages(remote);
				const ubi = new Map([...local, ...remote].map((bm) => [bm.id, bm]));
				bmcache = Array.from(ubi.values()).filter((bm) => !bm.deleted);
				await idbSetItem(colId + ':cache', JSON.stringify(bmcache));

				console.info('[head:bookmarks] Updated cache has ' + bmcache.length + ' bookmarks');
				window.SetHydrating('collection', false);
			})();
		}

		// - finishing up -
		bookmarks = bmcache;
		token = await pb.files.getToken();

		masonry?.recalculate();

		window.SetHydrating('collection', false);
	}

	async function UseMasonry(): Promise<void> {
		await tick();

		if (!Macy && browser) {
			const macyModule = await import('macy');
			Macy = macyModule.default;
		}

		const container = document.querySelector('#main-content');
		if (!container) return;

		if (masonry) {
			try {
				masonry.remove();
			} catch (e) {
				console.warn('Error removing masonry', e);
			}
		}

		try {
			masonry = new Macy({
				container,
				trueOrder: true,
				waitForImages: false,
				useOwnImageLoader: false,
				useContainerForBreakpoints: true,
				mobileFirst: true,
				columns: 1,
				margin: { y: 16, x: '2%' },
				breakAt: { 1400: 5, 1100: 4, 800: 3, 520: 2 }
			});

			requestAnimationFrame(() => masonry.recalculate());
		} catch (e) {
			console.error('Error initializing masonry', e);
		}
	}

	async function newBookmark(url: string): Promise<Error | void> {
		// - Helpers & Functions -
		const urlF = (input: string): string | null => {
			if (/\s/.test(input)) return null;

			try {
				let uri = new URL(input);
				if (!uri.protocol.startsWith('http')) {
					uri = new URL('https://' + input);
				}
				uri.protocol = 'https:';

				const hostnameParts = uri.hostname.split('.');
				const tld = hostnameParts[hostnameParts.length - 1];
				if (!tld || tld.length < 2) return null;

				return uri.toString();
			} catch {
				try {
					const uri = new URL('https://' + input);

					const hostnameParts = uri.hostname.split('.');
					const tld = hostnameParts[hostnameParts.length - 1];
					if (!tld || tld.length < 2) return null;

					return uri.toString();
				} catch {
					return null;
				}
			}
		};

		const fileF = async (
			url: string,
			filename: string
		): Promise<{ file: File | null; base64: string | null }> => {
			try {
				const r = await fetch(`${pb.baseURL}api/proxy?url=${encodeURIComponent(urlF(url))}`, {
					headers: { Authorization: `Bearer ${pb.authStore.token}` }
				});
				if (!r.ok) {
					console.warn('[bookmark:proxy] 1 asset download did not succeed', r.status);
					return { file: null, base64: null };
				}

				const blob = await r.blob();
				const file = new File([blob], filename, {
					type: blob.type || `image/${filename.split('.').pop() || 'png'}`
				});

				const base64 = await new Promise<string>((resolve, reject) => {
					const reader = new FileReader();
					reader.onloadend = () => resolve(reader.result as string);
					reader.onerror = reject;
					reader.readAsDataURL(blob);
				});

				return { file, base64 };
			} catch (e) {
				console.info('[bookmark:error] 1 asset download did not succeed', e);
				return { file: null, base64: null };
			}
		};

		// - Validation and setup -

		url = urlF(url);

		if (!url) {
			console.info('[bookmark:failure] Tried to add non-existing URL:', url);
			toast.error('Heads up!', {
				description: "That doesn't look like a valid URL"
			});
			return;
		}

		const loading = toast.loading('Working on it!', {
			description: 'This could take a few seconds...'
		});

		// - Crawling the link -

		const crawl = (await pb.send('/api/crawl', {
			method: 'GET',
			query: { url }
		})) as {
			title: string;
			favicon?: string;
			image?: string;
			url: string;
		};

		console.info('[bookmark:crawl] Result', crawl);

		// - Capturing metadata and images -

		['favicon', 'image'].forEach(async (key) => {
			if (crawl[key].startsWith('./')) {
				crawl[key] = crawl[key].replace('./', '/');
			}
			if (crawl[key].startsWith('/')) {
				crawl[key] = new URL(crawl.url).origin + crawl[key];
			}
		});

		const { file: crawl_favicon, base64: base64_favicon } = await fileF(
			crawl.favicon,
			'favicon.png'
		);
		const { file: crawl_image, base64: base64_image } = await fileF(crawl.image, 'cover.png');

		console.info('[bookmark:crawl] Assets downloaded');

		// - Storing data -

		const form: Record<string, any> = {
			label: crawl.title ?? new URL(url).hostname,
			link: url,
			collection: colId,
			user: pb.authStore.record.id,
			created: new Date().toISOString()
		};

		if (crawl_favicon) form.favicon = crawl_favicon;
		if (crawl_image) form.cover = crawl_image;

		const record = (await pb.collection('bookmarks').create(form)) as Bookmark;

		if (base64_favicon) record._favicon_base64 = base64_favicon;
		if (base64_image) record._cover_base64 = base64_image;

		// - Finishing up -
		bookmarks.unshift(record);
		await idbSetItem(colId + ':cache', JSON.stringify(bookmarks));
		toast.dismiss(loading);

		window.SetColored('green');
	}

	async function removeBookmark(item: Bookmark) {
		const bm = bookmarks.find((m) => m.id === item.id);
		if (!bm) return;

		bookmarks = bookmarks.filter((m) => m.id !== item.id);
		await idbSetItem(colId + ':cache', JSON.stringify(bookmarks));

		await tick();

		requestAnimationFrame(() => {
			masonry?.recalculate(true);
		});

		let confirmed = true;

		const confirmDelete = async () => {
			if (!confirmed) return;
			try {
				if (bm.id) {
					await pb.collection('bookmarks').update(bm.id, {
						deleted: true
					});
					await idbSetItem(colId + ':cache', JSON.stringify(bookmarks));
					window.SetColored('red');
				} else {
					toast.error('Failed to delete bookmark! (local error)');
				}
			} catch (err) {
				console.error('Failed to delete bookmark', err);
				toast.error('Failed to delete bookmark!');
			}
		};

		toast.info('Link is deleted!', {
			description: 'The link will deleted when this dialog is dismissed.',
			action: {
				label: 'Undo',
				onClick: async () => {
					confirmed = false;
					try {
						bookmarks = [bm, ...bookmarks];
						await idbSetItem(colId + ':cache', JSON.stringify(bookmarks));
						await tick();
						requestAnimationFrame(async () => {
							await masonry?.recalculate(true);
							await main();
						});
					} catch (err) {
						console.error('Failed to undo delete', err);
						toast.error('Failed to restore bookmark!');
					}
				}
			},
			onAutoClose: () => confirmDelete(),
			onDismiss: () => confirmDelete()
		});
	}

	onMount(async () => {
		await main();
		if (bookmarks.length > 0) {
			UseMasonry();
		}
	});
	$effect(() => void main());
	$effect(() => {
		if (bookmarks.length > 0) {
			UseMasonry();
		}
	});
</script>

<svelte:window onpaste={(e) => newBookmark(e.clipboardData?.getData('text/plain') || '')} />

{#if global === 'load'}
	<content class="h-full w-full flex justify-center items-center">
		<Loading />
	</content>
{:else if global === 'error'}
	<content class="h-full w-full flex justify-center items-center">
		<p>Something went wrong, please try again later.</p>
	</content>
{:else}
	<content
		class="h-full w-full overflow-y-auto flex flex-col gap-2 justify-start items-start p-1 pr-6"
		role="presentation"
		ondragover={(e: DragEvent) => {
			e.preventDefault();
		}}
		ondrop={(e: DragEvent) => {
			e.preventDefault();
			newBookmark(e.dataTransfer?.getData('text/plain'));
		}}
	>
		{#if bookmarks.length === 0}
			<div
				in:fade={{ duration: 500, delay: 750 }}
				class="flex flex-col h-full w-full gap-2 justify-center items-center opacity-65"
			>
				<img src="/logo.svg" alt="Dotpen" class="size-6 invert dark:invert-0" />
				<div class="text-center">
					<p class="text-lg font-medium">This collection is tidy!</p>
					<p class="text-xs w-64">
						You don't have any bookmarks saved here, drag and drop some links to get it to fill!
					</p>
				</div>
			</div>
		{:else}
			<div id="main-content" class="w-full min-h-full h-fit">
				{#each bookmarks as item, idx (item.id)}
					<Link
						removeItemCallback={(e: Event) => {
							requestAnimationFrame(() => {
								masonry?.recalculate(true);
							});
						}}
						removeItem={() => removeBookmark(item)}
						data={item}
						index={idx}
						filetoken={token}
					/>
				{/each}
			</div>
		{/if}
	</content>
{/if}
