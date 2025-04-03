<script lang="ts">
	import { ArrowSquareOut } from 'phosphor-svelte';
	import { fly } from 'svelte/transition';
	import type { LinkItem } from '$/lib/components/item/types';
	import { goto } from '$app/navigation';

	let {
		data: item = {
			label: '',
			icon: '',
			url: '',
			cover: '',
			date: ''
		},
		index = 0
	}: {
		data: LinkItem;
		index: number;
	} = $props();
</script>

<button
	onclick={() => {
		goto('/system/out?url=' + item.url);
	}}
	data-packed={index}
	class="text-wrap group cursor-pointer bg-white w-64 overflow-hidden h-min dark:bg-stone-950 rounded-xl p-5 shadow-sm hover:shadow-lg dark:shadow-none border border-zinc-200 dark:border-stone-700/50 duration-300 hover:scale-102 active:scale-98"
	in:fly|global={{ y: 20, duration: 500, delay: Math.min(index * 50, 1000) }}
	onmouseenter={(e) => {
		e.currentTarget.style.transform = `rotate(${Math.random() < 0.5 ? '-' : ''}${Math.floor(Math.random() * 1.5)}deg)`;
	}}
	onmouseleave={(e) => {
		e.currentTarget.style.transform = 'rotate(0deg)';
	}}
>
	<div class="relative flex flex-col items-start gap-3">
		<div class="absolute w-full h-full left-0 top-0">
			<a
				href={item.url}
				class="absolute duration-200 opacity-0 group-hover:opacity-100 bg-zinc-900 dark:bg-zinc-100 text-white dark:text-zinc-900 rounded-full p-1 z-[100] -right-1.5 -top-1.5"
			>
				<ArrowSquareOut class="size-4" weight="bold" />
			</a>
		</div>
		<div>
			{#if item.cover}
				<img
					src={item.cover}
					alt={item.label}
					class="w-full h-full outline outline-black/10 dark:outline-white/10 object-cover transform transition-transform duration-500 max-h-48 rounded-lg mb-4"
					onerror={(e: Event) => {
						const target = e.target as HTMLImageElement;
						if (target) {
							target.remove();
						}
					}}
				/>
			{/if}
			<div class="flex justify-between items-center w-52">
				<div>
					<p
						class="flex gap-2 items-center font-serif text-base leading-tight font-medium text-zinc-700 dark:text-zinc-300 duration-300"
						title={item.label}
					>
						{item.label.split('').slice(0, 28).join('')}{item.label.split('').length > 28
							? '...'
							: ''}
						{#if item.icon}
							<img src={item.icon} alt={item.label} class="size-4 rounded-sm" />
						{/if}
					</p>
					<div class="w-8 h-px bg-zinc-300 dark:bg-zinc-600 mt-2 mb-1"></div>
					<span class="text-xs font-serif italic text-zinc-500 dark:text-zinc-400">
						{new URL(item.url).hostname} - {new Date(item.date).toLocaleDateString()}
					</span>
				</div>
			</div>
		</div>
	</div>
</button>
