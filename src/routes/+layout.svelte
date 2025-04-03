<script lang="ts">
	import '../app.css';

	import { invalidate } from '$app/navigation';
	import { onMount } from 'svelte';
	import { ModeWatcher, mode } from 'mode-watcher';
	import { Toaster } from 'svelte-sonner';

	let { data, children } = $props();
	let { session, supabase } = $derived(data);

	onMount(() => {
		const { data } = supabase.auth.onAuthStateChange((_, newSession) => {
			if (newSession?.expires_at !== session?.expires_at) {
				invalidate('supabase:auth');
			}
		});

		return () => data.subscription.unsubscribe();
	});
</script>

<ModeWatcher defaultMode="system" />
<Toaster position="top-right" theme={$mode === 'dark' ? 'dark' : 'light'} richColors />
{@render children()}
