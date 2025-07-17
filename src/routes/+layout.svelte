<script>
	import '../app.css';

	import { ModeWatcher, mode } from 'mode-watcher';
	import { BugBeetle } from 'phosphor-svelte';
	import { onMount } from 'svelte';
	import { Toaster } from 'svelte-sonner';

	const { children } = $props();

	onMount(() => {
		console.log(
			`%c

 /$$$$$$$            /$$                               
| $$__  $$          | $$                               
| $$  \  $$ /$$$$$$ /$$$$$$   /$$$$$$  /$$$$$$ /$$$$$$$ 
| $$  | $$/$$__  $|_  $$_/  /$$__  $$/$$__  $| $$__  $$
| $$  | $| $$  \  $$ |$$   | $$  \  $| $$$$$$$|  $$  \  $$
| $$  | $| $$  | $$ | $$ /$| $$  | $| $$_____| $$   |$$
| $$$$$$$|  $$$$$$/ |  $$$$| $$$$$$$|  $$$$$$| $$   |$$
|_______/ \______/   \___/ | $$____/ \_______|__/  |__/
                           | $$                        
                           | $$                        
                           |__/                        


Dotpen Labs (C) ${new Date().getFullYear()}

Watch out! Pasting code into the development console can be dangerous and is not recommended. If you know what you're doing, come help us out iwanttodevelopdotpen@bijsven.nl
`,
			`font-family: monospace`
		);
	});
</script>

<ModeWatcher defaultMode="light" />
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
