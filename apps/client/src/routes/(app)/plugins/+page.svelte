<script>
	import { plgmgr } from '$/lib/plugins/mgr';
	import View from '$/lib/plugins/view.svelte';
	import { onMount } from 'svelte';

	plgmgr.register({
		name: 'main',
		url: '/plugins/simple/page.html',
		location: 'main'
	});
</script>

<content class="bg-red-500 h-full w-full absolute">
	<View location="main" />
</content>
