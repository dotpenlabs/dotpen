<script lang="ts">
	import { flyAndScale, focus, ParseParams } from '$lib/utils';
	import { goto } from '$app/navigation';
	import {
		CaretLeft,
		Envelope,
		EnvelopeOpen,
		Eye,
		EyeSlash,
		Package,
		Spinner
	} from 'phosphor-svelte';
	import { onMount } from 'svelte';
	import { toast } from 'svelte-sonner';
	import { fade, fly, scale, blur, slide } from 'svelte/transition';
	import { elasticOut, quintOut, quintIn, backOut } from 'svelte/easing';
	import ProviderGoogle from './provider_google.svelte';
	import ProviderEmail from './provider_email.svelte';
	import Emblem from '$lib/components/emblem.svelte';
	import Loading from '$lib/components/loading.svelte';
	import { page } from '$app/state';
	import ProviderNotion from './provider_notion.svelte';

	let view = $state('login') as 'login' | 'signup' | 'verify';
	let loginState = $state(0) as 0 | 1 | 2;
	let loading = $state(true);
	let logo_hovered = $state(false);

	let formVisible = $state(false);
	let buttonVisible = $state(false);
	let footerVisible = $state(false);

	let selectedProvider = $state<null | 'google' | 'notion' | 'email'>(null);

	let fields = $state({
		name: '',
		username: '',
		email: '',
		password: ''
	});

	onMount(async () => {
		const { intent, username, state, login_error } = ParseParams(window.location.href);

		view = intent as 'login' | 'signup' | 'verify';

		if (login_error) {
			switch (login_error) {
				case 'invalid_credentials':
					toast.error('Invalid credentials', {
						description: 'Please try again'
					});
					break;
				case 'otp_failed':
					toast.error('OTP verification failed', {
						description: 'Please try again'
					});
					break;
				case 'oauth':
					toast.error('Login failed', {
						description: 'Please try again'
					});
					break;
				default:
					toast.error('Login failed', {
						description: 'Please try again'
					});
					break;
			}
		}

		const img = new Image();
		img.src = '/assets/marketing/auth.jpg';
		await new Promise((resolve) => {
			img.onload = () => {
				resolve(true);
			};
		});

		if (username) {
			fields.email = username;
		}

		if (state !== null) {
			loginState = parseInt(state) as 0 | 1 | 2;
		}

		loading = false;

		setTimeout(() => (formVisible = true), 50);
		setTimeout(() => (buttonVisible = true), 100);
		setTimeout(() => (footerVisible = true), 150);
	});

	function toggleView() {
		formVisible = false;
		buttonVisible = false;
		footerVisible = false;
		selectedProvider = null;

		setTimeout(() => {
			view = view === 'login' ? 'signup' : 'login';
			setTimeout(() => (formVisible = true), 100);
			setTimeout(() => (buttonVisible = true), 400);
			setTimeout(() => (footerVisible = true), 700);
		}, 300);
	}

	function handleProviderSelect(provider: 'google' | 'email' | 'notion') {
		selectedProvider = provider;
	}

	function resetProvider() {
		selectedProvider = null;
	}
</script>

