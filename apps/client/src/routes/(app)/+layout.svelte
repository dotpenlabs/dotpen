<script lang="ts" module>
	declare global {
		interface Window {
			SetHydrating: (id: string, value: boolean) => void;
			SetColored: (isColored: 'default' | 'green' | 'red') => void;
		}
	}
</script>

<script lang="ts">
	import { onMount, type Snippet } from 'svelte';
	import type { LayoutData } from '../$types';

	import { PaneGroup, Pane, PaneResizer } from 'paneforge';
	import {
		Book,
		CaretDoubleDown,
		CaretDoubleUp,
		CaretUpDown,
		Check,
		FigmaLogo,
		Folder,
		Gear,
		MagnifyingGlass,
		Palette,
		PuzzlePiece,
		Sparkle,
		Tray,
		X
	} from 'phosphor-svelte';

	import Loading from '$/lib/components/loading.svelte';
	import { flyAndScale } from '$/lib/utils';
	import { fade, fly } from 'svelte/transition';
	import Item from '$lib/components/navigation/navitem.svelte';
	import { toast } from 'svelte-sonner';

	import { focus } from '$/lib/utils';
	import { goto } from '$app/navigation';
	import { pb } from '$/lib';
	import { page } from '$app/state';
	import View from '$/lib/plugins/view.svelte';

	let { children }: { children: Snippet } = $props();

	let online = $state(true);
	let loading = $state(true);
	let states = $state({
		quickActions: false,
		showSettings: false,
		newCollection: false,
		showSidebar: true,
		forceHidden: false,
		returnToClosed: false,
		collection_creating: false,
		collection_preview: ''
	});

	let collections = $state([]);
	let isHydrating = $state(false);
	let setColored = $state('default') as 'default' | 'green' | 'red';

	const hydrationStatus: Record<string, boolean> = {};
	const hydrationTimers: Record<string, NodeJS.Timeout> = {};

	const COLLECTIONS_CACHE_KEY = 'dotpen_collections';

	function GetCollectionsCache() {
		try {
			const cached = localStorage.getItem(COLLECTIONS_CACHE_KEY);
			if (cached) {
				const parsed = JSON.parse(cached);
				if (Array.isArray(parsed)) {
					collections = parsed;
				}
			}
		} catch (e) {
			console.error('Failed to load collections from cache', e);
		}
	}

	function SaveCollectionCache(cols: any[]) {
		try {
			localStorage.setItem(COLLECTIONS_CACHE_KEY, JSON.stringify(cols));
		} catch (e) {
			console.error('Failed to save collections to cache', e);
		}
	}

	async function fetchCollections() {
		if (!pb.authStore.isValid) return;
		try {
			const userId = pb.authStore.model?.id;
			if (!userId) return;
			const result = await pb.collection('collections').getFullList({
				filter: `user = "${userId}" && name != "system_inbox"`,
				sort: 'created'
			});
			const freshCollections = result.map((col) => ({
				id: col.id,
				name: col.name,
				icon: Folder,
				url: `/${col.id}`
			}));

			if (JSON.stringify(freshCollections) !== JSON.stringify(collections)) {
				collections = freshCollections;
				SaveCollectionCache(freshCollections);
			}
		} catch (err) {
			console.error('Failed to fetch collections', err);
			toast.error('Failed to load collections!');
		}
	}

	function BrowserTitleTick() {
		const currentPath = page.url.pathname;
		const currentCollection = collections.find((c) => c.url === currentPath);
		document.title = currentCollection ? `${currentCollection.name} • Dotpen` : 'Dotpen';

		switch (currentPath) {
			case '/inbox':
				document.title = 'Inbox • Dotpen';
				break;
			case '/extensions':
				document.title = 'Extensions • Dotpen';
				break;
		}
	}

	onMount(async () => {
		window.SetHydrating = (id: string, value: boolean) => {
			if (value) {
				if (hydrationTimers[id]) return;

				hydrationTimers[id] = setTimeout(() => {
					delete hydrationTimers[id];
					hydrationStatus[id] = true;

					const anyStillHydrating = Object.values(hydrationStatus).some((v) => v);
					isHydrating = anyStillHydrating;
				}, 125);
			} else {
				if (hydrationTimers[id]) {
					clearTimeout(hydrationTimers[id]);
					delete hydrationTimers[id];
				}

				hydrationStatus[id] = false;

				const anyStillHydrating = Object.values(hydrationStatus).some((v) => v);
				isHydrating = anyStillHydrating;
			}
		};

		window.SetColored = (isColored: 'default' | 'green' | 'red') => {
			setColored = isColored;
			setTimeout(() => {
				setColored = 'default';
			}, 1750);
		};

		loading = false;

		if (!pb.authStore.isValid) {
			window.location.href = '/auth';
			return;
		}

		await pb.collection('users').authRefresh();

		GetCollectionsCache();
		fetchCollections();
		BrowserTitleTick();
		UseUpdateManager();
		UseOnlineStatusMGR();
	});

	function UseOnlineStatusMGR() {
		setInterval(async () => {
			try {
				if (navigator.onLine === true) {
					online = true;
				} else if (navigator.onLine === false) {
					online = false;
				} else {
					await pb.health.check();
					online = true;
				}
			} catch (e) {
				console.error('[onlinemgr] device seems offline', e);
				online = false;
			}
		}, 3500);
	}

	function UseUpdateManager() {
		const run = () => {
			fetch('/api/version')
				.then((res) => res.json())
				.then((data) => {
					if (data.version !== APP_VERSION) {
						console.log('[updatemgr] Update available!');
						toast.info('Update available!', {
							description: 'A new version of Dotpen is available!',
							duration: Number.POSITIVE_INFINITY,
							action: {
								label: 'Update',
								onClick: () => {
									window.location.reload();
								}
							}
						});
					}
				});
		};

		run();
		setInterval(run, 1000 * 60 * 15);
	}

	$effect(BrowserTitleTick);

	let isMobile = $state(false);
	let mobileView = $state<'sidebar' | 'items'>('sidebar');
	let currentCollectionName = $state('');

	onMount(() => {
		const checkMobile = () => {
			isMobile = window.innerWidth < 768;
			if (page.url.pathname === '/share') mobileView = 'items';
		};
		checkMobile();
		window.addEventListener('resize', checkMobile);
		return () => window.removeEventListener('resize', checkMobile);
	});

	function handleSidebarItemClick({ url, label }) {
		if (isMobile) {
			currentCollectionName = label;
			mobileView = 'items';
		}
		goto(url);
	}
	function handleBack() {
		if (isMobile) {
			mobileView = 'sidebar';
		}
	}
