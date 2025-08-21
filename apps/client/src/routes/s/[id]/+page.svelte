<script lang="ts">
	import { onMount, tick } from 'svelte';
	import { pb } from '$/lib';
	import { page } from '$app/state';
	import { toast } from 'svelte-sonner';
	import { browser } from '$app/environment';

	import Loading from '$/lib/components/navigation/loading.svelte';
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

	let appstate = $state('loading') as 'loading' | 'show' | 'error';
	let snapshot = $state(null) as Snapshot | null;

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
			} catch (e) {
				console.warn('Error removing masonry', e);
			}
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

	onMount(async () => {
		try {
			const pid = page.url.pathname.split('/').at(-1) || '';
			console.log('share pid', pid);
			try {
				snapshot = (await pb.collection('snapshots').getFirstListItem(`slug = "${pid}"`, {
					query: {
						slug: pid
					}
				})) as unknown as Snapshot;
			} catch (e) {
				console.error('Failed to load snapshot', e);
				appstate = 'error';
			}

			appstate = snapshot ? 'show' : 'error';

			if (snapshot) {
				await UseMasonry();
			}
		} catch (e) {
			console.error('Failed to load snapshot', e);
			appstate = 'error';
		}
	});
</script>

<svelte:head>
	<title>Snapshot: {snapshot?.name || 'Snapshot'} • Dotpen</title>
</svelte:head>

{#if appstate === 'loading'}
	<div class="h-full w-full flex items-center justify-center"><Loading /></div>
{:else if appstate === 'error'}
	<div class="absolute h-full w-full flex gap-2 items-center justify-center opacity-65">
		<Ghost weight="bold" class="size-5 opacity-80 text-black dark:text-white" />
		<p>Snapshot does not exist.</p>
	</div>
{:else}
	<div
		class="z-[999] fixed bottom-6 left-1/2 transform -translate-x-1/2 flex flex-col sm:flex-row items-center justify-between w-full max-w-3xl h-auto bg-black/65 backdrop-blur-xl rounded-2xl border border-white/10 px-4 sm:px-6 py-4 sm:py-6 shadow-lg"
	>
		<div class="flex items-center gap-4 sm:gap-6 w-full sm:w-auto">
			<div class="-space-y-1">
				<div class="flex items-center gap-2">
					<h3 class="text-white text-lg font-medium !font-sans tracking-tight">
						Shared collections
					</h3>
				</div>
				<p class="text-white/60 !font-sans text-xs sm:w-96 font-normal leading-5">
					Create collections in your own style and instantly share them with friends, teammates, or
					anyone you like.<br /><br />
					Shared on {new Date(snapshot?.created).toLocaleDateString()} at {new Date(
						snapshot?.created
					).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}.
				</p>
			</div>
		</div>

		<div class="flex items-center gap-2 sm:gap-4 w-full sm:w-auto mt-4 sm:mt-0 justify-end">
			<a
				href="mailto:report@bijsven.nl"
				title="Report content"
				class="text-white/30 text-xs hover:text-orange-500 transition-colors duration-200 !font-sans text-nowrap"
			>
				<Warning weight="bold" class="w-4 h-4 mr-1" />
			</a>
			<a
				href="/"
				class="inline-flex items-center justify-center bg-white text-gray-900 text-sm font-medium px-4 py-2 rounded-xl hover:bg-white/90 transition-all duration-200 !font-sans whitespace-nowrap"
			>
				{#if !pb.authStore.isValid}
					Try Dotpen
				{:else}
					Open Dotpen
				{/if}
			</a>
		</div>
	</div>
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
