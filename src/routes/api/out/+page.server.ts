import { redirect } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';

export const load = (async ({ url }) => {
	const goto = url.searchParams.get('url');

	if (!goto) {
		return {
			error: 'No URL provided',
		};
	}

	return redirect(301, goto);
}) satisfies PageServerLoad;
