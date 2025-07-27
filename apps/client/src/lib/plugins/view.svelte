<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import Host from './host.svelte';
	import { _PluginKit, type PluginConfig } from './head';

	let { location } = $props();

	let plugins = $state([]) as PluginConfig[];

	const updatePlugins = () => {
		plugins = _PluginKit.getByLocation(location);
	};

	onMount(() => {
		updatePlugins();
		const unsubscribe = _PluginKit.subscribe(updatePlugins);

		onDestroy(() => {
			unsubscribe();
		});
	});
</script>

{#each plugins as plugin}
	<Host id={plugin.id} name={plugin.name} srcdoc={plugin.content} />
{/each}
