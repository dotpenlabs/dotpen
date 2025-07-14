import { error } from '@sveltejs/kit';

export async function GET({ url }) {
    const targetUrl = url.searchParams.get('url');
    if (!targetUrl || !/^https?:\/\//.test(targetUrl)) {
        throw error(400, 'Invalid or missing url parameter');
    }

    try {
        const response = await fetch(targetUrl);
        if (!response.ok) {
            throw error(502, 'Failed to fetch image');
        }
        const contentType = response.headers.get('content-type') || 'application/octet-stream';
        const body = new Uint8Array(await response.arrayBuffer());
        return new Response(body, {
            headers: {
                'content-type': contentType,
                'cache-control': 'public, max-age=86400'
            }
        });
    } catch (err) {
        throw error(500, 'Proxy error');
    }
} 