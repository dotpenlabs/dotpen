<div
	class="fixed inset-0 z-50 flex flex-col gap-4 items-center justify-center bg-white dark:bg-black"
>
	<img src="/full.png" alt="Dotpen Logo" class="h-8 opacity-90 invert-100 dark:invert-0" />
	<div class="w-64 h-1.5 bg-black/25 dark:bg-white/25 rounded-full overflow-hidden relative">
		<div
			class="absolute h-full bg-black/50 dark:bg-white/50 rounded-full animate-loadingBar-alternate"
		></div>
		<div class="absolute h-full bg-black/50 dark:bg-white/50 rounded-full animate-loadingBar"></div>
	</div>
</div>

<style>
	@keyframes loadingBar {
		0% {
			left: -10%;
			width: 5%;
		}
		50% {
			width: 65%;
		}
		100% {
			left: 105%;
			width: 5%;
		}
	}

	.animate-loadingBar-alternate {
		animation: loadingBar 1s ease-in-out infinite alternate;
	}

	.animate-loadingBar {
		animation: loadingBar 1.25s ease-in-out infinite alternate;
		animation-delay: 0.5s;
	}

	.animate-loadingBar,
	.animate-loadingBar-alternate {
		will-change: left, width;
	}
</style>
