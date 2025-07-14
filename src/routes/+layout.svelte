<script lang="ts">
	import '../app.css';
	import { onMount } from 'svelte';
	import { ModeWatcher, mode } from 'mode-watcher';
	import { Toaster } from 'svelte-sonner';

	let { data, children } = $props();
</script>

<ModeWatcher defaultMode="system" />
<Toaster position="top-right" theme={$mode === 'dark' ? 'dark' : 'light'} richColors />
<div class="lg:block hidden">
	{@render children()}
</div>
<div class="lg:hidden h-full w-full flex justify-center items-center">
	<p>Dotpen is currently unavailable on mobile.</p>
</div>
