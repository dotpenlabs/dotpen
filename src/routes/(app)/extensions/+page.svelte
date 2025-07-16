<script>
	import { onMount } from 'svelte';
	import StoreItem from '$/lib/components/storeItem/storeItem.svelte';
	import { fade } from 'svelte/transition';

	let search = '';
	let items = [];

	// Simulate loading extensions (replace with real data if available)
	onMount(() => {
		items = Array.from({ length: 200 }, (_, index) => ({
			label: `Extension ${index}`,
			icon: '/full.png',
			url: `/extension/${index}`
		}));
	});

	// Filter items based on search
	$: filteredItems = items.filter((item) =>
		item.label.toLowerCase().includes(search.toLowerCase())
	);
</script>

<content class="h-full w-full p-4 overflow-y-auto flex flex-col gap-4">
	<div class="flex justify-between items-center w-full min-h-16 flex-nowrap gap-2">
		<div class="flex flex-col gap-1" in:fade={{ duration: 400 }} out:fade={{ duration: 200 }}>
			<div class="flex gap-2 items-center">
				<p class="text-xl font-semibold">Extensions</p>
				<p
					class="text-xs px-3 select-none !font-sans h-5 flex justify-center items-center rounded-xl bg-orange-600 hover:bg-orange-500 duration-200 cursor-pointer text-white"
				>
					alpha
				</p>
			</div>
			<p class="text-sm text-black/65">
				Enhance your experience with powerful extensions and plugins
			</p>
		</div>
		<input
			class="rounded-lg border border-black/15 px-3 py-2 text-sm w-64 focus:outline-none focus:ring-2 focus:ring-blue-400"
			placeholder="Search extensions..."
			bind:value={search}
			aria-label="Search extensions"
		/>
	</div>

	{#if filteredItems.length > 0}
		<div class="flex flex-row gap-3 flex-wrap justify-center items-start">
			{#each filteredItems as item, index}
				<StoreItem {item} {index} removeItem={() => {}} />
			{/each}
		</div>
	{:else}
		<div class="text-center text-gray-500 mt-8 w-full">
			<p>No extensions match your search.</p>
		</div>
	{/if}
</content>
