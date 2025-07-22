<!-- <script lang="ts">
	import { fade, fly } from 'svelte/transition';
	import { onMount } from 'svelte';
	import { cubicOut, elasticOut } from 'svelte/easing';
	import { ParseParams, Reveal, Rotate, flyAndScale } from '$/lib/utils';
	import { mode, setMode } from 'mode-watcher';
	import { ArrowRight, Moon, Sun, UserCircle, Compass } from 'phosphor-svelte';
	import Loading from '$/lib/components/loading.svelte';
	import Emblem from '$/lib/components/emblem.svelte';
	import { toast } from 'svelte-sonner';
	import { goto } from '$app/navigation';

	let loading = $state(true);../$types.js

	let { data } = $props();
	const features = [
		{
			title: 'Capture',
			description: 'Save links and notes instantly with a single click'
		},
		{
			title: 'Organize',
			description: 'Organise all of your links, notes and ideas in one place!'
		},
		{
			title: 'Collaborate',
			description: 'Easily share links, with coworkers, friends and family!'
		}
	];

	function spring(node: Element, { delay = 0, duration = 400, easing = cubicOut, x = 0, y = 0 }) {
		const style = getComputedStyle(node);
		const target_opacity = +style.opacity;

		return {
			delay,
			duration,
			easing,
			css: (t: number) => `
                transform: translate(${(1 - t) * x}px, ${(1 - t) * y}px) scale(${0.5 + t * 0.5});
                opacity: ${t * target_opacity};
            `
		};
	}

	function scrollToSection(id: string) {
		const element = document.getElementById(id);
		if (element) {
			element.scrollIntoView({ behavior: 'smooth' });
		}
	}

	onMount(() => {
		const { error, error_code } = ParseParams(window.location.href);

		const loginTimestamp = sessionStorage.getItem('login:timestamp');
		const currentTime = Date.now();
		const timeDiff = currentTime - parseInt(loginTimestamp || '0');

		if (timeDiff < 1000 * 60 * 60 * 24) {
			sessionStorage.removeItem('login:timestamp');

			toast.success('Welcome back!', {
				description: "You're now logged in"
			});

			goto('/app');
		}

		if (error) {
			switch (error) {
				case 'access_denied':
					toast.error('Social Login denied', {
						description: "We can't authenicated via your social login. Please try again."
					});
					break;
				default:
					console.error(error);
					toast.error('There was an error logging in', {
						description: error.toString()
					});
					break;
			}
		}

		if (error_code) {
			switch (error_code) {
				case 'invalid_credentials':
					toast.error('Invalid credentials');
					break;
				case 'user_banned':
					toast.error('Your account is currently suspended', {
						description: 'Please contact support if you believe this is an error.'
					});
					break;
				case 'email_not_verified':
					toast.error('Email not verified');
					break;
				default:
					console.error(error_code);
					toast.error('There was an error logging in', {
						description: error_code.toString()
					});
					break;
			}
		}

		setTimeout(() => {
			loading = false;
		}, 300);
	});
</script>

<svelte:head>
	<title>Dotpen - Capture links, notes, and ideas</title>
	<meta
		name="description"
		content="Dotpen is een snelle en eenvoudige manier om ideeën, bladwijzers en meer op te slaan!"
	/>
</svelte:head>

