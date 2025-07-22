<script lang="ts">
	import { onMount, afterUpdate, type Snippet } from 'svelte';

	export let gridGap = '0.5em';
	export let colWidth = 'minmax(Min(20em, 100%), 1fr)';
	export let stretchFirst = false;
	export let children: Snippet;

	let container: HTMLElement;
	let childElements: HTMLElement[] = [];

	function applyMasonry() {
		if (!container) return;

		// Get all direct children
		childElements = Array.from(container.children) as HTMLElement[];
		if (childElements.length === 0) return;

		// Reset all margins
		childElements.forEach((child) => {
			child.style.marginTop = '0';
		});

		// Get the number of columns
		const computedStyle = getComputedStyle(container);
		const columnCount = computedStyle.gridTemplateColumns.split(' ').length;

		if (columnCount <= 1) return; // No need for masonry with just one column

		// Skip the first row
		for (let i = columnCount; i < childElements.length; i++) {
			const child = childElements[i];
			const itemAbove = childElements[i - columnCount];

			// Calculate the gap
			const itemAboveRect = itemAbove.getBoundingClientRect();
			const childRect = child.getBoundingClientRect();

			// Apply margin to position items correctly
			const gap = parseInt(computedStyle.rowGap || computedStyle.gridRowGap || '0', 10);
			const marginTop = itemAboveRect.bottom + gap - childRect.top;

			if (marginTop > 0) {
				child.style.marginTop = `${marginTop}px`;
			}
		}
	}

	// Apply masonry layout when DOM updates
	afterUpdate(() => {
		applyMasonry();
	});

	// Set up resize observer and initial layout
	onMount(() => {
		applyMasonry();

		// Add resize observer
		const resizeObserver = new ResizeObserver(() => {
			applyMasonry();
		});

		resizeObserver.observe(container);

		// Add window resize listener
		window.addEventListener('resize', applyMasonry);

		return () => {
			resizeObserver.disconnect();
			window.removeEventListener('resize', applyMasonry);
		};
	});
</script>

<div
	bind:this={container}
	class="masonry-grid {stretchFirst ? 'stretch-first' : ''}"
	style="--grid-gap: {gridGap}; --col-width: {colWidth};"
>
	{@render children()}
</div>

<style>
	.masonry-grid {
		display: grid;
		grid-template-columns: repeat(auto-fit, var(--col-width));
		grid-gap: var(--grid-gap);
		width: 100%;
	}

	/* For browsers that support native masonry */
	@supports (grid-template-rows: masonry) {
		.masonry-grid {
			grid-template-rows: masonry;
		}
	}

	/* Make sure all direct children align to start */
	.masonry-grid > :global(*) {
		align-self: start;
	}

	/* First item stretches across all columns if needed */
	.masonry-grid.stretch-first > :global(*:first-child) {
		grid-column: 1 / -1;
	}
</style>
