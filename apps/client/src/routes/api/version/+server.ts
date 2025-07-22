import type { RequestHandler } from './$types';

export const GET: RequestHandler = async () => {
	return new Response(
		JSON.stringify({
			version: APP_VERSION,
			build: process.env.npm_package_version,
			commit: process.env.COMMIT_SHA,
			branch: process.env.COMMIT_BRANCH,
			date: new Date().toISOString()
		})
	);
};
