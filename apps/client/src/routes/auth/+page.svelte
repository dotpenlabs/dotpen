<script lang="ts" module>
	declare global {
		interface Window {
			UseEmail?: (email: string, password: string) => Promise<void>;
		}
	}
</script>

<script lang="ts">
	import { pb } from '$/lib';
	import { goto } from '$app/navigation';
	import { GoogleLogo } from 'phosphor-svelte';
	import { onMount } from 'svelte';
	import { toast } from 'svelte-sonner';

	import type { ClientResponseError } from 'pocketbase';

	let isReinit = false;

	onMount(() => {
		let loading;

		window.UseEmail = async (email: string, password: string) => {
			loading = toast.loading('Authenticating...', {
				description: 'Please follow the instructions on your screen.',
				duration: Number.POSITIVE_INFINITY
			});

			try {
				await pb.collection('users').authWithPassword(email, password);
				toast.dismiss(loading);
				goto('/');
			} catch (e: any) {
				e = e as ClientResponseError;
				if (e.status === 400) {
					try {
						await pb.collection('users').create({
							email,
							password,
							passwordConfirm: password
						});

						await pb.collection('users').authWithPassword(email, password);
						toast.dismiss(loading);
						toast.success('Account created!');
						goto('/');
						return;
					} catch (err) {
						toast.dismiss(loading);
						toast.error('Account creation failed.', {
							description: err.message || 'Unknown error occurred.'
						});
						return;
					}
				}

				toast.dismiss(loading);
			}
		};

		return () => {
			toast.dismiss(loading);
			delete window.UseEmail;
		};
	});

	async function UseGoogle() {
		const loading = toast.loading('Authenicating...', {
			description: 'Please follow the instructions on your screen.',
			duration: Number.POSITIVE_INFINITY
		});
		try {
			const resp = await pb.collection('users').authWithOAuth2({ provider: 'google' });
			toast.dismiss(loading);
			goto('/');
		} catch (e) {
			toast.dismiss(loading);
			if (!isReinit) {
				toast.info('Hang tight!', {
					description:
						"There went something wrong while logging in. Hang tight, we'll will try again!"
				});
				UseGoogle();
			} else {
				toast.error('Oh no!', {
					description:
						"We tried our best to get you back on track, but we don't have magic. Please report this issue."
				});
			}
			isReinit = true;
		}
	}
</script>

<svelte:head>
	<title>Login • Dotpen</title>
</svelte:head>

<content
	class="bg-white dark:bg-stone-900 dark:text-white h-full w-full absolute flex flex-col justify-center items-center gap-6"
>
	<p class="italic text-4xl text-black/65 dark:text-white/65 text-center font-medium">
		Take charge of your <span class="text-black dark:text-white">mind</span><br /> like a
		<span class="text-black dark:text-white">superhero</span>.
	</p>

	<button
		onclick={UseGoogle}
		class="bg-black dark:bg-white dark:text-black hover:opacity-80 px-18 py-4 !font-sans text-[16px] cursor-pointer active:scale-97 active:opacity-75 duration-100 tracking-wide font-medium flex items-center justify-center gap-4 rounded-xl text-white"
	>
		<img src="/assets/login/g.png" class="invert-100 size-4.5 dark:invert-0" alt="Google icon" /> Continue
		with Google
	</button>
	<p class="text-xs opacity-65 w-96 text-center">
		By clicking "Continue with Google," you confirm that you have read and agree to Dotpen's <a
			href="/legal/terms"
			class="underline">Terms of Service</a
		>
		and <a href="/legal/privacy" class="underline">Privacy Policy</a>.
	</p>
</content>
