<script lang="ts">
	import { ContextMenu } from 'bits-ui';
	import Trash from 'phosphor-svelte/lib/Trash';
	import { ArrowSquareOut } from 'phosphor-svelte';
	import { fly } from 'svelte/transition';
	import type { LinkItem } from '$/lib/components/item/types';

	type StoreItem = {
		label: string;
		icon: string;
		url: string;
		id?: string;
	};

	let {
		item = {
			label: '',
			icon: '',
			url: ''
		},
		index = 0,
		removeItem
	}: {
		item: StoreItem;
		index: number;
		removeItem: () => void;
	} = $props();
</script>

<ContextMenu.Root>
	<ContextMenu.Trigger>
		<button
			title={item.label}
			onclick={() => {
				window.open(item.url, '_blank');
			}}
			data-packed={index}
			class="text-wrap group cursor-pointer bg-white w-full sm:w-64 overflow-hidden h-min dark:bg-stone-950 rounded-xl p-3 sm:p-5 shadow-sm hover:shadow-lg dark:shadow-none border border-zinc-200 dark:border-stone-700/50 duration-300 hover:scale-102 active:scale-98"
			in:fly|global={{ y: 20, duration: 300, delay: Math.min(index * 35, 1000) }}
			out:fly|global={{ y: 10, duration: 150 }}
			onmouseenter={(e) => {
				e.currentTarget.style.transform = `rotate(${Math.random() < 0.5 ? '-' : ''}${Math.floor(Math.random() * 1.5)}deg)`;
			}}
			onmouseleave={(e) => {
				e.currentTarget.style.transform = 'rotate(0deg)';
			}}
		>
			<div class="relative flex flex-col items-start gap-2 sm:gap-3">
				<div class="absolute w-full h-full left-0 top-0">
					<a
						href={item.url}
						onclick={(e) => e.preventDefault()}
						class="absolute duration-200 opacity-0 group-hover:opacity-100 bg-zinc-900 dark:bg-zinc-100 text-white dark:text-zinc-900 rounded-full p-1 z-[100] -right-1.5 -top-1.5"
					>
						<ArrowSquareOut class="size-3 sm:size-4" weight="bold" />
					</a>
				</div>
				<div class="w-full">
					<div class="flex justify-between items-center w-full">
						<div class="w-full">
							<p
								class="flex gap-2 items-center font-serif text-sm sm:text-base leading-tight font-medium text-zinc-700 dark:text-zinc-300 duration-300"
							>
								{#if item.icon}
									<img
										src={item.icon}
										alt={item.label}
										class="size-3 sm:size-4 rounded-sm flex-shrink-0"
									/>
								{/if}
								<span class="text-xs sm:text-sm font-semibold truncate max-w-36 sm:max-w-52">
									{item.label}
								</span>
							</p>
							<div class="w-6 sm:w-8 h-px bg-zinc-300 dark:bg-zinc-600 mt-1.5 sm:mt-2 mb-1"></div>
							<p
								class="text-xs w-full font-serif text-left italic text-zinc-500 dark:text-zinc-400"
							>
								Dotpen Team – Recommanded
							</p>
						</div>
					</div>
				</div>
			</div>
		</button>
	</ContextMenu.Trigger>

	<ContextMenu.Portal>
		<ContextMenu.Content
			class="p-2 bg-white/95 backdrop-blur-xl border border-gray-200/50 rounded-xl shadow-2xl shadow-black/10 text-gray-800 min-w-[200px] animate-in fade-in-0 zoom-in-95 duration-200 ring-1 ring-black/5"
		>
			<ContextMenu.Item
				onclick={() => {
					removeItem();
				}}
				class="group flex items-center justify-between h-9 px-3 py-2 text-sm font-medium rounded-lg cursor-pointer select-none transition-all duration-150 ease-out hover:bg-red-50 focus:bg-red-50 focus:outline-none active:scale-[0.98]"
			>
				<div class="flex items-center">
					<Trash
						class="mr-3 size-4 text-gray-600 group-hover:text-red-600 transition-colors duration-150"
					/>
					<span class="group-hover:text-red-700 transition-colors duration-150">Delete</span>
				</div>
			</ContextMenu.Item>
		</ContextMenu.Content>
	</ContextMenu.Portal>
</ContextMenu.Root>
