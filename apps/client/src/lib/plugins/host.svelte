<script lang="ts">
	import { onMount } from 'svelte';

	let { name, src } = $props();

	let renderError = $state(false);
	let render = $state() as HTMLIFrameElement;

	function SendMessage(message: any): void {
		render.contentWindow.postMessage(message, '*');
	}

	onMount(() => {
		const hdlr = (event: MessageEvent) => {
			if (event.origin !== new URL(src).origin || event.source !== render.contentWindow) return;
			console.log(`[${name}] ${event.data}`);

			if (event.data?.type === 'update.height') {
				if (event.data.height >= 500) {
					console.error(
						`(system) [${name}] tried to set a height of ${event.data.height}px, but was rejected because it was too large.`
					);
					SendMessage({ type: 'error', reason: 'Widget is too large' });
					renderError = true;
					return;
				}
				render.style.height = `${event.data.height}px`;
			}
		};

		window.addEventListener('message', hdlr);

		return () => {
			window.removeEventListener('message', hdlr);
		};
	});
</script>

{#if renderError}
	<div class="text-center text-red-500 mt-8">
		<p>Dotpen was unable to load this plugin. Check console.</p>
	</div>
{:else}
	<iframe
		bind:this={render}
		{src}
		sandbox="allow-scripts"
		style="border: none; width: 100%; height: 200px;"
		loading="lazy"
		title={'Plugin: ' + name}
	></iframe>
{/if}
