import { redirect } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ url, locals: { supabase } }) => {
    const email = url.searchParams.get('email');
    const token = url.searchParams.get('token_hash');

    if (email && token) {
        const { data: user, error: userError } = await supabase.auth.getUser();

        if (userError || !user) {
            console.error('Error fetching user:', userError?.message);
            throw redirect(303, '/auth?login_error=' + "otp_failed");
        }

        const { data, error } = await supabase.auth.verifyOtp({
            token_hash: token,
            type: 'email'
        });

        if (error) {
            console.error('Error verifying OTP:', error.message);
            throw redirect(303, '/auth?login_error=' + "otp_failed");
        } else {
            console.log('OTP verified successfully');
            throw redirect(303, '/app');
        }
    }

    console.log('No email or token provided');
    throw redirect(303, '/app');
}; 