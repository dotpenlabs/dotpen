<script lang="ts">
	import { flyAndScale } from '$/lib/utils';
	import { page } from '$app/state';
	import { ContextMenu } from 'bits-ui';
	import Trash from 'phosphor-svelte/lib/Trash';
	import TrashIcon from '$lib/components/movingicons/trash.svelte';

	interface NavItemProps {
		url: string;
		icon: typeof Trash;
		label: string;
		removeCollection?: () => void;
		onSelect?: (payload: { url: string; label: string }) => void;
	}

	let { url, icon: Icon, label, removeCollection = undefined, onSelect } = $props();
	let trashHovered = $state(false);

	function handleClick() {
		let accessed = JSON.parse(localStorage.getItem('track:accessed') || '{}');
		let currentTime = Date.now();

		for (const key in accessed) {
			if (currentTime - accessed[key].timestamp > 172800000) {
				delete accessed[key];
			}
		}

		// if (url === page.url.pathname) return;

		if (accessed[url]) {
			accessed[url].count += 1;
		} else {
			accessed[url] = { count: 1, timestamp: currentTime };
		}

		localStorage.setItem('track:accessed', JSON.stringify(accessed));
		onSelect?.({ url, label });
	}
</script>

{#if removeCollection}
	<ContextMenu.Root>
		<ContextMenu.Trigger>
			<a
				href={url}
				class="flex items-center relative gap-2 cursor-pointer w-full duration-200 hover:bg-stone-100 dark:hover:bg-stone-800 rounded-lg px-2 py-1"
				onclick={handleClick}
			>
				{#if page.url.pathname === url}
					<div
						transition:flyAndScale={{ x: 2.5, y: 0, duration: 150 }}
						class="bg-stone-300 dark:bg-stone-700 h-[65%] w-1 rounded-2xl right-0 mr-2 absolute"
					></div>
				{/if}
				<Icon class="size-4 opacity-80" />
				<p class="text-sm text-left">{label}</p>
			</a>
		</ContextMenu.Trigger>
		<ContextMenu.Portal>
			<ContextMenu.Content
				onmouseenter={() => {
					trashHovered = true;
				}}
				onmouseleave={() => {
					trashHovered = false;
				}}
				class="p-2 bg-white/95 backdrop-blur-xl border border-gray-200/50 rounded-xl shadow-2xl shadow-black/10 text-gray-800 min-w-[200px] animate-in fade-in-0 zoom-in-95 duration-200 ring-1 ring-black/5"
			>
				<ContextMenu.Item
					onclick={removeCollection}
					class="group flex items-center justify-between h-9 px-3 py-2 text-sm font-medium rounded-lg cursor-pointer select-none transition-all duration-150 ease-out hover:bg-red-50 focus:bg-red-50 focus:outline-none active:scale-[0.98]"
				>
					<div class="flex items-center gap-2">
						<TrashIcon isHovered={trashHovered} strokeWidth={3} size={12} class="opacity-80" />
						<span class="group-hover:text-red-700 transition-colors duration-150"> Delete </span>
					</div>
				</ContextMenu.Item>
			</ContextMenu.Content>
		</ContextMenu.Portal>
	</ContextMenu.Root>
{:else}
	<a
		href={url}
		class="flex items-center relative gap-2 cursor-pointer w-full duration-200 hover:bg-stone-100 dark:hover:bg-stone-800 rounded-lg px-2 py-1"
		onclick={handleClick}
	>
		{#if page.url.pathname === url}
			<div
				transition:flyAndScale={{ x: 2.5, y: 0, duration: 150 }}
				class="bg-stone-300 dark:bg-stone-700 h-[65%] w-1 rounded-2xl right-0 mr-2 absolute"
			></div>
		{/if}
		<Icon class="size-4 opacity-80" />
		<p class="text-sm text-left">{label}</p>
	</a>
{/if}
