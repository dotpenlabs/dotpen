<script lang="ts">
	import SvelteMarkdown from 'svelte-markdown';
	import privacy from './privacy.md?raw';
</script>

<svelte:head>
	<title>Privacy Policy • Dotpen</title>
</svelte:head>

<content class="flex justify-center w-full h-full absolute overflow-y-auto">
	<div class="max-w-[56rem] w-full pt-24 font-light">
		<p class="!text-3xl font-semibold">Privacy Policy</p>
		<div class="w-16 h-1 rounded-xl bg-black/25 my-3.5 mb-8"></div>
		<div id="markdown" class="opacity-95 pb-24">
			<SvelteMarkdown source={privacy} />
		</div>
	</div>
</content>