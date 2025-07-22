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

<content>
	{@render children()}
</content>
