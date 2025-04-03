<script lang="ts">
	import { Eye, EyeSlash, Mailbox, Spinner, CaretRight, ArrowLeft } from 'phosphor-svelte';
	import { toast } from 'svelte-sonner';
	import { fade, fly } from 'svelte/transition';
	import { quintOut, backOut } from 'svelte/easing';
	import { focus } from '$lib/utils';
	import { onMount } from 'svelte';
	import { page } from '$app/state';
	import type { SupabaseClient } from '@supabase/supabase-js';
	import { goto, invalidate, invalidateAll } from '$app/navigation';

	const {
		view,
		onSelect,
		form = false
	} = $props<{
		view: 'login' | 'signup';
		onSelect?: () => void;
		form?: boolean;
	}>();
	let viewState = $state(form ? 1 : (0 as 0 | 1 | 2));

	let showPassword = $state(false);
	let loading = $state(false);
	let email = $state('');
	let password = $state('');
	let name = $state('');

	const { supabase } = page.data as unknown as { supabase: SupabaseClient };

	onMount(() => {
		if (form && onSelect) onSelect();
	});

	function handleInitialClick() {
		viewState = 1;
		if (onSelect) onSelect();
	}

	async function handleSubmit(e: Event) {
		e.preventDefault();
		e.stopPropagation();

		if (view === 'signup' && viewState === 1) {
			if (!name && view === 'signup') {
				toast.error('Please enter your name');
				return;
			}

			if (!email) {
				toast.error('Please enter your email');
				return;
			}

			const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
			if (!emailRegex.test(email)) {
				toast.error('Please enter a valid email address');
				return;
			}

			if (!password) {
				toast.error('Please enter your password');
				return;
			}

			loading = true;

			try {
				const { error } = await supabase.auth.signUp({
					email: email,
					password: password,
					options: {
						data: {
							display_name: name
						},
						emailRedirectTo: window.location.origin
					}
				});

				if (error) {
					toast.error(error.message || 'Account creation failed');
				} else {
					toast.success('Your account has been created!', {
						description: 'Please check your email for verification.'
					});

					window.location.href = '/auth?intent=verify';
				}
			} catch (error: any) {
				toast.error(error.message || 'Account creation failed');
			} finally {
				loading = false;
			}
			return;
		}

		if (view === 'login' && viewState === 1) {
			if (!email) {
				toast.error('Please enter your email');
				return;
			}

			const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
			if (!emailRegex.test(email)) {
				toast.error('Please enter a valid email address');
				return;
			}

			viewState = 2;
			return;
		}

		if (view === 'login' && viewState === 2) {
			if (!password) {
				toast.error('Please enter your password');
				return;
			}

			loading = true;

			try {
				const { error } = await supabase.auth.signInWithPassword({
					email,
					password
				});

				if (error) {
					toast.error(error.message || 'Authentication failed', {
						description: 'Please check your email and password and try again.'
					});
				} else {
					toast.success('Welcome back!', {
						description: "You're now logged in"
					});

					goto('/app');
				}
			} catch (error: any) {
				toast.error(error.message || 'Authentication failed');
			} finally {
				loading = false;
			}
		}
	}

	function goBack() {
		if (viewState === 2) {
			viewState = 1;
		} else if (!form) {
			viewState = 0;
		}
	}
</script>

<form
	onsubmit={handleSubmit}
	class="flex flex-col gap-4 w-full"
	in:fly={{ y: 5, duration: 400, easing: quintOut }}