</script>

<svelte:window
	onkeydown={(e) => {
		if (e.key === 'Escape') {
			states.showSettings = false;
		}

		if (e.key === ',' && e.ctrlKey) {
			states.showSettings = true;
		}
	}}
/>

{#if loading}
	<Loading />
{:else if isMobile}
	<div
		class="absolute inset-0 flex h-full w-full bg-[#FAF5F2] dark:bg-[#171616] transition-colors duration-1200 overflow-hidden"
	>
		<div
			class="absolute inset-0 h-full w-full transition-transform duration-300"
			style="transform: translateX({mobileView === 'sidebar' ? '0%' : '-100%'})"
		>
			<div
				class="h-full w-full bg-white dark:bg-stone-950 border border-stone-200 dark:border-stone-800 rounded-xl overflow-hidden flex flex-col"
			>
				<div class="flex flex-col gap-1 h-full w-full min-h-0">
					<div class="h-full w-full flex flex-col justify-between min-h-0">
						<div
							class="flex-1 flex flex-col items-center justify-start gap-1 p-4 relative w-full min-h-0"
						>
							<div class="h-14 w-full">
								<p class="text-xl font-medium">
									{pb.authStore.record.name
										? 'Howdy, ' + pb.authStore.record.name.split(' ')[0] + '!'
										: 'Howdy!'}
								</p>
								{#if online}
									{#if !isHydrating}
										<p in:fly={{ duration: 400, y: 5 }} class="text-xs opacity-65">
											All changes are synchronized with the cloud.
										</p>
									{:else}
										<p in:fly={{ duration: 350, y: -10 }} class="text-xs opacity-65">
											Changes are currently being synced...
										</p>
									{/if}
								{:else}
									<p in:fly={{ duration: 350, y: -10 }} class="text-xs opacity-65">
										Not connected to internet, data may be outdated.
									</p>
								{/if}
							</div>
							<div class="flex items-center gap-3 w-full opacity-50 my-2">
								<p class="text-xs">Menu</p>
								<div class="flex-1 h-px bg-stone-300 dark:bg-stone-700 rounded-lg"></div>
							</div>
							<div class="flex h-8 items-center gap-4 mb-4 justify-between w-full">
								<button
									class="flex gap-2 text-[14.5px] justify-center rounded-xl cursor-pointer active:scale-98 active:opacity-90 duration-50 items-center border h-full w-full border-[#E4E4E4] dark:border-[#1d1d1d]"
								>
									<MagnifyingGlass />Search
								</button>
								<button
									onclick={() => {
										goto('/extensions');
									}}
									class="flex gap-2 text-[14.5px] justify-center rounded-xl cursor-pointer active:scale-98 active:opacity-90 duration-50 items-center border h-full w-full border-[#E4E4E4] dark:border-[#1d1d1d]"
								>
									<PuzzlePiece />Extensions
								</button>
							</div>
							<Item url="/inbox" icon={Tray} label="Inbox" onSelect={handleSidebarItemClick} />
							<div class="flex items-center gap-3 w-full opacity-50 my-2 mt-4 relative">
								<p class="text-xs">Collections</p>
								<div class="flex-1 h-px bg-stone-300 dark:bg-stone-700 rounded-lg"></div>
								<button
									aria-label="Create new collection"
									class="absolute right-0 top-1/2 -translate-y-1/2 text-xs px-2 py-1 rounded-full bg-stone-200 dark:bg-stone-700 hover:bg-stone-300 dark:hover:bg-stone-600 focus:outline-none focus:ring-2 focus:ring-stone-400"
									style="z-index:2;"
									onclick={() => {
										states.newCollection = true;
									}}
								>
									+
								</button>
							</div>
							<div class="flex-1 w-full min-h-0 overflow-y-auto pb-8">
								<button
									oncontextmenu={(e: Event) => {
										e.preventDefault();
										states.newCollection = !states.newCollection;
									}}
									class="flex flex-1 flex-col h-full justify-start items-center gap-1 w-full"
								>
									{#if collections.length === 0 && !states.newCollection && !states.collection_creating}
										<div
											in:fade={{ duration: 500 }}
											class="flex flex-col h-full w-full gap-2 justify-center items-center opacity-65"
										>
											<div class="text-center">
												<p class="text-lg font-medium">No collections?</p>
												<p class="text-xs w-64">
													Long hold the empty space to create one, so simple is it!
												</p>
											</div>
										</div>
									{/if}
									{#each collections as collection}
										<div class="w-full">
											<Item
												url={collection.url}
												icon={/^\p{Emoji}/u.test(collection.name) ? null : Folder}
												label={collection.name}
												onSelect={handleSidebarItemClick}
												removeCollection={async () => {
													await pb.collection('collections').delete(collection.id);
													collections = collections.filter((c) => c.id !== collection.id);
													SaveCollectionCache(collections);
													if (page.url.pathname.includes(collection.id)) {
														goto('/');
														return;
													}
													toast.success('Collection deleted!');
												}}
											/>
										</div>
									{/each}
									{#if states.collection_creating}
										<div
											role="dialog"
											class="flex items-center relative gap-2 w-full duration-200 ring-1 ring-stone-200 dark:ring-stone-800 bg-stone-100 dark:bg-stone-800 rounded-lg px-2 py-1 animate-pulse opacity-70"
										>
											<Folder class="size-4 opacity-80" />
											<span
												class="w-full text-sm p-1 text-left m-0 rounded-lg bg-transparent border-none focus:outline-none focus:ring-0"
												>{states.collection_preview || 'Creating...'}</span
											>
										</div>
									{/if}
									{#if states.newCollection}
										<div
											in:flyAndScale={{ duration: 100, y: -5 }}
											class="flex items-center relative gap-2 cursor-pointer w-full duration-200 ring-1 ring-stone-200 dark:ring-stone-800 bg-stone-100 dark:bg-stone-800 rounded-lg px-2 py-1"
										>
											<Folder class="size-4 opacity-80" />
											<input
												type="text"
												class="w-full text-sm p-1 m-0 rounded-lg bg-transparent border-none focus:outline-none focus:ring-0"
												placeholder="Start typing... (press esc to exit)"
												aria-label="Input field for typing"
												onkeyup={async (e) => {
													if (e.key === 'Enter') {
														states.newCollection = false;
														const target = e.target as HTMLInputElement;
														if (target.value) {
															if (
																collections.find((collection) => collection.name === target.value)
															) {
																toast.error('Collection already exists!');
																return;
															}
															try {
																states.collection_creating = true;
																states.collection_preview = target.value;
																const userId = pb.authStore.model?.id;
																const newCol = await pb
																	.collection('collections')
																	.create({ name: target.value, user: userId });
																collections.push({
																	id: newCol.id,
																	name: newCol.name,
																	icon: Folder,
																	url: `/${newCol.id}`
																});
																SaveCollectionCache(collections);
																states.collection_creating = false;
																states.collection_preview = '';
																goto(`/${newCol.id}`);
																toast.info('Collection created!');
															} catch (err) {
																states.collection_creating = false;
																states.collection_preview = '';
																console.error('Failed to create collection', err);
																toast.error('Failed to create collection!');
															}
														}
													}
													if (e.key === 'Escape') {
														states.newCollection = false;
													}
												}}
												use:focus
											/>
										</div>
									{/if}
								</button>
							</div>
						</div>
					</div>
				</div>
			</div>
		</div>
		<div
			class="absolute inset-0 h-full w-full bg-white dark:bg-stone-950 border border-stone-200 dark:border-stone-800 rounded-xl overflow-hidden flex flex-col transition-transform duration-300 min-h-0"
			style="transform: translateX({mobileView === 'items' ? '0%' : '100%'})"
		>
			<div
				class="flex items-center h-14 px-4 border-b border-stone-200 dark:border-stone-800 bg-white/80 dark:bg-stone-950/80 backdrop-blur-xl flex-shrink-0"
			>
				<button class="mr-2 text-xl" onclick={handleBack} aria-label="Back">←</button>
				<p class="text-base font-semibold flex-1 text-center">{currentCollectionName}</p>
				{#if pb.authStore.record?.avatar}
					<img
						src={pb.files.getURL(pb.authStore.record, pb.authStore.record.avatar, {
							thumb: '32x32'
						})}
						alt="avatar"
						class="ml-2 w-8 h-8 rounded-full object-cover border border-stone-200 dark:border-stone-800 cursor-pointer"
					/>
				{:else}
					<div class="ml-2 w-8 h-8 rounded-full bg-stone-200 dark:bg-stone-800"></div>
				{/if}
			</div>
			<div class="flex-1 min-h-0 overflow-y-auto">
				<div class="h-full w-full p-5 overflow-x-clip">
					{@render children()}
				</div>
			</div>
		</div>
		<nav
			class="fixed bottom-0 left-0 right-0 bg-white dark:bg-stone-950 border-t border-stone-300 dark:border-stone-700 flex justify-around py-2 z-50"
		>
			<button onclick={() => (mobileView = 'sidebar')}>
				<Tray class="mx-auto" />
				<span class="text-xs">Home</span>
			</button>

			<button onclick={() => (states.showSettings = true)}>
				<Gear class="mx-auto" />
				<span class="text-xs">Settings</span>
			</button>
		</nav>
	</div>
{:else}
	<content
		data-vaul-drawer-wrapper
		class={`bg-[#FAF5F2] dark:bg-[#171616] absolute flex h-full w-full p-3 overflow-hidden transition-colors duration-1200 ` +
			(setColored === 'green'
				? 'bg-green-100 dark:bg-green-900'
				: '' + setColored === 'red'
					? 'bg-red-100 dark:bg-red-900'
					: '')}
	>
		<div
			class="h-full w-full bg-white dark:bg-stone-950 border border-stone-200 dark:border-stone-800 rounded-xl overflow-hidden"
		>
			<PaneGroup direction="horizontal">
				<div class="flex flex-col gap-1 h-full w-128 px-4 pt-4 min-h-0">
					<div class="h-full w-full flex flex-col justify-between min-h-0">
						<div
							class="flex-1 flex flex-col items-center justify-start gap-1 p-4 relative w-full min-h-0"
						>
							<div class="h-14 w-full">
								<p class="text-xl font-medium">
									{pb.authStore.record.name
										? 'Howdy, ' + pb.authStore.record.name.split(' ')[0] + '!'
										: 'Howdy!'}
								</p>
								{#if online}
									{#if !isHydrating}
										<p in:fly={{ duration: 400, y: 5 }} class="text-xs opacity-65">
											All changes are synchronized with the cloud.
										</p>
									{:else}
										<p in:fly={{ duration: 350, y: -10 }} class="text-xs opacity-65">
											Changes are currently being synced...
										</p>
									{/if}
								{:else}
									<p in:fly={{ duration: 350, y: -10 }} class="text-xs opacity-65">
										Not connected to internet, data may be outdated.
									</p>
								{/if}
							</div>
							<div class="flex items-center gap-3 w-full opacity-50 my-2">
								<p class="text-xs">Menu</p>
								<div class="flex-1 h-px bg-stone-300 dark:bg-stone-700 rounded-lg"></div>
							</div>
							<div class="flex h-8 items-center gap-4 mb-4 justify-between w-full">
								<button
									class="flex gap-2 text-[14.5px] justify-center rounded-xl cursor-pointer active:scale-98 active:opacity-90 duration-50 items-center border h-full w-full border-[#E4E4E4] dark:border-[#1d1d1d]"
								>
									<MagnifyingGlass />
									Search
								</button>
								<button
									onclick={() => {
										goto('/plugins');
									}}
									class="flex gap-2 text-[14.5px] justify-center rounded-xl cursor-pointer active:scale-98 active:opacity-90 duration-50 items-center border h-full w-full border-[#E4E4E4] dark:border-[#1d1d1d]"
								>
									<PuzzlePiece />
									Plugins
								</button>
							</div>
							<Item url="/inbox" icon={Tray} label="Inbox" onSelect={handleSidebarItemClick} />
							<div class="flex items-center gap-3 w-full my-2 mt-4 relative">
								<p class="text-xs opacity-50">Collections</p>
								<div class="flex-1 h-px bg-stone-300 dark:bg-stone-700 rounded-lg opacity-50"></div>
								<button
									aria-label="Create new collection"
									class="absolute right-0 top-1/2 -translate-y-1/2 opacity-75 cursor-pointer active:scale-95 aspect-square text-base text-black px-2 rounded-full bg-white"
									style="z-index:2;"
									onclick={() => {
										states.newCollection = true;
									}}
								>
									+
								</button>
							</div>
							<div class="flex-1 w-full min-h-0 overflow-y-auto pb-8">
								<button
									class="flex flex-col gap-1 w-full h-full"
									onkeyup={(e) => {
										if (e.key === 'n' && !states.newCollection) {
											states.newCollection = true;
										}
									}}
									ondblclick={() => {
										states.newCollection = !states.newCollection;
									}}
								>
									{#if collections.length === 0 && !states.newCollection && !states.collection_creating}
										<div
											in:fade={{ duration: 500 }}
											class="flex flex-col h-full w-full gap-2 justify-center items-center opacity-65"
										>
											<div class="text-center">
												<p class="text-lg font-medium">No collections?</p>
												<p class="text-xs w-64">
													{#if !isMobile}
														Just double click (or press "n" when selected) here to create one, so
														simple is it!
													{:else}
														Just tap anywhere twice to create one, so simple is it!
													{/if}
												</p>
											</div>
										</div>
									{/if}
									{#each collections as collection}
										<Item
											url={collection.url}
											icon={/^\p{Emoji}/u.test(collection.name) ? null : Folder}
											label={collection.name}
											onSelect={handleSidebarItemClick}
											removeCollection={async () => {
												try {
													await pb.collection('collections').delete(collection.id);
													collections = collections.filter((c) => c.id !== collection.id);
													SaveCollectionCache(collections);
													if (page.url.pathname.includes(collection.id)) {
														goto('/');
													}
													toast.success('Collection deleted!');
												} catch (err) {
													console.error('Failed to delete collection', err);
													toast.error('Failed to delete collection!');
												}
											}}
										/>
									{/each}
									{#if states.collection_creating}
										<div
											role="dialog"
											class="flex items-center relative gap-2 w-full duration-200 ring-1 ring-stone-200 dark:ring-stone-800 bg-stone-100 dark:bg-stone-800 rounded-lg px-2 py-1 animate-pulse opacity-70"
										>
											<Folder class="size-4 opacity-80" />
											<span
												class="w-full text-sm p-1 text-left m-0 rounded-lg bg-transparent border-none focus:outline-none focus:ring-0"
												>{states.collection_preview || 'Creating...'}</span
											>
										</div>
									{/if}
									{#if states.newCollection}
										<div
											in:flyAndScale={{ duration: 100, y: -5 }}
											class="flex items-center relative gap-2 cursor-pointer w-full duration-200 ring-1 ring-stone-200 dark:ring-stone-800 bg-stone-100 dark:bg-stone-800 rounded-lg px-2 py-1"
										>
											<Folder class="size-4 opacity-80" />
											<input
												type="text"
												class="w-full text-sm p-1 m-0 rounded-lg bg-transparent border-none focus:outline-none focus:ring-0"
												placeholder="Start typing... (press esc to exit)"
												aria-label="Input field for typing"
												onkeyup={async (e) => {
													if (e.key === 'Enter') {
														states.newCollection = false;
														const target = e.target as HTMLInputElement;
														if (target.value) {
															if (
																collections.find((collection) => collection.name === target.value)
															) {
																toast.error('Collection already exists!');
																return;
															}
															try {
																states.collection_creating = true;
																states.collection_preview = target.value;
																const userId = pb.authStore.model?.id;
																const newCol = await pb.collection('collections').create({
																	name: target.value,
																	user: userId
																});
																collections.push({
																	id: newCol.id,
																	name: newCol.name,
																	icon: Folder,
																	url: `/${newCol.id}`
																});
																SaveCollectionCache(collections);
																states.collection_creating = false;
																states.collection_preview = '';
																goto(`/${newCol.id}`);
																toast.info('Collection created!');
															} catch (err) {
																states.collection_creating = false;
																states.collection_preview = '';
																console.error('Failed to create collection', err);
																toast.error('Failed to create collection!');
															}
														}
													}
													if (e.key === 'Escape') {
														states.newCollection = false;
													}
												}}
												use:focus
											/>
										</div>
									{/if}
								</button>
							</div>
							<div class="h-32 overflow-clip mb-7">
								<View location="sidebar-bottom" />
							</div>
							<div
								class="flex flex-row gap-3 items-center justify-start px-3 w-full -mt-5 text-xs opacity-80"
							>
								<p>✸ Dotpen Early Access {APP_VERSION}</p>
								<a class="opacity-65 hover:opacity-100 duration-200" href="/legal/privacy"
									>Privacy</a
								>
								<a class="opacity-65 hover:opacity-100 duration-200" href="/legal/terms">Terms</a>
								<a class="opacity-65 hover:opacity-100 duration-200" href="/auth/logout">Logout</a>
							</div>
						</div>
					</div>
				</div>
				<div class="bg-stone-200 dark:bg-stone-800 w-px h-full"></div>
				<div class="h-full w-full p-5 overflow-x-clip min-h-0">
					{@render children()}
				</div>
			</PaneGroup>
		</div>
	</content>
{/if}
