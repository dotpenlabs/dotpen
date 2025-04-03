<script lang="ts">
	import { fade } from 'svelte/transition';
	import { Rotate } from '../utils';
</script>

<div
	transition:fade
	class="absolute top-0 left-0 h-full w-full flex items-center justify-center bg-[#F7EEE9] dark:bg-[#212020]"
>
	<img
		src="/logo.svg"
		alt="Cognitum Logo"
		transition:Rotate={{ duration: 5000 }}
		class="size-5 invert-100 animate-spin dark:invert-0 hover:scale-110 active:scale-90 transition-all duration-1000 dark:opacity-90"
		on:introend={(event) => {
			event.currentTarget.classList.add('rotateAnim');
		}}
	/>
</div>

<style>
	@keyframes rotateAnimation {
		0% {
			transform: rotate(0deg);
		}
		100% {
			transform: rotate(360deg);
		}
	}

	:global(.rotateAnim) {
		animation: rotateAnimation 1.75s cubic-bezier(0.4, 0, 0.2, 1) infinite;
		animation-delay: 3s;
		animation-fill-mode: both;
	}
</style>
