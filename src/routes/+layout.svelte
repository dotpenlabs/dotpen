<script>
	import '../app.css';

	import { ModeWatcher, mode } from 'mode-watcher';
	import { BugBeetle } from 'phosphor-svelte';
	import { Toaster } from 'svelte-sonner';

	const { children } = $props();
</script>

<ModeWatcher defaultMode="system" />
<Toaster position="top-right" theme={$mode ? 'light' : 'dark'} richColors />

<content class="lg:block hidden">
	{@render children()}
</content>

<content
	class="lg:hidden absolute h-full w-full flex justify-center items-center gap-2 bg-black text-white"
>
	<BugBeetle class="text-red-500" weight="duotone" />
	<p>Dotpen is not yet available on mobile.</p>
</content>
