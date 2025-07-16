import type { PageLoad } from './$types';

export const load = (async ({ params }) => {
    return {
        collection: params.collection
    };
}) satisfies PageLoad;