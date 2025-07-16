<script>
	import '../app.css';

	import { ModeWatcher, mode } from 'mode-watcher';
	import { Toaster } from 'svelte-sonner';

	const { children } = $props();
</script>

<ModeWatcher defaultMode="system" />
<Toaster position="top-right" theme={$mode ? 'light' : 'dark'} richColors />

{@render children()}
