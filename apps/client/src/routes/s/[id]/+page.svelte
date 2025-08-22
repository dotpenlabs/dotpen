<script lang="ts">
	import { onMount, tick } from 'svelte';
	import { pb } from '$/lib';
	import { page } from '$app/state';
	import { toast } from 'svelte-sonner';
	import { browser } from '$app/environment';

	import PublicLink from '$/lib/components/item/public_link.svelte';
	import { Ghost, Warning } from 'phosphor-svelte';

	let Macy: any;
	let masonry: any | undefined;

	type Snapshot = {
		id: string;
		slug?: string;
		name?: string;
		created: string;
		data: {
			id?: string;
			label: string;
			link: string;
			created?: string;
			_cover_base64?: string;
			_favicon_base64?: string;
		}[];
	};

	let appstate = $state<'loading' | 'show' | 'error'>('loading');
	let snapshot = $state<Snapshot | null>(null);
	let progress = $state(0);

	async function UseMasonry() {
		await tick();
		if (!Macy && browser) {
			const macyModule = await import('macy');
			Macy = macyModule.default;
		}

		const container = document.querySelector('#snapshot-content');
		if (!container) return;

		if (masonry) {
			try {
				masonry.remove();
			} catch {}
		}

		masonry = new Macy({
			container,
			trueOrder: false,
			waitForImages: true,
			useOwnImageLoader: false,
			useContainerForBreakpoints: true,
			mobileFirst: true,
			columns: 1,
			margin: { y: 16, x: '2%' },
			breakAt: {
				6000: 18,
				5500: 17,
				5000: 16,
				4500: 15,
				4200: 14,
				4000: 13,
				3500: 12,
				3000: 11,
				2500: 10,
				2200: 9,
				2000: 8,
				1800: 6,
				1600: 6,
				1450: 6,
				1300: 5,
				1200: 5,
				1000: 4,
				900: 3,
				800: 3,
				700: 3,
				600: 2,
				525: 2,
				500: 1
			}
		});

		requestAnimationFrame(() => masonry.recalculate());
	}

	async function fetchSnapshot(slug: string) {
		appstate = 'loading';
		progress = 0;

		try {
			snapshot = await pb.collection('snapshots').getFirstListItem(`slug = "${slug}"`, {
				query: {
					slug
				},
				fetch: async (url, config) => {
					const response = await fetch(url, config);

					if (!response.body) {
						throw new Error('No response body');
					}

					// Reader gebruiken om progress te meten
					const reader = response.body.getReader();
					const contentLength = +response.headers.get('Content-Length')!;
					let receivedLength = 0;
					let chunks: Uint8Array[] = [];

					while (true) {
						const { done, value } = await reader.read();
						if (done) break;

						chunks.push(value);
						receivedLength += value.length;
						if (contentLength) {
							progress = Math.round((receivedLength / contentLength) * 100);
						} else {
							// fallback als Content-Length ontbreekt
							progress = Math.min(progress + 5, 95);
						}
					}

					// Alles samenvoegen
					const chunksAll = new Uint8Array(receivedLength);
					let position = 0;
					for (let chunk of chunks) {
						chunksAll.set(chunk, position);
						position += chunk.length;
					}

					return new Response(chunksAll, {
						status: response.status,
						statusText: response.statusText,
						headers: response.headers
					});
				}
			});

			appstate = snapshot ? 'show' : 'error';

			if (snapshot) {
				await UseMasonry();
			}
		} catch (e) {
			console.error('Failed to load snapshot', e);
			appstate = 'error';
		}
	}

	onMount(async () => {
		const pid = page.url.pathname.split('/').at(-1) || '';
		await fetchSnapshot(pid);
	});
</script>

<svelte:head>
	<title>Snapshot: {snapshot?.name || 'Snapshot'} • Dotpen</title>
</svelte:head>

{#if appstate === 'loading'}
	<!-- Loader met percentage -->
	<div
		class="fixed inset-0 z-50 flex flex-col gap-4 items-center justify-center bg-white dark:bg-black"
	>
		<img src="/full.png" alt="Dotpen Logo" class="h-8 opacity-90 invert-100 dark:invert-0" />
		<div class="w-64 h-1.5 bg-black/25 dark:bg-white/25 rounded-full overflow-hidden relative">
			{#if progress === 0}
				<div
					class="absolute h-full bg-black/50 dark:bg-white/50 rounded-full animate-loadingBar-alternate"
				></div>
				<div
					class="absolute h-full bg-black/50 dark:bg-white/50 rounded-full animate-loadingBar"
				></div>
			{:else}
				<div
					class="absolute h-full bg-black/50 dark:bg-white/50 rounded-full left-0 transition-all duration-200"
					style={`width: ${progress}%;`}
				></div>
			{/if}
		</div>
	</div>
{:else if appstate === 'error'}
	<div class="absolute h-full w-full flex gap-2 items-center justify-center opacity-65">
		<Ghost weight="bold" class="size-5 opacity-80 text-black dark:text-white" />
		<p>Snapshot does not exist.</p>
	</div>
{:else}
	<div class="min-h-full w-full flex flex-col pt-6 gap-6 px-12 pb-44">
		<div id="snapshot-content" class="w-full px-4">
			{#each snapshot?.data || [] as item, idx (item.id || idx)}
				<div>
					<PublicLink data={item} index={idx} />
				</div>
			{/each}
		</div>
	</div>
{/if}

<style>
	@keyframes loadingBar {
		0% {
			left: -10%;
			width: 5%;
		}
		50% {
			width: 65%;
		}
		100% {
			left: 105%;
			width: 5%;
		}
	}
	.animate-loadingBar-alternate {
		animation: loadingBar 1s ease-in-out infinite alternate;
	}
	.animate-loadingBar {
		animation: loadingBar 1.25s ease-in-out infinite alternate;
		animation-delay: 0.5s;
	}
	.animate-loadingBar,
	.animate-loadingBar-alternate {
		will-change: left, width;
	}
</style>