<svelte:head>
	{#if view === 'login'}
		<title>Login • Cognitum</title>
	{:else if view === 'signup'}
		<title>Register • Cognitum</title>
	{:else if view === 'verify'}
		<title>Verify • Cognitum</title>
	{/if}
	<meta
		name="viewport"
		content="width=device-width, initial-scale=1.0, maximum-scale=1.0, user-scalable=no"
	/>
	<link
		href="https://fonts.googleapis.com/css2?family=Roboto:ital,wght@0,100..900;1,100..900&display=swap"
		rel="stylesheet"
	/>
</svelte:head>

{#if !loading}
	<content
		in:fade={{ duration: 800, easing: quintOut }}
		class="bg-[#FAF5F2] dark:bg-[#171616] overflow-hidden flex h-full w-full absolute top-0 left-0"
	>
		<div
			in:fly={{ y: 20, duration: 1000, delay: 200, easing: elasticOut }}
			class="w-full sm:w-[65%] md:w-[50%] lg:w-[32.5%] h-full flex flex-col justify-center items-center relative"
		>
			<div
				class="w-full h-full flex flex-col justify-between items-center p-4 sm:p-8 md:p-12 lg:p-16 xl:p-24 overflow-y-auto"
			>
				<button
					onmouseover={() => {
						logo_hovered = true;
					}}
					onmouseout={() => {
						logo_hovered = false;
					}}
					onfocus={() => {
						logo_hovered = true;
					}}
					onblur={() => {
						logo_hovered = false;
					}}
					class="w-full flex h-12 sm:h-16 justify-start items-center"
				>
					<a
						in:fly={{ y: -30, duration: 700, easing: backOut }}
						class="w-fit flex h-min justify-start px-2 items-center cursor-pointer text-nowrap"
						href="/"
					>
						{#if !logo_hovered}
							<div class="pointer-events-none" in:scale={{ duration: 300, start: 0.8, opacity: 0 }}>
								<Emblem />
							</div>
						{:else}
							<div
								class="flex gap-2 justify-start items-center pointer-events-none"
								in:flyAndScale={{ x: -10, y: 0, duration: 250 }}
							>
								<CaretLeft class="size-5 text-stone-700 dark:text-stone-300" />
								<p
									in:fly={{ x: -10, duration: 200, easing: quintOut }}
									class="text-stone-700 dark:text-stone-300"
								>
									Back to homepage
								</p>
							</div>
						{/if}
					</a>
				</button>
				<div
					class="flex flex-col gap-3 my-auto w-full max-w-md mx-auto"
					in:fly={{ y: 30, duration: 800, easing: elasticOut }}
				>
					{#if selectedProvider && formVisible}
						<button
							type="button"
							onclick={resetProvider}
							class="text-stone-500 hover:text-stone-700 dark:hover:text-stone-300 cursor-pointer mb-4 text-sm self-start"
						>
							← Back to all options
						</button>
					{/if}

					{#if view === 'login'}
						<p
							in:fly={{ y: -20, duration: 500, delay: 200, easing: quintOut }}
							out:fly={{ y: 20, duration: 300, easing: quintIn }}
							class="text-2xl sm:text-3xl font-medium text-stone-900 dark:text-stone-100"
						>
							Welcome back!
						</p>
						<p
							in:fly={{ y: -15, duration: 500, delay: 400, easing: quintOut }}
							out:fly={{ y: 15, duration: 300, delay: 100, easing: quintIn }}
							class="text-sm sm:text-base text-stone-600 dark:text-stone-400"
						>
							Login to Sublink, the platform designed to make finding and organizing your bookmarks
							effortless and hassle-free.
						</p>

						{#if formVisible}
							<div
								in:fly={{ y: -10, duration: 500, delay: 600, easing: quintOut }}
								out:fly={{ y: 10, duration: 300, delay: 200, easing: quintIn }}
								class="flex flex-col gap-4 mt-6"
							>
								{#if !selectedProvider}
									<ProviderGoogle onSelect={() => handleProviderSelect('google')} />
									<div class="w-full border-t border-stone-300 dark:border-stone-700 my-2 relative">
										<span
											class="absolute top-1/2 left-1/2 transform -translate-x-1/2 -translate-y-1/2 bg-[#FAF5F2] dark:bg-[#171616] px-3 text-sm text-stone-500 dark:text-stone-400"
											>or</span
										>
									</div>
									<ProviderNotion onSelect={() => handleProviderSelect('notion')} />
								{:else if selectedProvider === 'email'}
									<ProviderEmail {view} form={true} />
								{/if}
							</div>
						{/if}
					{:else if view === 'signup'}
						<p
							in:fly={{ y: -20, duration: 500, delay: 200, easing: quintOut }}
							out:fly={{ y: 20, duration: 300, easing: quintIn }}
							class="text-2xl sm:text-3xl font-medium text-stone-900 dark:text-stone-100"
						>
							Create an account
						</p>
						<p
							in:fly={{ y: -15, duration: 500, delay: 400, easing: quintOut }}
							out:fly={{ y: 15, duration: 300, delay: 100, easing: quintIn }}
							class="text-sm sm:text-base text-stone-600 dark:text-stone-400"
						>
							Sign up for Sublink and easily organize and access your bookmarks anytime, anywhere.
						</p>

						{#if formVisible}
							<div
								in:fly={{ y: -10, duration: 500, delay: 600, easing: quintOut }}
								out:fly={{ y: 10, duration: 300, delay: 200, easing: quintIn }}
								class="flex flex-col gap-4 mt-6"
							>
								{#if !selectedProvider}
									<ProviderGoogle onSelect={() => handleProviderSelect('google')} />
									<div class="w-full border-t border-stone-300 dark:border-stone-700 my-2 relative">
										<span
											class="absolute top-1/2 left-1/2 transform -translate-x-1/2 -translate-y-1/2 bg-[#FAF5F2] dark:bg-[#171616] px-3 text-sm text-stone-500 dark:text-stone-400"
											>or</span
										>
									</div>
									<ProviderNotion onSelect={() => handleProviderSelect('notion')} />
								{/if}
							</div>
						{/if}
					{:else if view === 'verify'}
						<Envelope class="size-12" />
						<p
							in:fly={{ y: -20, duration: 500, delay: 200, easing: quintOut }}
							out:fly={{ y: 20, duration: 300, easing: quintIn }}
							class="text-2xl sm:text-3xl font-medium text-stone-900 dark:text-stone-100"
						>
							Verify your email
						</p>
						<p
							in:fly={{ y: -15, duration: 500, delay: 400, easing: quintOut }}
							out:fly={{ y: 15, duration: 300, delay: 100, easing: quintIn }}
							class="text-sm sm:text-base text-stone-600 dark:text-stone-400"
						>
							Please check your email for a verification link, and continue from here. You can now
							close this page.
						</p>
						<button
							type="button"
							onclick={() => goto('/')}
							class="bg-zinc-900 cursor-pointer dark:bg-zinc-100 text-white dark:text-zinc-900 py-3 px-6 rounded-xl font-medium text-sm md:text-base inline-flex items-center justify-center hover:opacity-90 transition-all"
						>
							Close tab
						</button>
					{/if}
				</div>
				{#if footerVisible && !selectedProvider}
					<div
						in:fly={{ y: 30, duration: 750, delay: 500, easing: elasticOut }}
						class="w-full text-center sm:text-left"
					>
						{#if view === 'login'}
							<div class="text-stone-500 dark:text-stone-400 text-sm">
								New here? Click
								<button
									type="button"
									onclick={toggleView}
									class="underline select-none cursor-pointer hover:text-stone-800 dark:hover:text-stone-200 hover:no-underline duration-200 active:text-stone-700 dark:active:text-stone-300 bg-transparent border-none p-0 text-inherit"
								>
									here
								</button>
								to create an account.
							</div>
						{:else if view === 'signup'}
							<div class="text-stone-600 dark:text-stone-300 text-sm mb-2">
								Already have an account? Click
								<button
									type="button"
									onclick={toggleView}
									class="underline select-none cursor-pointer hover:text-stone-800 dark:hover:text-stone-200 hover:no-underline duration-200 active:text-stone-700 dark:active:text-stone-300 bg-transparent border-none p-0 text-inherit"
								>
									here
								</button>
								to log in.
							</div>
							<p class="text-stone-500 dark:text-stone-400 text-sm">
								By creating an account, you agree to our <a href="/policy/terms" class="underline"
									>Terms of Service</a
								>, <a href="/policy/privacy" class="underline">Privacy Policy</a> and
								<a href="/policy/usage" class="underline">Usage Policy</a>.
							</p>
						{/if}
					</div>
				{/if}
			</div>
		</div>
		<div class="sm:w-[35%] md:w-[50%] lg:w-[67.5%] h-full relative hidden sm:block">
			<div
				class="absolute h-full w-[25%] left-0 top-0 z-50 bg-gradient-to-l from-transparent to-[#FAF5F2] dark:to-[#171616]"
			></div>
			<img
				in:fade={{ duration: 1200, delay: 300 }}
				src={`/assets/marketing/auth.jpg`}
				aria-hidden="true"
				alt=""
				class="w-full h-full object-cover"
			/>
		</div>
	</content>
{:else}
	<Loading />
{/if}
