<script lang="ts">
	import { fade } from 'svelte/transition';
	import LinkComponent from '$/lib/components/item/link.svelte';
	import type { LinkItem } from '$/lib/components/item/types';
	import { toast } from 'svelte-sonner';
	import { onMount } from 'svelte';
	import Loading from '$/lib/components/loading.svelte';
	import { page } from '$app/stores';
	import type { SupabaseClient, User } from '@supabase/supabase-js';

	let loading = true;
	let Marks: LinkItem[] = [];

	onMount(async () => {
		// Haal de supabase client en de user op uit de $page store
		const { supabase, user } = $page.data as unknown as { supabase: SupabaseClient; user: User };

		try {
			// Query de database om de opgeslagen marks/decks op te halen van de ingelogde user
			const { data, error } = await supabase
				.from('marks')
				.select('*')
				.eq('user_id', user?.id)
				.order('date', { ascending: false });

			if (error) {
				console.error('Error fetching marks:', error);
				toast.error('Kon de opgeslagen items niet ophalen');
			} else if (data) {
				Marks = data as LinkItem[];
			}
		} catch (err) {
			console.error('Unexpected error:', err);
			toast.error('Er is een onverwachte fout opgetreden bij het ophalen van de opgeslagen items.');
		} finally {
			loading = false;
		}
	});

	async function addMark(userData: any) {
		if (!userData) return;

		const { supabase, user } = $page.data as unknown as { supabase: SupabaseClient; user: User };

		const formattedUrl =
			userData.startsWith('http://') || userData.startsWith('https://')
				? userData
				: 'https://' + userData;
		const urlRegex =
			/(https?:\/\/(?:www\.|(?!www))[a-zA-Z0-9][a-zA-Z0-9-]+[a-zA-Z0-9]\.[^\s]{2,}|www\.[a-zA-Z0-9][a-zA-Z0-9-]+[a-zA-Z0-9]\.[^\s]{2,}|https?:\/\/(?:www\.|(?!www))[a-zA-Z0-9]+\.[^\s]{2,}|www\.[a-zA-Z0-9]+\.[^\s]{2,})/;
		const isValidUrl = urlRegex.test(formattedUrl);

		if (isValidUrl) {
			const loading_toast = toast.loading('Capturing link', {
				description: 'This may take a few seconds...'
			});

			try {
				const res = await fetch('/system/crawl?url=' + formattedUrl);
				const data = await res.json();
				if (!data.title || !data.image) {
					toast.warning('Some fields are not resolved.', {
						description: 'Learn why on our help page.'
					});
				}

				const checkFavicon = async (url: string) => {
					const img = new Image();
					img.src = url;
					return new Promise<string>((resolve) => {
						img.onload = () => resolve(url);
						img.onerror = () => resolve('');
					});
				};

				const faviconUrl = await checkFavicon(new URL(formattedUrl).origin + '/favicon.ico').then(
					async (icon) => {
						if (icon) return icon;
						return await checkFavicon(new URL(formattedUrl).origin + '/favicon.png');
					}
				);

				const newMark: LinkItem = {
					label: data.title || 'Untitled',
					icon: faviconUrl,
					url: formattedUrl,
					cover:
						data.image && data.image.startsWith('/')
							? new URL(formattedUrl).origin + data.image
							: data.image,
					date: new Date().toISOString()
				};

				// Voeg de nieuwe mark toe aan de lijst
				Marks = [newMark, ...Marks];

				const { error } = await supabase.from('marks').insert({
					user_id: user?.id,
					label: newMark.label,
					icon: newMark.icon,
					url: newMark.url,
					cover: newMark.cover,
					date: newMark.date
				});

				if (error) {
					console.error('DB insert error:', error);
					toast.error('Error saving link in database', {
						description: 'Er is een fout opgetreden bij het opslaan in de database.'
					});
				} else {
					toast.success('Link captured and saved!', {
						description:
							'The link has been successfully added to your inbox and saved in the database!'
					});
				}
			} catch (error) {
				console.error(error);
				toast.error('We failed to capture this link...', {
					description: 'Our crawlers are not able to access this link.'
				});
			} finally {
				toast.dismiss(loading_toast);
			}
		}
	}

	async function handlePaste(e: ClipboardEvent) {
		const userData = e.clipboardData?.getData('text');
		addMark(userData);
	}

	function drop(event: DragEvent) {
		event.preventDefault();
		addMark(event.dataTransfer?.getData('text/plain'));
	}
</script>

<svelte:head>
	<title>Inbox • Cognitum</title>
	<script defer src="/modules/masonry.js"></script>
</svelte:head>

<svelte:window
	onpaste={handlePaste}
	ondragover={(event: Event) => event.preventDefault()}
	ondrop={(event: DragEvent) => drop(event)}
/>

{#if loading}
	<Loading />
{:else}
	<content
		class="h-full w-full flex flex-col gap-2 justify-start items-start p-1"
		role="presentation"
	>
		{#if Marks.length === 0}
			<div
				in:fade={{ duration: 500, delay: 750 }}
				class="flex flex-col h-full w-full gap-2 justify-center items-center opacity-65"
			>
				<img src="/logo.svg" alt="Cognitum" class="size-6 invert dark:invert-0" />
				<div class="text-center">
					<p class="text-lg font-medium">Your inbox is tidy!</p>
					<p class="text-xs w-64">
						Nice work! You've got 0 unsorted items in your inbox. Keep it up!
					</p>
				</div>
			</div>
		{:else if !loading}
			<div
				id="grid"
				class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-2 pb-12"
			>
				{#each Marks as item, idx}
					<LinkComponent data={item} index={idx} />
				{/each}
			</div>
		{:else}
			<div in:fade={{ duration: 250, delay: 500 }}>
				<Loading />
			</div>
		{/if}
	</content>
{/if}
