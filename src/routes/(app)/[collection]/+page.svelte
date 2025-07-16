<script lang="ts">
	import { fade } from 'svelte/transition';
	import LinkComponent from '$/lib/components/item/link.svelte';
	import type { LinkItem } from '$/lib/components/item/types';
	import { toast } from 'svelte-sonner';
	import { onDestroy, onMount, tick } from 'svelte';
	import type { PageData } from './$types';
	import { pb } from '$/lib';
	import { page } from '$app/state';

	let Marks: LinkItem[] = $state([]);
	let { data: pageData }: { data: PageData } = $props();
	let pasteTries = 0;

	$effect(() => {
		(async () => {
			if (pageData.collection) {
				if (pageData.collection == 'inbox') {
					try {
						const resp = await pb
							.collection('collections')
							.getFirstListItem("name = 'system_inbox'");
						pageData.collection = resp.id;
					} catch {
						const resp = await pb.collection('collections').create({
							name: 'system_inbox',
							user: pb.authStore.record.id
						});

						pageData.collection = resp.id;
					}
				}
				await fetchBookmarks();
			}
		})();
	});

	let Macy: any;
	let masonry: any | undefined;

	async function urlToFile(url: string, filename: string): Promise<File | null> {
		try {
			const res = await fetch(url);
			if (!res.ok) return null;
			const blob = await res.blob();
			const ext = filename.split('.').pop() || 'png';
			return new File([blob], filename, { type: blob.type || `image/${ext}` });
		} catch {
			return null;
		}
	}

	import type { RecordModel } from 'pocketbase';

	interface LinkRecordItem extends RecordModel {
		label: string;
		link: string;
		favicon: string;
		cover: string;
	}

	async function fetchBookmarks() {
		const result = await pb.collection('bookmarks').getFullList({
			filter: `collection = "${pageData.collection}"`,
			sort: '-created'
		});
		const fileToken = await pb.files.getToken();
		Marks = result.map((bm: LinkRecordItem) => ({
			label: bm.label,
			icon: pb.files.getURL(bm, bm.favicon, {
				thumb: '50x50',
				token: fileToken
			}),
			url: bm.link,
			cover: pb.files.getURL(bm, bm.cover, {
				thumb: '0x200',
				token: fileToken
			}),
			date: bm.created,
			id: bm.id
		}));
	}

	onMount(async () => {
		Macy = (await import('macy')).default as any;
		await fetchBookmarks();
	});

	async function initMasonry() {
		await tick();
		const container = document.querySelector('#main-content');
		if (!container) return;

		if (masonry) {
			masonry.remove();
		}

		var masonry = new Macy({
			container: '#main-content',
			trueOrder: true,
			waitForImages: false,
			useOwnImageLoader: false,
			useContainerForBreakpoints: true,
			mobileFirst: true,
			columns: 1,
			margin: {
				y: 16,
				x: '2%'
			},
			breakAt: {
				1400: 5,
				1100: 4,
				800: 3,
				520: 2
			}
		});

		setInterval(() => {
			masonry?.recalculate();
		}, 1000);
	}

	$effect(() => {
		if (Marks.length > 0) {
			initMasonry();
		}
	});

	async function fetchImageFileThroughProxy(url: string, filename: string): Promise<File | null> {
		try {
			const res = await fetch(pb.baseURL + `api/proxy?url=${encodeURIComponent(url)}`, {
				headers: {
					Authorization: `Bearer ${pb.authStore.token}`
				}
			});
			if (!res.ok) return null;
			const blob = await res.blob();
			const ext = filename.split('.').pop() || 'png';
			return new File([blob], filename, { type: blob.type || `image/${ext}` });
		} catch (e) {
			console.error('Failed to fetch image through proxy', e);
			return null;
		}
	}

	const isValidUrl = (url: string): boolean => {
		const urlRegex =
			/(https?:\/\/(?:www\.|(?!www))[a-zA-Z0-9][a-zA-Z0-9-]+[a-zA-Z0-9]\.[^\s]{2,}|www\.[a-zA-Z0-9][a-zA-Z0-9-]+[a-zA-Z0-9]\.[^\s]{2,}|https?:\/\/(?:www\.|(?!www))[a-zA-Z0-9]+\.[^\s]{2,}|www\.[a-zA-Z0-9]+\.[^\s]{2,})/;
		return urlRegex.test(url);
	};

	const normalizeUrl = (url: string): string => {
		const trimmed = url.trim();
		return trimmed.startsWith('http://') || trimmed.startsWith('https://')
			? trimmed
			: `https://${trimmed}`;
	};

	async function addBookmark(requestUrl: string): Promise<void> {
		if (!requestUrl?.trim()) return;
		const startTime = performance.now();

		const isDuplicate = (url: string): boolean => {
			return Marks.some((mark) => mark.url === url);
		};

		const confirmDuplicate = (): Promise<boolean> => {
			return new Promise((resolve) => {
				toast.warning('Duplicate item found', {
					description: 'This link already exists. Add anyway?',
					action: {
						label: 'Add',
						onClick: () => resolve(true)
					}
				});
			});
		};

		const showWarningToast = (bookmark: LinkItem): void => {
			const hasTitle = bookmark.label !== bookmark.url;

			toast.warning('Link has been captured', {
				description: hasTitle
					? `Added: "${bookmark.label}" but some data might be missing`
					: "We couldn't get the title of the link. It can be harder to find it back.",
				action: {
					label: 'Remove',
					onClick: async () => {
						await pb.collection('bookmarks').delete(bookmark.id);
						await fetchBookmarks();
					}
				}
			});
		};

		const showSuccessToast = (bookmark: LinkItem, startTime: number): void => {
			const hasTitle = Boolean(bookmark.label && bookmark.label !== bookmark.url);

			if (!hasTitle) {
				showWarningToast(bookmark);
			} else {
				const duration = (performance.now() - startTime).toFixed(2);
				toast.success('Link has been captured!', {
					description: `Added: "${bookmark.label}" in ${duration}ms`
				});
			}
		};

		const formattedUrl = normalizeUrl(requestUrl);
		if (!isValidUrl(formattedUrl)) {
			toast.error('Invalid link', {
				description: 'The link you`ve tried to add was invalid.'
			});
			return;
		}

		if (isDuplicate(formattedUrl)) {
			const proceed = await confirmDuplicate();
			if (!proceed) return;
		}

		const loadingToast = toast.loading('Capturing link', {
			description: 'This may take a few seconds...'
		});

		try {
			const response = await pb.send('/api/crawl', {
				query: {
					url: formattedUrl
				}
			});

			const data = await response;

			['favicon', 'image'].forEach((key) => {
				if (typeof data?.[key] === 'string' && data[key].startsWith('/')) {
					data[key] = formattedUrl + data[key];
				}
			});

			const faviconFile = data.favicon
				? await fetchImageFileThroughProxy(data.favicon, 'favicon.png')
				: null;
			const coverFile = data.image
				? await fetchImageFileThroughProxy(data.image, 'cover.png')
				: null;

			const newBookmark: LinkItem = {
				label: data.title || formattedUrl,
				icon: '',
				url: formattedUrl,
				cover: '',
				date: new Date().toISOString()
			};

			window.dispatchEvent(new CustomEvent('network_activity:start'));

			const formData: any = {
				label: newBookmark.label,
				link: newBookmark.url,
				collection: pageData.collection,
				created: new Date().toISOString()
			};
			if (faviconFile) formData.favicon = faviconFile;
			if (coverFile) formData.cover = coverFile;

			const bookmark = await pb.collection('bookmarks').create(formData);
			await fetchBookmarks();
			await tick();
			masonry?.recalculate(true);

			if (!Boolean(newBookmark.label && newBookmark.label !== newBookmark.url)) {
				showWarningToast(bookmark as unknown as LinkItem);
			} else {
				showSuccessToast(bookmark as unknown as LinkItem, startTime);
			}
		} catch (error) {
			console.error('Failed to add bookmark:', error);
			toast.error('We failed to capture this link...', {
				description: 'Our crawlers are not able to access this link.'
			});
		} finally {
			window.dispatchEvent(new CustomEvent('network_activity:stop'));
			toast.dismiss(loadingToast);
		}
	}

	function removeBookmark(item: LinkItem) {
		let confirmed = true;

		const bm = Marks.find((mark) => mark.id === item.id);
		if (!bm) return;

		Marks = Marks.filter((mark) => mark.id !== item.id);
		masonry?.recalculate(true);

		const confirmDelete = () => {
			(async () => {
				if (!confirmed) return;
				try {
					if (bm.id) {
						await pb.collection('bookmarks').delete(bm.id);
						await fetchBookmarks();
					}
				} catch (err) {
					console.error('Failed to delete bookmark', err);
					toast.error('Failed to delete bookmark!');
				}
			})();
		};

		toast.info('Link has been removed', {
			description: 'This link has been removed from the collection.',
			action: {
				label: 'Undo',
				onClick: async () => {
					confirmed = false;
					try {
						Marks = [bm, ...Marks];
						masonry?.recalculate(true);
						await fetchBookmarks();
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

	onDestroy(() => {
		masonry?.remove();
	});
</script>

<svelte:window
	onpaste={(e) => {
		const text = e.clipboardData.getData('text/plain');
		const textnrml = normalizeUrl(text);

		if (text && isValidUrl(textnrml)) {
			addBookmark(textnrml);
		} else {
			if (pasteTries > 2) {
				toast.warning('Pastedata is not valid!', {
					description: "You've tried to paste: " + text
				});
				return;
			}

			console.log('Invalid user action, probably not meant.');
			pasteTries++;

			setTimeout(() => {
				pasteTries = 0;
			}, 10000);
		}
	}}
/>

<content
	class="h-full w-full overflow-y-auto flex flex-col gap-2 justify-start items-start p-1 pr-6"
	role="presentation"
	ondragover={(e: DragEvent) => {
		e.preventDefault();
	}}
	ondrop={(e: DragEvent) => {
		e.preventDefault();
		addBookmark(e.dataTransfer?.getData('text/plain'));
	}}
>
	{#if Marks.length === 0}
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
			{#each Marks as item, idx}
				<LinkComponent removeItem={() => removeBookmark(item)} data={item} index={idx} />
			{/each}
		</div>
	{/if}
</content>
