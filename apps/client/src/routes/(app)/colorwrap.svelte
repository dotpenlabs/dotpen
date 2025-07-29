<script lang="ts">
	import { twMerge } from 'tailwind-merge';

	let { appplatform, appcolor, children } = $props();

	const colorClass = {
		green: 'bg-green-100 dark:bg-[#002B00]',
		red: 'bg-red-100 dark:bg-[#300000]'
	};
</script>

<content class="absolute h-full w-full inset-0 overflow-clip">
	{#if appplatform === 'desktop'}
		<div
			class={twMerge(
				'h-full w-full transition-colors duration-1200 overflow-hidden bg-[#E3E3E3] dark:bg-stone-900 p-4',
				colorClass[appcolor]
			)}
		>
			{@render children()}
		</div>
	{:else}
		{@render children()}
	{/if}
</content>
