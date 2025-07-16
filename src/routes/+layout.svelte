<script>
	import '../app.css';

	import { ModeWatcher } from 'mode-watcher';
	import { Toaster } from 'svelte-sonner';

	const { children } = $props();
</script>

<ModeWatcher defaultMode="system" />
<Toaster position="top-right" theme="system" richColors />

{@render children()}
