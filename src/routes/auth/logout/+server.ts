import { redirect } from '@sveltejs/kit';
import type { RequestHandler } from './$types';

export const GET: RequestHandler = async ({ locals: { supabase } }) => {
    const { error } = await supabase.auth.signOut();

    if (error) return new Response('Logout failed', { status: 500 });

    return redirect(303, '/');
};