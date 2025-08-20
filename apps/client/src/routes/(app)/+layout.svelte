<script lang="ts">
	import { onMount, type Snippet } from 'svelte';
	import { Folder } from 'phosphor-svelte';

	import { goto } from '$app/navigation';
	import { pb, pk } from '$/lib';
	import { page } from '$app/state';

	import Loading from '$/lib/components/navigation/loading.svelte';
	import type { Collection } from '$/lib/types';
	import { navigating } from '$app/state';
	import Colorwrap from './colorwrap.svelte';
	import Sidebar from './sidebar.svelte';

	let { children }: { children: Snippet } = $props();

	let appstate = $state('app_loading') as 'hydrating' | 'app_loading' | 'app' | 'offline';
	let appplatform = $state('desktop') as 'desktop' | 'mobile';
	let appcolor = $state('default') as 'default' | 'green' | 'red';

	let sidebarloc = $state('left') as 'left' | 'right';

	let collections = $state([]) as Collection[];

	const KEYS = {
		COLCACHE: 'layout:collections:cache'
	};

	onMount(async () => {
		// console.log('[/app/auth] checking authenication');
		// if (!pb.authStore.isValid) {
		// 	console.warn('[/app/auth] not authenticated, redirecting to /auth');
		// 	goto('/auth');
		// 	return;
		// } else {
		// 	console.log(`[/app/auth] Welcome back, ${pb.authStore.record.name}!`);
		// }

		// console.log('[/app/colcache] checking collection-cache');
		// const cached = localStorage.getItem(KEYS.COLCACHE);
		// let ucache = false;
		// if (cached) {
		// 	const d = JSON.parse(cached);
		// 	if (Array.isArray(d)) {
		// 		collections = d as Collection[];
		// 	}
		// 	ucache = true;
		// 	console.log('[/app/colcache] cache found, loading collections');
		// } else {
		// 	console.log('[/app/colcache] cache not found, starting bg fetch');
		// 	const result = await pb
		// 		.collection('collections')
		// 		.getFullList({
		// 			filter: `user = "${pb.authStore.model?.id}" && name != "system_inbox"`,
		// 			sort: 'created'
		// 		})
		// 		.then((res) => {
		// 			return res.map((col) => ({
		// 				id: col.id,
		// 				name: col.name,
		// 				icon: Folder,
		// 				url: `/${col.id}`
		// 			}));
		// 		});
		// 	collections = result as Collection[];

		// 	console.log('[/app/colcache] saving collection-cache');
		// 	localStorage.setItem(KEYS.COLCACHE, JSON.stringify(collections));
		// 	console.log('[/app/colcache] collection-cache saved');
		// 	console.log('Collections loaded:', collections);
		// }

		// (async () => {
		// 	if (ucache) {
		// 		console.log('[/app/colcache] updating cache');
		// 		const result = await pb
		// 			.collection('collections')
		// 			.getFullList({
		// 				filter: `user = "${pb.authStore.model?.id}" && name != "system_inbox"`,
		// 				sort: 'created'
		// 			})
		// 			.then((res) => {
		// 				return res.map((col) => ({
		// 					id: col.id,
		// 					name: col.name,
		// 					icon: Folder,
		// 					url: `/${col.id}`
		// 				}));
		// 			});
		// 		collections = result as Collection[];

		// 		console.log('[/app/colcache] saving collection-cache');
		// 		localStorage.setItem(KEYS.COLCACHE, JSON.stringify(collections));

		// 		console.log('Collections loaded:', collections);
		// 	}
		// })();

		// window.SetHydrating = (value: boolean) => {
		// 	appstate = value ? 'hydrating' : appstate === 'hydrating' ? 'app' : appstate;
		// };

		// window.SetColored = (isColored: 'default' | 'green' | 'red') => {
		// 	appcolor = isColored;
		// 	setTimeout(() => {
		// 		appcolor = 'default';
		// 	}, 1500);
		// };

		// setTitle();

		// appstate = 'app';

		// appplatform = (() => {
		// 	// use md (in tailwind) for the correct width breakpoint.
		// 	if (window.innerWidth < 768) {
		// 		return 'mobile';
		// 	} else {
		// 		return 'desktop';
		// 	}
		// })();

		await pk.init();
		await pk.add('nl.bijsven.dropbox');
	});

	$effect(() => {
		if (navigating.complete) {
			console.log('[app/effect] navigating.complete, setting title');
			setTitle();
		}
	});

	function setTitle() {
		const currentPath = page.url.pathname;
		const currentCollection = collections.find((c) => c.id === currentPath);
		document.title = currentCollection ? `${currentCollection.name} • Dotpen` : 'Dotpen';

		switch (currentPath) {
			case '/inbox':
				document.title = 'Inbox • Dotpen';
				break;
			case '/plugins':
				document.title = 'Plugins • Dotpen';
				break;
		}
	}
</script>

{#if appstate === 'app_loading'}
	<Loading />
{:else}
	<pk-widget location="sidebar:bottom"></pk-widget>
	<!-- <Colorwrap {appplatform} {appcolor}>
		<content
			data-vaul-drawer-wrapper
			class="h-full w-full flex bg-white dark:bg-stone-950 overflow-y-auto rounded-xl border border-stone-200 dark:border-stone-800"
		>
			{#if sidebarloc === 'left'}
				<Sidebar bind:collections {appstate} {appplatform} {sidebarloc} />
				<div class="h-full w-px bg-stone-200 dark:bg-stone-800 mr-5"></div>
			{/if}
			<div class="inset-0 w-full h-full">
				{@render children()}
			</div>
			{#if sidebarloc === 'right'}
				<div class="h-full w-px bg-stone-200 dark:bg-stone-800 ml-5"></div>
				<Sidebar bind:collections {appstate} {appplatform} {sidebarloc} />
			{/if}
		</content>
	</Colorwrap> -->
{/if}
