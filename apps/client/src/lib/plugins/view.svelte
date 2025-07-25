<script lang="ts">
	import { onMount } from 'svelte';
	import Host from './host.svelte';
	import { plgmgr } from './mgr';

	let { location } = $props();

	let plugins = $state([]);

	onMount(() => {
		plugins = plgmgr.getPluginsAt(location);
	});
</script>

<div>
	{#each plugins as plugin}
		<p>{plugin.name}</p>
		<Host name={plugin.name} src={plugin.url} />
	{/each}
</div>