>
	{#if viewState === 0}
		<button
			type="button"
			onclick={handleInitialClick}
			class="flex justify-center hover:scale-[101.75%] active:scale-[98.5%] duration-200 items-center gap-3 bg-white text-stone-700 border border-stone-300 rounded-xl p-3 hover:bg-stone-100 dark:bg-stone-800 dark:text-stone-200 dark:border-stone-700 dark:hover:bg-stone-700 transition-colors cursor-pointer min-h-[44px]"
		>
			<Mailbox class="size-5" weight="duotone" />
			<span class="font-medium mt-0.5 text-sm sm:text-base">Continue with Email</span>
		</button>
	{:else if viewState === 1}
		<div in:fade={{ duration: 300 }} class="flex flex-col gap-4">
			{#if !form}
				<div class="flex items-center mb-2">
					<button
						type="button"
						onclick={goBack}
						class="text-stone-500 hover:text-stone-700 dark:hover:text-stone-300 mr-2 text-sm cursor-pointer"
					>
						← Back
					</button>
				</div>
			{/if}

			{#if view === 'signup'}
				<div class="flex flex-col gap-4">
					<div class="flex flex-col gap-2">
						<label for="name" class="text-sm text-stone-700 dark:text-stone-300">Full Name</label>
						<input
							use:focus
							id="name"
							type="text"
							bind:value={name}
							placeholder="John Doe"
							class="w-full p-3 rounded-xl border border-stone-300 dark:border-stone-700 bg-white dark:bg-stone-800 text-stone-900 dark:text-stone-100 focus:outline-none focus:ring-2 focus:ring-stone-400 dark:focus:ring-stone-600"
						/>
					</div>

					<div class="flex flex-col gap-2">
						<label for="email" class="text-sm text-stone-700 dark:text-stone-300">Email</label>
						<input
							id="email"
							type="email"
							bind:value={email}
							placeholder="you@example.com"
							class="w-full p-3 rounded-xl border border-stone-300 dark:border-stone-700 bg-white dark:bg-stone-800 text-stone-900 dark:text-stone-100 focus:outline-none focus:ring-2 focus:ring-stone-400 dark:focus:ring-stone-600"
						/>
					</div>

					<div class="flex flex-col gap-2">
						<label for="password" class="text-sm text-stone-700 dark:text-stone-300">Password</label
						>
						<div class="relative">
							<input
								id="password"
								type={showPassword ? 'text' : 'password'}
								bind:value={password}
								placeholder="••••••••"
								class="w-full p-3 pr-12 rounded-xl border border-stone-300 dark:border-stone-700 bg-white dark:bg-stone-800 text-stone-900 dark:text-stone-100 focus:outline-none focus:ring-2 focus:ring-stone-400 dark:focus:ring-stone-600"
							/>
							<button
								type="button"
								onclick={() => (showPassword = !showPassword)}
								class="absolute right-3 top-1/2 transform -translate-y-1/2 text-stone-500 hover:text-stone-700 dark:hover:text-stone-300 cursor-pointer"
							>
								{#if showPassword}
									<EyeSlash size={20} />
								{:else}
									<Eye size={20} />
								{/if}
							</button>
						</div>
					</div>

					<button
						type="submit"
						disabled={loading}
						class="w-full mt-2 p-3 rounded-xl bg-stone-800 dark:bg-stone-200 text-white dark:text-stone-900 hover:bg-stone-700 dark:hover:bg-stone-300 hover:scale-[101%] active:scale-[99%] transition-all duration-200 cursor-pointer disabled:opacity-50 disabled:cursor-not-allowed flex justify-center items-center"
					>
						{#if loading}
							<Spinner size={20} class="animate-spin" />
						{:else}
							Sign up
						{/if}
					</button>
				</div>
			{:else}
				<div class="flex flex-col gap-2">
					<label for="email" class="text-sm text-stone-700 dark:text-stone-300">Email</label>
					<div class="flex items-center gap-2">
						<input
							use:focus
							id="email"
							type="email"
							bind:value={email}
							placeholder="you@example.com"
							class="flex-1 p-3 rounded-xl border border-stone-300 dark:border-stone-700 bg-white dark:bg-stone-800 text-stone-900 dark:text-stone-100 focus:outline-none focus:ring-2 focus:ring-stone-400 dark:focus:ring-stone-600"
						/>
						<button
							type="submit"
							aria-label="Continue"
							class="flex justify-center items-center bg-stone-800 dark:bg-stone-200 text-white dark:text-stone-900 h-[46px] w-[46px] aspect-square rounded-xl hover:bg-stone-700 dark:hover:bg-stone-300 hover:scale-[101%] active:scale-[99%] transition-all duration-200 cursor-pointer"
						>
							<CaretRight size={20} weight="bold" />
						</button>
					</div>
				</div>
			{/if}
		</div>
	{:else if viewState === 2}
		<div in:fade={{ duration: 300 }} class="flex flex-col gap-4">
			<div class="flex flex-col gap-2">
				<label for="password" class="text-sm text-stone-700 dark:text-stone-300">Password</label>
				<div class="flex items-center gap-2">
					<div class="relative flex-1">
						<input
							use:focus
							id="password"
							type={showPassword ? 'text' : 'password'}
							bind:value={password}
							placeholder="••••••••"
							class="w-full p-3 pr-12 rounded-xl border border-stone-300 dark:border-stone-700 bg-white dark:bg-stone-800 text-stone-900 dark:text-stone-100 focus:outline-none focus:ring-2 focus:ring-stone-400 dark:focus:ring-stone-600"
						/>
						<button
							type="button"
							onclick={() => (showPassword = !showPassword)}
							class="absolute right-3 top-1/2 transform -translate-y-1/2 text-stone-500 hover:text-stone-700 dark:hover:text-stone-300 cursor-pointer"
						>
							{#if showPassword}
								<EyeSlash size={20} />
							{:else}
								<Eye size={20} />
							{/if}
						</button>
					</div>
					<button
						type="submit"
						disabled={loading}
						aria-label="Log in"
						class="flex justify-center items-center bg-stone-800 dark:bg-stone-200 text-white dark:text-stone-900 h-[46px] w-[46px] aspect-square rounded-xl hover:bg-stone-700 dark:hover:bg-stone-300 hover:scale-[101%] active:scale-[99%] transition-all duration-200 disabled:opacity-50 disabled:cursor-not-allowed cursor-pointer"
					>
						{#if loading}
							<Spinner size={20} class="animate-spin" />
						{:else}
							<CaretRight size={20} weight="bold" />
						{/if}
					</button>
				</div>
			</div>
			<div class="flex items-center mb-2">
				<button
					type="button"
					onclick={goBack}
					class="text-stone-500 hover:text-stone-700 dark:hover:text-stone-300 mr-2 flex gap-1 items-center text-sm cursor-pointer"
				>
					<ArrowLeft /> Back
				</button>
				<span class="text-sm text-stone-500 dark:text-stone-400">{email}</span>
			</div>
		</div>
	{/if}
</form>