{#if !loading}
	<content
		in:fade={{ duration: 800 }}
		class="flex flex-col bg-[#FAF5F2] dark:bg-[#171616] transition-colors duration-500 ease-in-out justify-center items-center min-h-screen w-full"
	>
		<nav
			class="w-full max-w-6xl px-4 mx-auto flex justify-between items-center absolute top-0 left-0 right-0 p-4"
		>
			<a
				href="/"
				class="flex items-center group hover:scale-[102%] active:opacity-50 active:scale-95 duration-200"
			>
				<Emblem />
			</a>

			<div class="flex items-center gap-4">
				<a
					href="/auth?intent=login"
					class="flex items-center gap-2 bg-stone-200/40 dark:bg-stone-800/40 hover:bg-stone-200/60 dark:hover:bg-stone-800/60 px-3 py-2 rounded-xl text-sm md:text-base transition-all"
				>
					<UserCircle class="size-5" />
					<span>Login</span>
				</a>
				<button
					class="flex justify-center items-center cursor-pointer size-10 hover:scale-110 active:scale-90 transition-all duration-300"
					onclick={() => {
						if ($mode === 'light') {
							setMode('dark');
						} else {
							setMode('light');
						}
					}}
				>
					{#if $mode === 'dark'}
						<div in:Rotate={{ duration: 3000, easing: elasticOut }}>
							<Sun class="size-5" />
						</div>
					{:else}
						<div in:Rotate={{ duration: 3000, easing: elasticOut }}>
							<Moon class="size-5" />
						</div>
					{/if}
				</button>
			</div>
		</nav>

		<section class="flex justify-center items-center w-full min-h-[80vh] px-4 mt-16">
			<div class="flex flex-col justify-center items-center gap-8 max-w-4xl w-full text-center">
				<div
					class="flex flex-col justify-center items-center gap-6 w-full md:w-[90%] lg:w-[95%] xl:w-[60rem] text-center"
				>
					<h1
						class="text-4xl sm:text-5xl md:text-6xl lg:text-7xl font-medium text-wrap leading-20"
						in:Reveal={{ duration: 1200 }}
					>
						Be the new <span class="text-orange-700 dark:text-orange-300">superhero</span> of your mind.
					</h1>
					<p
						class="text-sm md:text-base lg:text-lg opacity-65 font-medium w-full sm:w-9/12 md:w-8/12"
						in:fade={{ delay: 800, duration: 800 }}
					>
						Experience a platform that captures what matters most and enhances your productivity
						with intuitive tools to work smarter.
					</p>
					<div class="flex gap-4 mt-4 w-full justify-center">
						<button
							onclick={() => {
								const email = prompt('Enter your email for the waitlist');
								if (email) {
									toast.success('Successfully added to the waitlist!', {
										description: `Thank you, ${email.split('@')[0][0].toUpperCase() + email.split('@')[0].slice(1)}! You have been successfully added to our waitlist.`
									});
								} else {
									toast.error('Email entry was cancelled or invalid.');
								}
							}}
							in:flyAndScale={{ y: 20, duration: 500, delay: 150 }}
							class="bg-zinc-900 cursor-pointer active:scale-98 dark:bg-zinc-100 text-white dark:text-zinc-900 py-3 px-6 rounded-xl font-medium text-sm md:text-base inline-flex items-center justify-center hover:opacity-90 transition-all"
						>
							Join waitlist <ArrowRight class="size-4 ml-2" weight="bold" />
						</button>

						<a
							href="#features"
							onclick={() => scrollToSection('features')}
							in:flyAndScale={{ y: 20, duration: 500, delay: 150 }}
							class="bg-zinc-200 active:scale-98 dark:bg-zinc-800 text-zinc-800 dark:text-zinc-200 py-3 px-6 rounded-xl font-medium text-sm md:text-base inline-flex items-center justify-center hover:opacity-75 duration-200 transition-all"
						>
							Explore <Compass class="size-4 ml-2" weight="bold" />
						</a>
					</div>
				</div>
			</div>
		</section>

		<section class="flex justify-center items-start h-[90vh] w-full px-4">
			<div class="flex flex-col justify-center items-center w-full text-center">
				<video
					in:fly={{ y: 50, duration: 1000, delay: 600, easing: elasticOut }}
					class="rounded-2xl md:rounded-3xl w-[95%] sm:w-[90%] md:w-[80%] lg:w-[70%] shadow-lg"
					muted={true}
					loop={true}
					autoplay
					src="/assets/marketing/video.webm"
				></video>
			</div>
		</section>

		<section id="features" class="flex justify-center items-center w-full px-4 mb-20">
			<div class="max-w-5xl w-full">
				<div class="text-center mb-12">
					<h2
						class="text-2xl md:text-3xl font-medium mb-4"
						in:fly={{ y: 30, duration: 600, delay: 200 }}
					>
						Simple. <span class="text-orange-700 dark:text-orange-300">Powerful</span>. Secure.
					</h2>
				</div>

				<div class="grid grid-cols-1 md:grid-cols-3 gap-8 mt-8">
					{#each features as feature, i}
						<div class="p-6 rounded-2xl" in:fly={{ y: 20, duration: 500, delay: 200 + i * 100 }}>
							<h3 class="text-xl font-medium mb-3">{feature.title}</h3>
							<p class="text-stone-600 dark:text-stone-400">{feature.description}</p>
						</div>
					{/each}
				</div>
			</div>
		</section>

		<footer class="w-full bg-stone-100 dark:bg-zinc-800/30 py-12 mt-auto">
			<div class="max-w-6xl mx-auto px-4">
				<div class="grid grid-cols-1 md:grid-cols-4 gap-8 w-full">
					<div class="col-span-1 md:col-span-1">
						<Emblem />
						<p class="text-sm text-stone-600 dark:text-stone-400 mt-4">
							A minimal way to save and organize your digital life.
						</p>
					</div>

					<div>
						<h4 class="text-sm font-semibold mb-4">Product</h4>
						<ul class="space-y-2">
							<li>
								<a
									href="#features"
									class="text-sm text-stone-600 dark:text-stone-400 hover:text-orange-700 dark:hover:text-orange-300"
									>Features</a
								>
							</li>
							<li>
								<a
									href="#pricing"
									class="text-sm text-stone-600 dark:text-stone-400 hover:text-orange-700 dark:hover:text-orange-300"
									>Pricing</a
								>
							</li>
							<li>
								<a
									href="/download"
									class="text-sm text-stone-600 dark:text-stone-400 hover:text-orange-700 dark:hover:text-orange-300"
									>Download</a
								>
							</li>
						</ul>
					</div>

					<div>
						<h4 class="text-sm font-semibold mb-4">Company</h4>
						<ul class="space-y-2">
							<li>
								<a
									href="https://bijsven.nl"
									class="text-sm text-stone-600 dark:text-stone-400 hover:text-orange-700 dark:hover:text-orange-300"
									>Founder</a
								>
							</li>
							<li>
								<a
									href="/blog"
									class="text-sm text-stone-600 dark:text-stone-400 hover:text-orange-700 dark:hover:text-orange-300"
									>Blog</a
								>
							</li>
							<li>
								<a
									href="mailto:mail@bijsven.nl"
									class="text-sm text-stone-600 dark:text-stone-400 hover:text-orange-700 dark:hover:text-orange-300"
									>Contact</a
								>
							</li>
						</ul>
					</div>

					<div>
						<h4 class="text-sm font-semibold mb-4">Legal</h4>
						<ul class="space-y-2">
							<li>
								<a
									href="/policy/terms"
									class="text-sm text-stone-600 dark:text-stone-400 hover:text-orange-700 dark:hover:text-orange-300"
									>Terms of Service</a
								>
							</li>
							<li>
								<a
									href="/policy/privacy"
									class="text-sm text-stone-600 dark:text-stone-400 hover:text-orange-700 dark:hover:text-orange-300"
									>Privacy Policy</a
								>
							</li>
							<li>
								<a
									href="/policy/cookies"
									class="text-sm text-stone-600 dark:text-stone-400 hover:text-orange-700 dark:hover:text-orange-300"
									>Cookie Policy</a
								>
							</li>
							<li>
								<a
									href="/policy/gdpr"
									class="text-sm text-stone-600 dark:text-stone-400 hover:text-orange-700 dark:hover:text-orange-300"
									>GDPR</a
								>
							</li>
						</ul>
					</div>
				</div>

				<div
					class="border-t border-stone-200 dark:border-stone-700 mt-8 pt-8 flex flex-col md:flex-row justify-between items-center"
				>
					<p class="text-xs text-stone-500 dark:text-stone-400">
						&copy; {new Date().getFullYear()} Dotpen. All rights reserved.
					</p>
					<div class="flex gap-4 mt-4 md:mt-0">
						<a
							href="https://twitter.com/cgntm"
							class="text-stone-500 dark:text-stone-400 hover:text-orange-700 dark:hover:text-orange-300"
							>Twitter</a
						>
					</div>
				</div>
			</div>
		</footer>
	</content>
{:else}
	<Loading />
{/if} -->
