<script>
	import { PluginKit } from '$/lib/plugins/head';
	import View from '$/lib/plugins/view.svelte';
	import { onMount } from 'svelte';

	onMount(() => {
		try {
			PluginKit.deregister('dropbox');

			PluginKit.register({
				id: 'dropbox',
				url: '/plugins/dropbox/page.html',
				location: 'sidebar-bottom',

				perms: ['app:sidebar-bottom']
			});
		} catch (e) {
			console.error(e);
		}
	});
</script>

<content class="h-full w-full absolute">
	<p>New plugins are now loaded!</p>
</content>
