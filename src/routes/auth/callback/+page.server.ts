import { redirect } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ url, locals: { supabase } }) => {
    const code = url.searchParams.get('code');

    if (code) {
        await supabase.auth.exchangeCodeForSession(code);

        const { data: user, error: userError } = await supabase.auth.getUser();
        if (userError) {
            throw redirect(303, '/auth?login_error=oauth');
        } else {
            throw redirect(303, '/app');
        }
    }

    throw redirect(303, '/app');
}; 