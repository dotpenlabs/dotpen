<script>
	import '../app.css';

	import { ModeWatcher, mode, setMode } from 'mode-watcher';
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

Watch out! Pasting code into the development console can be dangerous and is not recommended. If you know what you're doing, come help us out https://github.com/dotpenlabs/dotpen!

Your currently running v${APP_VERSION}, share this number when reporting bugs!
`,
			`font-family: monospace`
		);

		setMode('system');
	});
</script>

<ModeWatcher defaultMode="system" />
<Toaster position="top-right" theme={$mode == 'dark' ? 'dark' : 'light'} richColors />

<content>
	{@render children()}
</content>
