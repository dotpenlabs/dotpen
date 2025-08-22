<script lang="ts">
	import { page } from '$app/state';
	import { BugBeetle } from 'phosphor-svelte';
</script>

<svelte:head>
    <title>Dotpen</title>
</svelte:head>

<content
	class="absolute top-0 left-0 inset-0 flex flex-col justify-center items-center opacity-65 gap-2"
>
	<div class="flex flex-row items-center gap-1">
		<BugBeetle color="red" />
		<h1 class="!font-sans">{page.status} {page.error.message}</h1>
	</div>
	<button
		class="opacity-45 cursor-pointer underline text-xs hover:opacity-100 transition-opacity duration-300 ease-in-out"
		onclick={() => {
			window.history.back();
		}}
	>
		Go back
	</button>
</content>
