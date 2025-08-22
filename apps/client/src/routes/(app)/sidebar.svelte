<script lang="ts">
	import { onMount, type Snippet } from 'svelte';
	import { Folder, Gear, MagnifyingGlass, PuzzlePiece, Tray } from 'phosphor-svelte';
	import { fly } from 'svelte/transition';

	import { goto } from '$app/navigation';
	import { pb, kv } from '$/lib';
	import type { Collection } from '$/lib/types';
	import { twMerge } from 'tailwind-merge';
	import Navitem from '$/lib/components/navigation/collection.svelte';
	import { page } from '$app/state';
	import { toast } from 'svelte-sonner';
	import { flyAndScale, focus } from '$/lib/utils';

	import axios from 'axios';

	let intro = $state('');
	const intros = [
		'Hey there, $0',
		'Howdy, $0!',
		'Sup, $0!',
		'Welcome back, $0!',
		'Ready for it, $0?',
		'Look who’s here!',
		'Well well well, if it isn’t $0!',
		'You again? Yay!',
		'Guess who’s back? (It’s you!)',
		'Ahoy, matey!',
		'Live long and prosper, $0 🖖',
		'You’ve got the power, $0!',
		'The force is strong with you.',
		'Beam me up, $0!',
		'To infinity… and beyond!',
		"You're the chosen one, $0!",
		'Here we go again!',
		'Let’s roll, $0!',
		'It’s morphin’ time!',
		'Assemble, $0!',
		'You rang?',
		'Cheers, $0!',
		'Back in action!',
		'Wubba lubba dub dub!',
		'Danger zone... just kidding!',
		'Hey you cool cat',
		'I see you’ve returned, my precious...',
		'Did someone say adventure?',
		'Great Scott, it’s $0!',
		'Feeling the need for speed, $0?',
		'Your mission is ready, $0!',
		'No worries, nothing will self-destruct today, $0!',
		'Your a wizard, $0!',
		'Riddikulus Boo',
		'Stay cool like Captain Man!'
	];

	let {
		collections = $bindable(),
		appstate,
		appplatform,
		sidebarloc
	}: {
		collections: Collection[];
		appstate: string;
		appplatform: string;
		sidebarloc: string;
	} = $props();

	let nCollection = $state({
		name: '',
		open: false
	});

	onMount(() => {
		intro = intros[Math.floor(Math.random() * intros.length)];
		intro = intro.replace('$0', pb.authStore.record.name.split(' ')[0]);

		setTimeout(() => {
			intro = (() => {
				const hours = new Date().getHours();

				if (hours >= 5 && hours < 12) {
					return 'Goodmorning, $0!';
				} else if (hours >= 12 && hours < 18) {
					return 'Goodafternoon, $0!';
				} else if (hours >= 18 && hours < 22) {
					return 'Goodevening, $0!';
				} else {
					return 'Goodnight, $0!';
				}
			})();

			intro = intro.replace('$0', pb.authStore.record.name.split(' ')[0]);
		}, 2500);
	});
</script>

