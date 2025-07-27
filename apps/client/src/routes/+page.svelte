<script lang="ts">
	import { pb } from '$/lib';
	import Loading from '$/lib/components/loading.svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import { onMount } from 'svelte';
	import { getItem as idbGetItem, setItem as idbSetItem } from '$/lib/idb';

	onMount(async () => {
		if (!pb.authStore.isValid) {
			goto('/auth');
		} else {
			const userId = pb.authStore.record?.id;
			if (userId) {
				const collections = await pb.collection('collections').getFullList({
					filter: `user = "${userId}" && name != "system_inbox"`,
					sort: 'created'
				});
				let needsCache = false;
				for (const col of collections) {
					const cache = await idbGetItem('collection:' + col.id + ':cache');
					if (!cache) {
						needsCache = true;
						const bookmarks = await pb.collection('bookmarks').getFullList({
							filter: `collection = "${col.id}" && deleted = false`,
							sort: '-created'
						});
						await idbSetItem('collection:' + col.id + ':cache', JSON.stringify(bookmarks));
					}
				}

				try {
					const inboxId = await pb
						.collection('collections')
						.getFirstListItem("name = 'system_inbox'")
						.then((e) => e.id);
					const inboxCache = await idbGetItem('collection:' + inboxId + ':cache');
					if (!inboxCache) {
						needsCache = true;
						const inboxBookmarks = await pb.collection('bookmarks').getFullList({
							filter: `collection = "${inboxId}" && deleted = false`,
							sort: '-created'
						});
						await idbSetItem('collection:' + inboxId + ':cache', JSON.stringify(inboxBookmarks));
						await idbSetItem('inbox:id', inboxId);
					}
					if (needsCache) {
						location.reload();
						return;
					}
				} catch (e) {
					console.info('[head] No inbox found, skipping cache...');
				}
			}
			goto('/inbox');
		}
	});
</script>

<svelte:head>
	<title>Dotpen</title>
</svelte:head>

<content class="absolute h-full w-full">
	<Loading />
</content>
