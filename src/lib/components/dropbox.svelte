<script lang="ts">
	import { onMount } from 'svelte';
	import { fade, fly, scale } from 'svelte/transition';
	import { cubicOut, elasticOut } from 'svelte/easing';
	import { HandGrabbing } from 'phosphor-svelte';
	import { toast } from 'svelte-sonner';

	let { ondrop } = $props<{
		ondrop?: (link: { url: string; label?: string }) => void;
	}>();

	let heldLinks = $state<{ url: string; label?: string }[]>([]);
	let isDragging = $state(false);
	let isVisible = $state(false);
	let draggedLink = $state<{ url: string; label?: string } | null>(null);
	let draggedOutside = $state(false);

	// Derived values using $derived rune
	let linkCount = $derived(heldLinks.length);

	// Event handlers
	function handleDrop(e: DragEvent) {
		e.preventDefault();
		isDragging = false;
		const url = e.dataTransfer?.getData('text/plain');
		if (url && !heldLinks.some((l) => l.url === url)) {
			const newLink = { url };
			heldLinks = [...heldLinks, newLink];
		}
	}

	function handlePaste(e: ClipboardEvent) {
		const url = e.clipboardData?.getData('text');
		if (url && !heldLinks.some((l) => l.url === url)) {
			const newLink = { url };
			heldLinks = [...heldLinks, newLink];
		}

		toast('Link moved to dropbox', {
			position: 'bottom-center',
			style: 'text-align: center; display: flex; justify-content: center;',
			description:
				'Drag & drop links out of the dropbox to organise them in the correct collection.'
		});
	}

	function handleDragOver(e: DragEvent) {
		e.preventDefault();
		isDragging = true;
	}

	function handleDragLeave() {
		isDragging = false;
	}

	function handleLinkDragStart(e: DragEvent, link: { url: string; label?: string }) {
		e.dataTransfer?.setData('text/plain', link.url);
		if (link.label) e.dataTransfer?.setData('text/label', link.label);
		draggedLink = link;
		draggedOutside = false;
	}

	function handleLinkDragEnd(e: DragEvent) {
		if (draggedLink && draggedOutside) {
			// Add to collection
			ondrop?.(draggedLink);
			// Remove from held links with animation
			heldLinks = heldLinks.filter((l) => l.url !== draggedLink!.url);
		}
		draggedLink = null;
		draggedOutside = false;
	}

	function handleGlobalDragOver(e: DragEvent) {
		if (draggedLink) {
			const dropZone = document.querySelector('.drop-zone');
			if (dropZone) {
				const rect = dropZone.getBoundingClientRect();
				const isInsideDropZone =
					e.clientX >= rect.left &&
					e.clientX <= rect.right &&
					e.clientY >= rect.top &&
					e.clientY <= rect.bottom;
				draggedOutside = !isInsideDropZone;
			}
		}
	}

	function getHostname(url: string) {
		try {
			return new URL(url).hostname;
		} catch {
			return url;
		}
	}

	function getFavicon(url: string) {
		try {
			const hostname = new URL(url).hostname;
			return `https://www.google.com/s2/favicons?domain=${hostname}&sz=32`;
		} catch {
			return null;
		}
	}

	$effect(() => {
		const timer = setTimeout(() => {
			isVisible = true;
		}, 100);

		return () => clearTimeout(timer);
	});

	onMount(() => {
		window.addEventListener('paste', handlePaste);
		window.addEventListener('dragover', handleGlobalDragOver);

		return () => {
			window.removeEventListener('paste', handlePaste);
			window.removeEventListener('dragover', handleGlobalDragOver);
		};
	});
</script>

<div class="flex justify-center pb-8 pr-1 duration-200">
	<div
		class={'drop-zone relative w-80 overflow-x-hidden rounded-3xl border border-dashed border-black/25 dark:border-white/25 transition-all duration-1000 ease-in-out overflow-hidden ' +
			(linkCount ? 'h-48' : 'h-28')}
		ondrop={handleDrop}
		ondragover={handleDragOver}
		ondragleave={handleDragLeave}
		role="button"
		tabindex="0"
	>
		<div class="p-4 h-full flex flex-col">
			{#if linkCount === 0}
				<div class="flex-1 flex flex-col justify-center items-center">
					<HandGrabbing class="size-6 opacity-45" />
					<p class="opacity-65 text-xs text-center w-[65%] mt-2">
						Drag & drop links, to organise them.
					</p>
				</div>
			{:else}
				<div class="flex-1 overflow-y-auto space-y-2 pr-1 overflow-x-hidden">
					{#each heldLinks as link, index (link.url)}
						<div
							role="button"
							tabindex="0"
							class="group relative bg-white/80 dark:bg-black/80 backdrop-blur-sm rounded-xl p-3 shadow-sm border border-gray-100/50 dark:border-stone-900/50 hover:shadow-md hover:border-gray-200/70 dark:hover:border-gray-700/70 hover:bg-white/90 dark:hover:bg-black/90 darktransition-all duration-300 cursor-move"
							draggable="true"
							ondragstart={(e) => handleLinkDragStart(e, link)}
							ondragend={handleLinkDragEnd}
							in:fly={{ x: 20, duration: 400, delay: index * 100, easing: cubicOut }}
							out:scale={{ duration: 300, easing: cubicOut, start: 0.8 }}
						>
							<div class="flex items-center gap-3">
								<div
									class="flex-shrink-0 w-6 h-6 rounded-md bg-gray-100 dark:bg-stone-900 flex items-center justify-center overflow-hidden"
								>
									{#if getFavicon(link.url)}
										<img src={getFavicon(link.url)} alt="" class="w-4 h-4 object-contain" />
									{:else}
										<div class="w-2 h-2 bg-gray-400 dark:bg-stone-700 rounded-full"></div>
									{/if}
								</div>

								<div class="flex-1 min-w-0">
									<p
										class="text-sm font-medium text-gray-900 dark:text-gray-100 truncate leading-tight"
									>
										{link.label || getHostname(link.url)}
									</p>
									<p
										class="text-xs text-gray-500 dark:text-stone-500 truncate mt-0.5 leading-tight"
									>
										{link.url}
									</p>
								</div>

								<div
									class="flex-shrink-0 opacity-0 group-hover:opacity-100 transition-opacity duration-200"
								>
									<div
										class="w-1 h-8 bg-gray-300 dark:bg-stone-700 rounded-full flex items-center justify-center"
									>
										<div class="w-0.5 h-4 bg-gray-400 dark:bg-stone-600 rounded-full"></div>
									</div>
								</div>
							</div>
						</div>
					{/each}
				</div>
			{/if}
		</div>
	</div>
</div>
