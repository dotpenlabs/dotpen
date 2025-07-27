<script lang="ts">
	import { PluginKit } from '$/lib/plugins/client';
	import { onMount } from 'svelte';

	onMount(() => {
		PluginKit.Register('nl.bijsven.dropbox');
	});
</script>
