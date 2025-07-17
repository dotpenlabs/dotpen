<script lang="ts" module>
	declare global {
		interface Window {
			SetHydrating: (id: string, value: boolean) => void;
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

	import { Select } from 'bits-ui';

	import Loading from '$/lib/components/loading.svelte';
	import { flyAndScale } from '$/lib/utils';
	import { fade, fly } from 'svelte/transition';
	import { mode, setMode } from 'mode-watcher';
	import Item from '$lib/components/navigation/navitem.svelte';
	import { toast } from 'svelte-sonner';

	import { focus } from '$/lib/utils';
	import { goto } from '$app/navigation';
	import { pb } from '$/lib';
	import { page } from '$app/state';
	import DropBox from '$/lib/components/dropbox.svelte';

	let { data, children }: { data: LayoutData; children: Snippet } = $props();

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
				sort: '-created'
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

		loading = false;

		if (!pb.authStore.isValid) {
			goto('/');
		}

		await pb.collection('users').authRefresh();

		GetCollectionsCache();
		fetchCollections();

		BrowserTitleTick();
	});

	$effect(BrowserTitleTick);
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
{:else}
	<content class="bg-[#FAF5F2] dark:bg-[#171616] absolute flex h-full w-full p-3 overflow-hidden">
		<div
			class="h-full w-full bg-white dark:bg-stone-950 border border-stone-200 dark:border-stone-800 rounded-xl overflow-hidden"
		>
			<PaneGroup direction="horizontal">
				<div class="flex flex-col gap-1 h-full w-128">
					<div class="h-full w-full flex flex-col justify-between">
						<div class="flex-1 flex flex-col items-center justify-start gap-1 p-4 relative">
							<div class="h-full w-full flex flex-col justify-between">
								<div class="flex-1 flex flex-col items-center justify-start gap-1 p-4 relative">
									<div class="h-14 w-full">
										<p class="text-xl font-medium">
											{pb.authStore.record.name
												? 'Howdy, ' + pb.authStore.record.name.split(' ')[0] + '!'
												: 'Howdy!'}
										</p>
										{#if !isHydrating}
											<p in:fly={{ duration: 400, y: 5 }} class="text-xs opacity-65">
												All changes are synchronized with the cloud.
											</p>
										{:else}
											<p in:fly={{ duration: 350, y: -10 }} class="text-xs opacity-65">
												Changes are currently being synced...
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
												goto('/extensions');
											}}
											class="flex gap-2 text-[14.5px] justify-center rounded-xl cursor-pointer active:scale-98 active:opacity-90 duration-50 items-center border h-full w-full border-[#E4E4E4] dark:border-[#1d1d1d]"
										>
											<PuzzlePiece />
											Extensions
										</button>
									</div>
									<Item url="/inbox" icon={Tray} label="Inbox" />
									<div class="flex items-center gap-3 w-full opacity-50 my-2 mt-4">
										<p class="text-xs">Collections</p>
										<div class="flex-1 h-px bg-stone-300 dark:bg-stone-700 rounded-lg"></div>
									</div>
									<button
										class="flex flex-col gap-1 w-full h-96 overflow-y-auto pb-8"
										onkeyup={(e) => {
											if (e.key === 'n') {
												states.newCollection = !states.newCollection;
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
														Just double click (or press "n" when selected) here to create one, so
														simple is it!
													</p>
												</div>
											</div>
										{/if}
										{#each collections as collection}
											<Item
												url={collection.url}
												icon={collection.icon}
												label={collection.name}
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
												transition:flyAndScale={{ duration: 100, y: -5 }}
												class="flex items-center relative gap-2 cursor-pointer w-full duration-200 ring-1 ring-stone-200 dark:ring-stone-800 bg-stone-100 dark:bg-stone-800 rounded-lg px-2 py-1"
											>
												<Folder class="size-4 opacity-80" />
												<input
													type="text"
													class="w-full text-sm p-1 m-0 rounded-lg bg-transparent border-none focus:outline-none focus:ring-0"
													placeholder="Start typing..."
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
							</div>
							<DropBox />
							<div
								class="flex gap-3 items-center opacity-80 justify-start px-7 w-full -mt-6 text-xs"
							>
								<p>✸ Dotpen Early Access</p>
								<a class="opacity-65 hover:opacity-100 duration-200" href="/legal/privacy">
									Privacy
								</a>
								<a class="opacity-65 hover:opacity-100 duration-200" href="/legal/terms">Terms</a>
								<a class="opacity-65 hover:opacity-100 duration-200" href="/auth/logout">Logout</a>
							</div>
						</div>
					</div>
				</div>
				<div class="bg-stone-200 dark:bg-stone-800 w-px h-full"></div>

				<div class="h-full w-full p-5 overflow-x-clip">
					{@render children()}
				</div>
			</PaneGroup>
		</div>
	</content>
{/if}
