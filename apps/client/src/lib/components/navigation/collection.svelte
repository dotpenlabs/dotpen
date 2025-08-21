<script lang="ts">
	import { flyAndScale } from '$/lib/utils';
	import { page } from '$app/state';
	import { ContextMenu } from 'bits-ui';
	import Trash from 'phosphor-svelte/lib/Trash';
	import TrashIcon from '$/lib/components/icons/trash.svelte';
	import { Folder } from 'phosphor-svelte';

	interface NavItemProps {
		id: string;
		icon: typeof Folder;
		label: string;
		removeCollection?: () => void;
		onSelect?: (payload: { id: string; label: string }) => void;
		shareSnapshot?: () => void;
	}

	let { id, label, icon: Icon, removeCollection, onSelect, shareSnapshot }: NavItemProps = $props();
	let trashHovered = $state(false);

	function handleClick() {
		let accessed = JSON.parse(localStorage.getItem('track:accessed') || '{}');
		let currentTime = Date.now();

		for (const key in accessed) {
			if (currentTime - accessed[key].timestamp > 172800000) {
				delete accessed[key];
			}
		}

		if (accessed[id]) {
			accessed[id].count += 1;
		} else {
			accessed[id] = { count: 1, timestamp: currentTime };
		}

		localStorage.setItem('track:accessed', JSON.stringify(accessed));
		onSelect?.({ id, label });
	}
</script>

<ContextMenu.Root>
	<ContextMenu.Trigger>
		<a
			href={id}
			class="flex h-7 items-center relative gap-2 cursor-pointer w-full duration-200 hover:bg-stone-100 dark:hover:bg-stone-800 rounded-lg p-1 pl-2"
			onclick={handleClick}
		>
			{#if page.url.pathname === `/${id}`}
				<div
					transition:flyAndScale={{ x: 2.5, y: 0, duration: 150 }}
					class="bg-stone-300 dark:bg-stone-700 h-[65%] w-1 rounded-2xl right-0 mr-2 absolute"
				></div>
			{/if}
			<Icon class="min-h-4 min-w-4 size-4 opacity-80" />
			<p class="text-sm text-left truncate pr-12">{label}</p>
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
			class="p-2 bg-white/95 dark:bg-stone-950 backdrop-blur-xl border border-gray-200/50 dark:border-stone-700/50 rounded-xl shadow-2xl shadow-black/10 dark:shadow-stone-950/10 text-gray-800 dark:text-stone-200 min-w-[200px] animate-in fade-in-0 zoom-in-95 duration-200 ring-1 ring-black/5 dark:ring-white/5"
		>
			<ContextMenu.Item
				onclick={() => shareSnapshot?.()}
				class="group flex items-center justify-between h-9 px-3 py-2 text-sm font-medium rounded-lg cursor-pointer select-none transition-all duration-150 ease-out hover:bg-stone-50 dark:hover:bg-stone-900 focus:outline-none active:scale-[0.98]"
			>
				<div class="flex items-center gap-2">
					<span class="opacity-80">🔗</span>
					<span
						class="group-hover:text-stone-900 group-hover:dark:text-stone-100 transition-colors duration-150"
					>
						Share snapshot
					</span>
				</div>
			</ContextMenu.Item>
			<ContextMenu.Item
				onclick={removeCollection}
				class="group flex items-center justify-between h-9 px-3 py-2 text-sm font-medium rounded-lg cursor-pointer select-none transition-all duration-150 ease-out hover:bg-red-50 dark:hover:bg-red-950 focus:bg-red-50 focus:outline-none active:scale-[0.98]"
			>
				<div class="flex items-center gap-2">
					<TrashIcon isHovered={trashHovered} strokeWidth={3} size={12} class="opacity-80" />
					<span class="group-hover:text-red-700 transition-colors duration-150"> Delete </span>
				</div>
			</ContextMenu.Item>
		</ContextMenu.Content>
	</ContextMenu.Portal>
</ContextMenu.Root>
