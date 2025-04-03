import type { PageLoad } from './$types';

// This page is needed as a client entry point
// for the OAuth callback redirect but doesn't actually
// need to do anything client-side as all the work
// happens in +page.server.ts
export const load: PageLoad = async () => {
    return {};
}; 