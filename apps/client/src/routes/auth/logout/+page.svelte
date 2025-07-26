<script lang="ts">
	import { pb } from '$/lib';
	import Loading from '$/lib/components/loading.svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import { onMount } from 'svelte';
	import { toast } from 'svelte-sonner';
	import { clearAll as idbClearAll } from '$/lib/idb';

	onMount(() => {
		console.log('[System] Logging out...');

		localStorage.clear();
		sessionStorage.clear();

		document.cookie.split(';').forEach(function (c) {
			document.cookie = c.trim().split('=')[0] + '=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/';
		});

		if (indexedDB.databases) {
			indexedDB.databases().then((dbs) => {
				dbs.forEach((db) => indexedDB.deleteDatabase(db.name));
			});
		}

		goto('/');
	});
</script>

<svelte:head>
	<title>Dotpen User Management</title>
</svelte:head>

<content class="absolute h-full w-full">
	<Loading />
</content>
