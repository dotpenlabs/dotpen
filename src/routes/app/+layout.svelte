<script lang="ts">
	import { onMount, type Snippet } from 'svelte';
	import type { LayoutData } from './$types';

	import { PaneGroup, Pane, PaneResizer } from 'paneforge';
	import { Book, FigmaLogo, Gear, SignOut, X } from 'phosphor-svelte';

	import Loading from '$/lib/components/loading.svelte';
	import {
		AppWindowMac,
		ArrowLeftToLine,
		ArrowRightToLine,
		Brush,
		Folder,
		Gift,
		Inbox,
		Sparkle,
		Sparkles
	} from 'lucide-svelte';
	import { flyAndScale } from '$/lib/utils';
	import { fade, fly } from 'svelte/transition';
	import { mode, setMode } from 'mode-watcher';
	import Item from './navitem.svelte';
	import { toast } from 'svelte-sonner';

	import { focus } from '$/lib/utils';
	import { goto } from '$app/navigation';

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

	let collections = $state([
		{
			name: 'Drawing inspiration',
			icon: Brush,
			url: '/app/drawing'
		},
		{
			name: 'Present ideas for John',
			icon: Gift,
			url: '/app/present'
		},
		{
			name: 'Writing a book',
			icon: Book,
			url: '/app/book'
		}
	]);

	onMount(() => {
		loading = false;
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
				{#if states.showSidebar}
					<Pane
						class="flex flex-col gap-1 h-full w-full"
						defaultSize={20}
						minSize={15}
						maxSize={25}
					>
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
											<Item url="/app/inbox" icon={Inbox} label="Inbox" />
											<div class="flex items-center gap-3 w-full opacity-50 my-2">
												<p class="text-xs">Suggestions</p>
												<div class="flex-1 h-px bg-stone-300 dark:bg-stone-700 rounded-lg"></div>
											</div>
											<Item url="/app/recent" icon={Sparkles} label="Recently added" />
											<Item url="/app/suggestions" icon={FigmaLogo} label="AI Suggested" />
											<div class="flex items-center gap-3 w-full opacity-50 my-2">
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
															onkeyup={(e) => {
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
																		collections.push({
																			name: target.value,
																			icon: Folder,
																			url: `/app/${target.value}`
																		});
																		goto(`/app/${target.value}`);
																	}
																	toast.info('Collection created!');
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
										<div
											class="h-16 border-t border-stone-200 dark:border-stone-800 w-full flex justify-between px-8 py-6 items-center"
										>
											{data.user?.user_metadata.name}
											<div class="flex gap-2">
												<button
													onclick={() => {
														states.showSettings = !states.showSettings;
													}}
													class="outline group cursor-pointer outline-stone-200 dark:outline-stone-800 hover:bg-stone-100 dark:hover:bg-stone-950 rounded-lg p-2"
												>
													<Gear
														class="size-4 group-active:scale-95 group-hover:scale-105 duration-200"
													/>
												</button>
												<button
													onclick={() => {
														states.returnToClosed = !states.returnToClosed;
													}}
													class="outline group cursor-pointer outline-stone-200 dark:outline-stone-800 hover:bg-stone-100 dark:hover:bg-stone-950 rounded-lg p-2"
													title={states.returnToClosed ? 'Open sidebar' : 'Auto close sidebar'}
												>
													{#if !states.returnToClosed}
														<ArrowLeftToLine
															class="size-4 group-active:scale-95 group-hover:scale-105 duration-200"
														/>
													{:else}
														<ArrowRightToLine
															class="size-4 group-active:scale-95 group-hover:scale-105 duration-200"
														/>
													{/if}
												</button>
											</div>
										</div>
									</div>
								</div>
							</div>
						</div>
					</Pane>
				{/if}
				<PaneResizer class="w-5 flex justify-start items-center">
					<div class="bg-stone-200 dark:bg-stone-800 w-px h-full"></div>
				</PaneResizer>
				<Pane class={'pt-3'}>
					{@render children()}
				</Pane>
			</PaneGroup>
		</div>
	</content>

	{#if states.showSettings}
		<button
			transition:fade
			aria-label="Close settings"
			class="absolute flex h-full w-full justify-center items-center inset-0 bg-stone-900/50 backdrop-blur-sm cursor-pointer"
			onclick={() => {
				states.showSettings = false;
			}}
		></button>
		<div class="absolute flex justify-center items-center inset-0 pointer-events-none">
			<div
				in:flyAndScale={{ duration: 250, delay: 200 }}
				out:flyAndScale={{ duration: 100 }}
				class="bg-white flex dark:bg-stone-950 outline outline-stone-200 dark:outline-stone-800 h-[80%] w-[60%] rounded-lg pointer-events-auto"
			>
				<div
					class="flex flex-col gap-2 w-64 border-r border-stone-200 dark:border-stone-800 h-full p-6"
				>
					<div class="flex items-center gap-2">
						<p class="text-sm">Settings</p>
						<div class="flex-1 h-px bg-stone-300 dark:bg-stone-700 rounded-lg"></div>
					</div>
					<div
						class="flex items-center active:opacity-65 outline outline-stone-200 dark:outline-stone-800 select-none gap-2 cursor-pointer w-full duration-200 hover:bg-stone-300 dark:hover:bg-stone-800 rounded-lg px-2 py-1"
					>
						<AppWindowMac class="size-4 opacity-80" />
						<p class="text-sm">General</p>
					</div>
				</div>
				<div class="relative h-full w-full p-6 bg-white dark:bg-stone-950 rounded-lg shadow-lg">
					<button
						class="absolute top-4 right-4 text-gray-600 hover:text-gray-800 dark:text-gray-300 dark:hover:text-gray-100 transition duration-200"
						onclick={() => {
							states.showSettings = false;
						}}
					>
						<X class="size-5" />
					</button>

					<h2 class="text-lg font-semibold">General</h2>
					<div class="w-full h-px bg-stone-200 dark:bg-stone-800 rounded-lg my-4"></div>

					<div class="flex justify-between items-center gap-4">
						<div>
							<p class="text-sm font-semibold">Theme</p>
							<p class="text-sm opacity-65 w-96">
								Choose your preferred theme to ensure the UI adapts dynamically to your personal
								preferences.
							</p>
						</div>

						<select
							onchange={(e) => {
								const target = e.target as HTMLInputElement;
								setMode(target.value as 'dark' | 'light' | 'system');
							}}
							value={$mode}
							class="bg-stone-200 dark:bg-stone-800 rounded-lg px-4 py-1 w-40"
						>
							<option value="light">Light</option>
							<option value="dark">Dark</option>
							<option value="system">System</option>
						</select>
					</div>
				</div>
			</div>
		</div>
	{/if}
{/if}
