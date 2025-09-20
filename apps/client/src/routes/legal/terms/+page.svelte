<script lang="ts">
	import terms from './terms.html?raw';
</script>

<svelte:head>
	<title>Terms of Service • Dotpen</title>
</svelte:head>

<content class="flex justify-center w-full h-full absolute overflow-y-auto">
	<div class="max-w-[56rem] w-full pt-24 font-light">
		<p class="!text-3xl font-semibold">Terms of Service</p>
		<div class="w-16 h-1 rounded-xl bg-black/25 my-3.5 mb-8"></div>
		<div id="markdown" class="opacity-95 pb-24">
			{@html terms}
		</div>
	</div>
</content>
