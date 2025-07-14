<script lang="ts">
	import { onMount, type Snippet } from 'svelte';
	import type { LayoutData } from './$types';

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
		Sparkle,
		Tray,
		X
	} from 'phosphor-svelte';

	import { Select } from 'bits-ui';

	import Loading from '$/lib/components/loading.svelte';
	import { flyAndScale } from '$/lib/utils';
	import { fade, fly } from 'svelte/transition';
	import { mode, setMode } from 'mode-watcher';
	import Item from './navitem.svelte';
	import { toast } from 'svelte-sonner';

	import { focus } from '$/lib/utils';
	import { goto } from '$app/navigation';
	import { pb } from '$/lib';
	import { page } from '$app/state';

	let { data, children }: { data: LayoutData; children: Snippet } = $props();

	let loading = $state(true);
	let states = $state({
		quickActions: false,
		showSettings: false,
		newCollection: false,
		showSidebar: true,
		forceHidden: false,
		returnToClosed: false
	});

	let collections = $state([]);

	async function fetchCollections() {
		if (!pb.authStore.isValid) return;
		try {
			const userId = pb.authStore.model?.id;
			if (!userId) return;
			const result = await pb.collection('collections').getFullList({
				filter: `user = "${userId}" && name != "system_inbox"`,
				sort: '-created'
			});
			collections = result.map((col) => ({
				id: col.id,
				name: col.name,
				icon: Folder,
				url: `/app/${col.id}`
			}));
		} catch (err) {
			console.error('Failed to fetch collections', err);
			toast.error('Failed to load collections!');
		}
	}

	function updateBrowserTitle() {
		const currentPath = page.url.pathname;
		const currentCollection = collections.find((c) => c.url === currentPath);
		document.title = currentCollection ? `${currentCollection.name} • Dotpen` : 'Dotpen';
		if (currentPath.endsWith('/inbox')) {
			document.title = 'Inbox • Dotpen';
		}
	}

	onMount(async () => {
		loading = false;

		if (!pb.authStore.isValid) {
			const email = prompt('Please enter your email');
			const password = prompt('Please enter your password');
			await pb.collection('users').authWithPassword(email, password);
		}
		await pb.collection('users').authRefresh();
		await fetchCollections();

		updateBrowserTitle();
	});

	$effect(() => {
		updateBrowserTitle();
	});
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
		{#if !states.showSidebar && !states.forceHidden}
			<button
				class="z-10 absolute left-0 top-0 h-full w-4"
				aria-label="Open sidebar"
				onclick={() => {
					states.showSidebar = true;
				}}
				onmouseenter={() => {
					states.showSidebar = true;
				}}
			></button>
		{/if}
		<div
			class="h-full w-full bg-white dark:bg-stone-950 border border-stone-200 dark:border-stone-800 rounded-xl overflow-hidden"
		>
			<PaneGroup direction="horizontal">
				<div class="flex flex-col gap-1 h-full w-128">
					<div
						role="navigation"
						onmouseleave={() => {
							if (states.returnToClosed) states.showSidebar = false;
						}}
						onblur={() => {
							if (states.returnToClosed) states.showSidebar = false;
						}}
						class="h-full w-full"
						transition:fly={{ duration: 250, y: 0, x: -25 }}
					>
						<div class="h-full w-full flex flex-col justify-between">
							<div class="flex-1 flex flex-col items-center justify-start gap-1 p-4 relative">
								<div class="h-full w-full flex flex-col justify-between">
									<div class="flex-1 flex flex-col items-center justify-start gap-1 p-4 relative">
										<div class="h-20 w-full">
											<p class="text-xl font-medium">
												{pb.authStore.record.name
													? 'Howdy, ' + pb.authStore.record.name + '!'
													: 'Howdy!'}
											</p>
											<p class="text-xs opacity-65">All changes are synchronized with the cloud.</p>
										</div>
										<div class="flex items-center gap-3 w-full opacity-50 my-2">
											<p class="text-xs">Menu</p>
											<div class="flex-1 h-px bg-stone-300 dark:bg-stone-700 rounded-lg"></div>
										</div>
										<div class="flex h-12 items-center gap-4 mb-4 justify-between w-full">
											<button
												class="flex gap-2 text-[14.5px] justify-center rounded-xl cursor-pointer active:scale-98 active:opacity-90 duration-50 items-center border h-full w-full border-[#E4E4E4] dark:border-[#1d1d1d]"
											>
												<MagnifyingGlass />
												Search
											</button>
											<button
												class="flex gap-2 text-[14.5px] justify-center rounded-xl cursor-pointer active:scale-98 active:opacity-90 duration-50 items-center border h-full w-full border-[#E4E4E4] dark:border-[#1d1d1d]"
											>
												<Sparkle />
												Inspire me
											</button>
										</div>
										<Item url="/app/inbox" icon={Tray} label="Inbox" />
										<div class="flex items-center gap-3 w-full opacity-50 my-2 mt-4">
											<p class="text-xs">Collections</p>
											<div class="flex-1 h-px bg-stone-300 dark:bg-stone-700 rounded-lg"></div>
										</div>
										<button
											class="flex flex-col gap-1 w-full h-full"
											ondblclick={() => {
												states.newCollection = !states.newCollection;
											}}
										>
											{#each collections as collection}
												<Item
													url={collection.url}
													icon={collection.icon}
													label={collection.name}
													removeCollection={async () => {
														try {
															await pb.collection('collections').delete(collection.id);
															collections = collections.filter((c) => c.id !== collection.id);
															toast.success('Collection deleted!');
														} catch (err) {
															console.error('Failed to delete collection', err);
															toast.error('Failed to delete collection!');
														}
													}}
												/>
											{/each}
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
																		collections.find(
																			(collection) => collection.name === target.value
																		)
																	) {
																		toast.error('Collection already exists!');
																		return;
																	}
																	try {
																		const userId = pb.authStore.model?.id;
																		const newCol = await pb.collection('collections').create({
																			name: target.value,
																			user: userId
																		});
																		collections.push({
																			id: newCol.id,
																			name: newCol.name,
																			icon: Folder,
																			url: `/app/${newCol.id}`
																		});
																		goto(`/app/${newCol.id}`);
																		toast.info('Collection created!');
																	} catch (err) {
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
				<div class="bg-stone-200 dark:bg-stone-800 w-px h-full"></div>

				<div class="h-full w-full p-5 overflow-x-clip">
					{@render children()}
				</div>
			</PaneGroup>
		</div>
	</content>
{/if}
