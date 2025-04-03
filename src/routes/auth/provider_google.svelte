<script lang="ts">
	import { toast } from 'svelte-sonner';
	import { page } from '$app/stores';
	import type { SupabaseClient } from '@supabase/supabase-js';
	import { time } from 'motion';

	const { onSelect } = $props<{ onSelect?: () => void }>();

	async function signInWithGoogle() {
		if (onSelect) onSelect();

		const { supabase } = $page.data as unknown as { supabase: SupabaseClient };

		try {
			const { data, error } = await supabase.auth.signInWithOAuth({
				provider: 'google',
				options: {
					redirectTo: window.location.origin
				}
			});

			sessionStorage.setItem('login:timestamp', Date.now().toString());

			if (error) throw error;

			const { data: user, error: userError } = await supabase.auth.getUser();
			if (userError) throw userError;
		} catch (error) {
			console.error('Error signing in with Google:', error);
			toast.error('Failed to connect', {
				description: 'Unable to connect to Google, please try again later'
			});
		}
	}
</script>

<button
	onclick={signInWithGoogle}
	class="flex justify-center hover:scale-[101.75%] active:scale-[98.5%] duration-200 items-center gap-3 bg-white text-stone-700 border border-stone-300 rounded-xl p-3 hover:bg-stone-100 dark:bg-stone-800 dark:text-stone-200 dark:border-stone-700 dark:hover:bg-stone-700 cursor-pointer min-h-[44px]"
>
	<img src="/assets/providers/google.svg" alt="Google" class="size-5" />
	<span class={'!font-["Roboto"] font-medium mt-0.5 text-sm sm:text-base'}
		>Continue with Google</span
	>
</button>
