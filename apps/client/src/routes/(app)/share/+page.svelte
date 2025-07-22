<script lang="ts">
	import { Drawer } from 'vaul-svelte';
	import Loading from '$/lib/components/loading.svelte';
	import { onMount } from 'svelte';
	import { pb } from '$/lib';
	import { toast } from 'svelte-sonner';

	let collections = $state<{ id: string; name: string }[]>([]);
	let selectedCollection = $state<string>('');
	let shareData = $state({ title: '', link: '' });
	let isSaving = $state(false);

	onMount(async () => {
		const url = new URL(window.location.href);
		shareData.title = url.searchParams.get('title') || '';
		shareData.link = url.searchParams.get('url') || '';
		if (!shareData.link) {
			toast.error("Couldn't share with Dotpen!", {
				description: 'Sorry! Dotpen currently only supports links.'
			});
			return;
		}
		try {
			const userId = pb.authStore.model?.id;
			if (!userId) return;
			const result = await pb.collection('collections').getFullList({
				filter: `user = "${userId}" && name != "system_inbox"`,
				sort: 'created'
			});
			collections = result.map((col) => ({ id: col.id, name: col.name }));
			if (collections.length) selectedCollection = collections[0].id;
		} catch (e) {
			toast.error('Failed to load collections');
		}
	});

	async function saveBookmark(e: Event) {
		e.preventDefault();

		if (!selectedCollection || !shareData.link) return;
		isSaving = true;
		const urlF = (input: string): string | null => {
			if (/\s/.test(input)) return null;
			try {
				let uri = new URL(input);
				if (!uri.protocol.startsWith('http')) uri = new URL('https://' + input);
				uri.protocol = 'https:';
				const hostnameParts = uri.hostname.split('.');
				const tld = hostnameParts[hostnameParts.length - 1];
				if (!tld || tld.length < 2) return null;
				return uri.toString();
			} catch {
				try {
					const uri = new URL('https://' + input);
					const hostnameParts = uri.hostname.split('.');
					const tld = hostnameParts[hostnameParts.length - 1];
					if (!tld || tld.length < 2) return null;
					return uri.toString();
				} catch {
					return null;
				}
			}
		};
		const fileF = async (url: string, filename: string) => {
			try {
				const r = await fetch(`${pb.baseURL}api/proxy?url=${encodeURIComponent(urlF(url))}`, {
					headers: { Authorization: `Bearer ${pb.authStore.token}` }
				});
				if (!r.ok) return { file: null, base64: null };
				const blob = await r.blob();
				const file = new File([blob], filename, {
					type: blob.type || `image/${filename.split('.').pop() || 'png'}`
				});
				const base64 = await new Promise<string>((resolve, reject) => {
					const reader = new FileReader();
					reader.onloadend = () => resolve(reader.result as string);
					reader.onerror = reject;
					reader.readAsDataURL(blob);
				});
				return { file, base64 };
			} catch {
				return { file: null, base64: null };
			}
		};
		const url = urlF(shareData.link);
		if (!url) {
			toast.error('Invalid URL');
			isSaving = false;
			return;
		}
		const loading = toast.loading('Saving bookmark...');
		try {
			const crawl = (await pb.send('/api/crawl', { method: 'GET', query: { url } })) as {
				title: string;
				favicon?: string;
				image?: string;
				url: string;
			};
			['favicon', 'image'].forEach((key) => {
				if (crawl[key]?.startsWith('./')) crawl[key] = crawl[key].replace('./', '/');
				if (crawl[key]?.startsWith('/')) crawl[key] = new URL(crawl.url).origin + crawl[key];
			});
			const { file: crawl_favicon } = await fileF(crawl.favicon, 'favicon.png');
			const { file: crawl_image } = await fileF(crawl.image, 'cover.png');
			const form: Record<string, any> = {
				label: shareData.title || crawl.title,
				link: url,
				collection: selectedCollection,
				user: pb.authStore.record.id,
				created: new Date().toISOString()
			};
			if (crawl_favicon) form.favicon = crawl_favicon;
			if (crawl_image) form.cover = crawl_image;
			await pb.collection('bookmarks').create(form);
			toast.success('Bookmark saved!');
		} catch (e) {
			toast.error('Failed to save bookmark');
		} finally {
			isSaving = false;
			toast.dismiss(loading);
		}
	}
</script>

<div class="h-full w-full flex flex-col items-center justify-center bg-background p-4">
	<img src="/logo.svg" alt="Dotpen" class="size-6 invert dark:invert-0 mb-4" />
	<div class="w-full max-w-md">
		<Drawer.Root shouldScaleBackground>
			<Drawer.Trigger class="w-full px-4 py-3 bg-primary text-primary-foreground rounded mb-4"
				>Share to Dotpen</Drawer.Trigger
			>
			<Drawer.Portal>
				<Drawer.Overlay class="fixed inset-0 bg-black/40" />
				<Drawer.Content
					class="fixed bottom-0 left-0 right-0 mt-24 flex flex-col rounded-t-[10px] bg-card p-6 max-w-md mx-auto"
				>
					<div class="mx-auto mb-4 h-1.5 w-12 flex-shrink-0 rounded-full bg-zinc-300"></div>
					<h2 class="text-lg font-semibold mb-2">Save to Dotpen</h2>
					<form class="flex flex-col gap-4" onsubmit={saveBookmark}>
						<label class="flex flex-col gap-1">
							<span class="text-sm font-medium">Collection</span>
							<select class="rounded border px-3 py-2" bind:value={selectedCollection} required>
								{#each collections as col}
									<option value={col.id}>{col.name}</option>
								{/each}
							</select>
						</label>
						<label class="flex flex-col gap-1">
							<span class="text-sm font-medium">Title</span>
							<input
								class="rounded border px-3 py-2"
								type="text"
								bind:value={shareData.title}
								placeholder="Title"
								required
							/>
						</label>
						<label class="flex flex-col gap-1">
							<span class="text-sm font-medium">Link</span>
							<input
								class="rounded border px-3 py-2"
								type="url"
								bind:value={shareData.link}
								placeholder="https://..."
								required
							/>
						</label>
						<button
							type="submit"
							class="w-full py-2 rounded bg-primary text-primary-foreground font-semibold disabled:opacity-60"
							disabled={isSaving || !selectedCollection || !shareData.link}
						>
							{isSaving ? 'Saving...' : 'Save'}
						</button>
					</form>
				</Drawer.Content>
			</Drawer.Portal>
		</Drawer.Root>
	</div>
</div>
