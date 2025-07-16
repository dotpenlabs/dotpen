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
	class="bg-amber-50/65 h-full w-full absolute flex flex-col justify-center items-center gap-6"
>
	<p class="italic text-4xl text-black/65 text-center font-medium">
		Be the <span class="text-black">superhero</span><br /> of your
		<span class="text-black">mind</span>.
	</p>

	<button
		onclick={UseGoogle}
		class="bg-[#CCA281] hover:opacity-95 w-72 h-12 !font-sans text-[16px] cursor-pointer active:scale-99 active:opacity-90 duration-100 tracking-wide font-medium flex items-center justify-center gap-3 rounded-xl text-white"
	>
		<GoogleLogo weight="bold" class="opacity-80 size-5" />
		Continue with Google
	</button>
</content>
