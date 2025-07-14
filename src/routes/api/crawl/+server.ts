import { json } from '@sveltejs/kit';
import fetch from 'node-fetch';
import { createRequire } from 'module';

interface MetaData {
	title?: string;
	description?: string;
	image?: string;
	url?: string;
	author?: string;
	date?: string | null;
	[key: string]: any;
}

const require = createRequire(import.meta.url);
const metascraper = require('metascraper')([
	require('metascraper-title')(),
	require('metascraper-description')(),
	require('metascraper-image')(),
	require('metascraper-url')(),
	require('metascraper-author')(),
	require('metascraper-date')(),
	require('metascraper-image')(),
	require('metascraper-twitter')()
]);

const cache = new Map<string, any>();

async function useTwitter(url: string) {
	const data = await fetch(`https://publish.twitter.com/oembed?url=${encodeURIComponent(url)}`);
	if (data.ok) {
		const { html, author_name } = (await data.json()) as {
			html?: string;
			author_name?: string;
		};
		const tweetText = html.match(/<p>(.*?)<\/p>/)?.[1] || '';

		const response = {
			title: tweetText ? `${author_name}: ${tweetText}` : `${author_name || 'Unknown'}'s tweet`,
			url,
			author: author_name || null,
			date: null
		};

		console.log(response);
		return response;
	}
}

async function fetchMetadata(url: string) {
	let meta: MetaData | undefined;

	if (url.includes('twitter.com') || url.includes('t.co') || url.includes('x.com')) {
		meta = await useTwitter(url);
	} else {
		const res = await fetch(url, {
			headers: {
				'User-Agent':
					'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.0.0 Safari/537.36'
			}
		});
		if (!res.ok) throw new Error('Fetch failed');
		const html = await res.text();
		const scrapedMeta = await metascraper({ html, url });
		meta = {
			...scrapedMeta
		};
	}

	return { ...meta, url };
}

export async function GET({ url }: { url: URL }) {
	const targetUrl = url.searchParams.get('url');
	if (!targetUrl) return json({ error: 'URL query parameter is required' }, { status: 400 });

	if (cache.has(targetUrl)) return json(cache.get(targetUrl));

	try {
		const metadata = await fetchMetadata(targetUrl);
		cache.set(targetUrl, metadata);
		return json(metadata);
	} catch {
		return json({ error: 'Failed to fetch or parse the URL' }, { status: 500 });
	}
}
