<script lang="ts">
	import { onMount } from 'svelte';

	import { Reveal, Rotate } from '../utils';

	let loading = $state(true);

	onMount(() => {
		loading = false;
	});
</script>

{#if !loading}
	<div class="flex items-center gap-1 w-auto group">
		<img
			in:Rotate={{ delay: 150, duration: 5000 }}
			src="/logo.svg"
			alt="Cognitum Logo"
			class="size-5 invert-100 group-hover:rotate-[360deg] group-active:rotate-[720deg] group-hover:scale-110 dark:invert-0 hover:scale-110 active:scale-90 transition-all duration-1000 dark:opacity-90"
		/>
		<span in:Reveal={{ duration: 800 }} class="font-medium text-2xl ml-2">Cognitum</span>
	</div>
{/if}
