<script lang="ts">
	import { flyAndScale } from '$/lib/utils';
	import { page } from '$app/state';

	let { url, icon: Icon, label } = $props();
	function handleClick() {
		let accessed = JSON.parse(localStorage.getItem('track:accessed') || '{}');
		let currentTime = Date.now();

		for (const key in accessed) {
			if (currentTime - accessed[key].timestamp > 172800000) {
				delete accessed[key];
			}
		}

		if (url === page.url.pathname) return;

		if (accessed[url]) {
			accessed[url].count += 1;
		} else {
			accessed[url] = { count: 1, timestamp: currentTime };
		}

		localStorage.setItem('track:accessed', JSON.stringify(accessed));
	}
</script>

<a
	href={url}
	class="flex items-center relative gap-2 cursor-pointer w-full duration-200 hover:bg-stone-100 dark:hover:bg-stone-800 rounded-lg px-2 py-1"
	onclick={handleClick}
>
	{#if page.url.pathname === url}
		<div
			transition:flyAndScale={{ x: 2.5, y: 0, duration: 150 }}
			class="bg-stone-300 dark:bg-stone-700 h-[65%] w-1 rounded-2xl right-0 mr-2 absolute"
		></div>
	{/if}
	<Icon class="size-4 opacity-80" />
	<p class="text-sm text-left">{label}</p>
</a>
