<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import Host from './host.svelte';
	import { PluginKit } from './head';

	let { location } = $props();

	let plugins = $state([]);

	const updatePlugins = () => {
		plugins = PluginKit.getByLocation(location);
	};

	onMount(() => {
		updatePlugins();
		const unsubscribe = PluginKit.subscribe(updatePlugins);

		onDestroy(() => {
			unsubscribe();
		});
	});
</script>

{#each plugins as plugin}
	<Host id={plugin.id} name={plugin.name} src={plugin.url} />
{/each}
