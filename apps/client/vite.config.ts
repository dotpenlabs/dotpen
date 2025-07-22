import tailwindcss from '@tailwindcss/vite';
import devtoolsJson from 'vite-plugin-devtools-json';

import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';

export default defineConfig({
	plugins: [tailwindcss(), sveltekit(), devtoolsJson()],
	define: {
		APP_VERSION: JSON.stringify(process.env.npm_package_version)
	},
	server: {
		watch: {
			ignored: ['**/server/**']
		}
	}
});
