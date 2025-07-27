<script lang="ts">
	import { onMount } from 'svelte';
	import { _PluginKit } from './head';

	let { id, name, srcdoc } = $props();

	let renderError = $state(false);
	let render = $state() as HTMLIFrameElement;

	const sandbox_default = 'allow-scripts';

	onMount(() => {
		_PluginKit.bind(id, render);
	});
</script>

{#if renderError}
	<div class="text-center text-red-500 mt-8">
		<p>Dotpen was unable to load this plugin. Check console.</p>
	</div>
{:else}
	<iframe
		bind:this={render}
		{srcdoc}
		sandbox={sandbox_default}
		data-plugin-id={id}
		style="border: none; width: 100%; height: 100%;"
		loading="lazy"
		title={'Plugin: ' + name}
	></iframe>
{/if}
