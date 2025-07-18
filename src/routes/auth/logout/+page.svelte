<script lang="ts">
	import { pb } from '$/lib';
	import Loading from '$/lib/components/loading.svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import { onMount } from 'svelte';
	import { toast } from 'svelte-sonner';
	import { clearAll as idbClearAll } from '$/lib/idb';

	onMount(() => {
		console.log('[DUM] Logging out...');

		pb.authStore.clear();
		idbClearAll().then(() => {
			toast('Byee!');
			goto('/');
		});
	});
</script>

<svelte:head>
	<title>Dotpen User Management</title>
</svelte:head>

<content class="absolute h-full w-full">
	<Loading />
</content>