{#snippet divider(name: string)}
	{#if sidebarloc === 'left'}
		<p class="text-xs">{name}</p>
		<div class="flex-1 h-px bg-stone-300 dark:bg-stone-700 rounded-lg"></div>
	{:else}
		<div class="flex-1 h-px bg-stone-300 dark:bg-stone-700 rounded-lg"></div>
		<p class="text-xs">{name}</p>
	{/if}
{/snippet}

<div class="h-full w-full max-w-sm px-8 flex flex-col">
	<div
		class={twMerge(
			'flex gap-1 justify-between items-center pt-8',
			sidebarloc === 'left' ? 'flex-row-reverse' : 'flex-row'
		)}
	>
		<button
			onclick={() => {
				goto('/settings');
			}}
			aria-label="Settings"
			title="Settings"
		>
			<Gear
				class="size-4.5 opacity-25 hover:opacity-100 hover:duration-1200 active:duration-100 active:ease-out hover:ease-in-out cursor-pointer active:opacity-50 hover:rotate-360"
			/>
		</button>
		<div
			class={twMerge(
				'flex flex-col gap-1 w-full',
				sidebarloc === 'left' ? 'text-left' : 'text-right'
			)}
		>
			<div class="max-h-8 min-h-8 h-8 overflow-clip">
				{#key intro}
					<p
						in:fly={{ duration: 500, delay: 750, y: -5 }}
						out:fly={{ duration: 750, y: 5 }}
						class="text-xl font-medium"
					>
						{intro}
					</p>
				{/key}
			</div>
			<div class="max-h-4 min-h-4 h-4 overflow-clip">
				{#if appstate === 'app'}
					<p
						in:fly={{ duration: 500, y: -5 }}
						out:fly={{ duration: 750, y: 5 }}
						class="text-xs opacity-65"
					>
						You're all set! All changes are synced.
					</p>
				{:else if appstate === 'hydrating'}
					<p
						in:fly={{ duration: 500, delay: 750, y: -5 }}
						out:fly={{ duration: 750, y: 5 }}
						class="text-xs opacity-65"
					>
						Hang tight! Almost ready...
					</p>
				{:else}
					<p
						in:fly={{ duration: 500, delay: 750, y: -5 }}
						out:fly={{ duration: 750, y: 5 }}
						class="text-xs opacity-65"
					>
						You are not connected, data might be outdated!
					</p>
				{/if}
			</div>
		</div>
	</div>
	<div class="flex items-center gap-3 w-full opacity-50 mt-5">
		{@render divider('Menu')}
	</div>
	<div class="flex h-8 items-center gap-4 my-4 justify-between w-full">
		<button
			class="flex gap-2 text-[14.5px] justify-center rounded-xl cursor-pointer active:scale-98 active:opacity-90 duration-50 items-center border h-full w-full border-[#E4E4E4] dark:border-[#282828]"
		>
			<MagnifyingGlass />Search
		</button>
		<button
			onclick={() => {
				goto('/plugins');
			}}
			class="flex gap-2 text-[14.5px] justify-center rounded-xl cursor-pointer active:scale-98 active:opacity-90 duration-50 items-center border h-full w-full border-[#E4E4E4] dark:border-[#282828]"
		>
			<PuzzlePiece />Plugins
		</button>
	</div>
	<Navitem
		id="inbox"
		icon={Tray}
		label="Inbox"
		onSelect={() => {}}
		removeCollection={() => {
			toast.warning("You can't delete your inbox!", {
				description: 'What would you do without an inbox? 😨'
			});
		}}
		shareSnapshot={async () => {
			try {
				let inboxId = '';
				try {
					inboxId = (await kv.get('inbox:id')) || '';
				} catch {}
				if (!inboxId) {
					try {
						inboxId = await pb
							.collection('collections')
							.getFirstListItem("name = 'system_inbox'")
							.then((e) => e.id as string);
					} catch {}
				}

				if (!inboxId) {
					toast.error('Kon Inbox niet vinden om te delen.');
					return;
				}

				let data: any[] = [];
				try {
					data = JSON.parse((await kv.get(inboxId + ':cache')) || '[]') || [];
				} catch {}

				const items = (data || []).map((b: any) => ({
					label: b.label,
					link: b.link,
					created: b.created,
					_favicon_base64: b._favicon_base64,
					_cover_base64: b._cover_base64
				}));

				if (!items.length) {
					toast.error('Niets om te delen in deze collectie.');
					return;
				}

				const toastId = toast.loading('Starting upload...', {
					description: 'Please wait, connecting...',
					duration: Infinity
				});

				const snapshot = await pb.collection('snapshots').create(
					{
						name: 'Inbox',
						user: pb.authStore.record?.id,
						collection: inboxId,
						data: items
					},
					{
						fetch: async (url, config) => {
							const axiosConfig: any = {
								method: config?.method || 'GET',
								url,
								headers: config?.headers,
								data: config?.body,
								onUploadProgress: (progressEvent: ProgressEvent) => {
									if (progressEvent.lengthComputable) {
										const percent = (progressEvent.loaded / progressEvent.total) * 100;
										console.log(`Upload: ${percent.toFixed(2)}%`);
										toast.loading(`Uploading data... ${percent.toFixed(2)}%`, {
											description: 'This could take a while, please wait...',
											duration: Infinity,
											id: toastId
										});
									}
								},
								responseType: 'text'
							};

							try {
								const response = await axios(axiosConfig);
								return new Response(response.data, {
									status: response.status,
									statusText: response.statusText,
									headers: new Headers(response.headers)
								});
							} catch (error: any) {
								throw new TypeError('Network request failed');
							}
						}
					}
				);

				toast.dismiss(toastId);
				const shareUrl = `${location.origin}/s/${snapshot.slug || snapshot.id}`;
				await navigator.clipboard.writeText(shareUrl);
				toast.success('Snapshot link copied to clipboard!', { description: shareUrl });
			} catch (e) {
				console.error('Failed to create snapshot (inbox)', e);
				toast.error('Failed to share snapshot');
			}
		}}
	/>
	<div class="flex items-center gap-3 w-full opacity-50 mt-5 mb-1">
		{@render divider('Collections')}
		<button
			class="cursor-pointer"
			onclick={() => {
				nCollection.open = true;
			}}
			>+
		</button>
	</div>

	<div class="flex-1 overflow-y-auto mb-8">
		<div class="flex-1 overflow-y-auto flex flex-col gap-2 min-h-0">
			{#each collections as collection}
				<Navitem
					id={collection.id}
					icon={/^\p{Emoji}/u.test(collection.name) ? null : Folder}
					label={collection.name}
					onSelect={() => {
						goto(`/${collection.id}`);
					}}
					removeCollection={async () => {
						try {
							await pb.collection('collections').delete(collection.id);
							window.SetColored('red');
							collections = collections.filter((c: Collection) => c.id !== collection.id);
							if (page.url.pathname.includes(collection.id)) {
								goto('/inbox');
							}
						} catch (err) {
							console.error('Failed to delete collection', err);
							toast.error('Failed to delete collection!');
						}
					}}
					shareSnapshot={async () => {
						try {
							const cacheKey = `${collection.id}:cache`;
							let data: any[] = [];
							try {
								data = JSON.parse((await kv.get(cacheKey)) || '[]') || [];
							} catch {}

							const items = (data || []).map((b: any) => ({
								label: b.label,
								link: b.link,
								created: b.created,
								_favicon_base64: b._favicon_base64,
								_cover_base64: b._cover_base64
							}));

							if (!items.length) {
								toast.error('Niets om te delen in deze collectie.');
								return;
							}

							const snapshot = await pb.collection('snapshots').create({
								name: collection.name,
								user: pb.authStore.record?.id,
								collection: collection.id,
								data: items
							});

							const shareUrl = `${location.origin}/s/${snapshot.slug || snapshot.id}`;
							await navigator.clipboard.writeText(shareUrl);
							toast.success('Snapshot link copied to clipboard!', { description: shareUrl });
						} catch (e) {
							console.error('Failed to create snapshot', e);
							toast.error('Failed to create snapshot');
						}
					}}
				/>
			{/each}
			{#if nCollection.open}
				<div
					in:flyAndScale={{ x: 2.5, y: 0, duration: 150 }}
					class="flex items-center relative gap-2 w-full duration-200 bg-stone-100 dark:bg-stone-800/60 rounded-lg px-2 py-1"
				>
					<Folder class="size-4 opacity-80" />
					<input
						type="text"
						class="w-full text-sm p-1 m-0 rounded-lg bg-transparent border-none focus:outline-none focus:ring-0"
						placeholder="Start typing..."
						aria-label="Input field for typing"
						onkeyup={async (e) => {
							if (e.key === 'Enter') {
								nCollection.name = e.currentTarget.value || '';
								nCollection.open = false;
								const target = e.target as HTMLInputElement;
								if (target.value) {
									if (
										collections.find((collection: Collection) => collection.name === target.value)
									) {
										toast.error('Collection already exists!');
										return;
									}
									try {
										const userId = pb.authStore.model?.id;
										const newCol = await pb
											.collection('collections')
											.create({ name: target.value, user: userId });
										collections.push({
											id: newCol.id,
											name: newCol.name
										});
										goto(`/${newCol.id}`);
										window.SetColored('green');
										toast.info('Collection created!');
									} catch (err) {
										console.error('Failed to create collection', err);
										window.SetColored('red');

										toast.error('Failed to create collection!');
									}
								}
							}
							if (e.key === 'Escape') {
								nCollection.open = false;
							}
						}}
						use:focus
					/>
				</div>
			{/if}
		</div>
	</div>
</div>
