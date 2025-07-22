<script lang="ts">
	const { children } = $props();
</script>

{@render children()}

<style lang="postcss">
	@reference "tailwindcss";

	:global(#markdown h1) {
		@apply text-xl font-bold mt-4;
	}

	:global(#markdown h2) {
		@apply text-2xl font-bold;
	}

	:global(#markdown h3) {
		@apply text-xl font-bold;
	}

	:global(#markdown p) {
		@apply text-base mt-1;
	}

	:global(#markdown ul) {
		@apply list-disc;
	}
</style>
