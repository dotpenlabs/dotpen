<script lang="ts">
	import { ArrowSquareOut } from 'phosphor-svelte';

	let {
		data: item = {
			label: '',
			link: '',
			_cover_base64: '',
			_favicon_base64: '',
			created: new Date().toISOString()
		},
		index = 0
	}: {
		data: {
			label: string;
			link: string;
			_cover_base64?: string;
			_favicon_base64?: string;
			created?: string;
		};
		index?: number;
	} = $props();
</script>

<button
	title={item.label}
	onclick={() => {
		window.open(item.link, '_blank');
	}}
	data-packed={index}
	class="text-wrap group cursor-pointer bg-white w-full sm:w-64 overflow-hidden h-min dark:bg-stone-950 rounded-xl p-3 sm:p-5 shadow-sm hover:shadow-lg dark:shadow-none border border-zinc-200 dark:border-stone-700/50 duration-300 hover:scale-102 active:scale-98"
>
	<div class="relative flex flex-col items-start gap-2 sm:gap-3">
		<div class="absolute w-full h-full left-0 top-0">
			<a
				aria-label="Open link"
				href={item.link}
				onclick={(e) => e.preventDefault()}
				class="absolute duration-200 opacity-0 group-hover:opacity-100 bg-zinc-900 dark:bg-zinc-100 text-white dark:text-zinc-900 rounded-full p-1 z-[100] -right-1.5 -top-1.5"
			>
				<ArrowSquareOut />
			</a>
		</div>
		<div class="w-full">
			{#if item._cover_base64}
				<img
					draggable="false"
					src={item._cover_base64}
					alt={item.label}
					class="w-full h-full outline outline-black/10 dark:outline-white/10 object-cover transform transition-transform duration-500 max-h-32 sm:max-h-48 rounded-lg mb-3 sm:mb-4"
					onerror={(e: Event) => {
						const target = e.target as HTMLImageElement;
						if (target) target.remove();
					}}
				/>
			{/if}
			<div class="flex justify-between items-center w-full">
				<div class="w-full">
					<p
						class="flex gap-2 items-center font-serif text-sm sm:text-base leading-tight font-medium text-zinc-700 dark:text-zinc-300 duration-300"
					>
						<span class="text-xs sm:text-sm font-semibold truncate max-w-36 sm:max-w-52"
							>{item.label}</span
						>
						{#if item._favicon_base64}
							<img
								src={item._favicon_base64}
								alt={item.label}
								class="size-3 sm:size-4 rounded-sm flex-shrink-0"
							/>
						{/if}
					</p>
					<div class="w-6 sm:w-8 h-px bg-zinc-300 dark:bg-zinc-600 mt-1.5 sm:mt-2 mb-1"></div>
					<p class="text-xs w-full font-serif text-left italic text-zinc-500 dark:text-zinc-400">
						{new URL(item.link).hostname}{#if item.created}&nbsp;-&nbsp;{new Date(
								item.created
							).toLocaleDateString()}{/if}
					</p>
				</div>
			</div>
		</div>
	</div>
</button>
