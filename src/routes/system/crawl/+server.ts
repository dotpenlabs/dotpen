import { json } from '@sveltejs/kit';
import * as cheerio from 'cheerio';

const cache = new Map();

export async function GET({ url }) {
    const targetUrl = url.searchParams.get('url');
    if (!targetUrl) {
        return json({ error: 'URL query parameter is required' }, { status: 400 });
    }

    if (cache.has(targetUrl)) {
        return json(cache.get(targetUrl));
    }

    try {
        const response = await fetch(targetUrl, {
            headers: {
                'User-Agent': 'CognitumBot/1.0'
            }
        });
        const html = await response.text();

        const $ = cheerio.load(html);
        const result = {
            title: $('meta[property="og:title"]').attr('content'),
            description: $('meta[property="og:description"]').attr('content'),
            image: $('meta[property="og:image"]').attr('content'),
            url: $('meta[property="og:url"]').attr('content'),
            author: $('meta[name="author"]').attr('content') || $('a[rel="author"]').text(),
            publishedDate: $('meta[property="article:published_time"]').attr('content') || $('time').attr('datetime') || new Date().toISOString(),
        };

        cache.set(targetUrl, result);

        return json(result);
    } catch (error) {
        console.error(error);
        return json({ error: 'Failed to fetch or parse the URL' }, { status: 500 });
    }
}