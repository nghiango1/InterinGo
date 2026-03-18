import { mdsvex } from 'mdsvex';
import adapter from '@sveltejs/adapter-static';

const config = {
	kit: {
		// adapter-auto only supports some environments, see https://svelte.dev/docs/kit/adapter-auto for a list.
		// If your environment is not supported, or you settled on a specific environment, switch out the adapter.
		// See https://svelte.dev/docs/kit/adapters for more information about adapters.
		adapter: adapter({
			fallback: '200.html', // may differ from host to host
			pages: 'dist',
			assets: 'dist',
			precompress: false,
			strict: true
		}),
		prerender: {
			crawl: true,
			entries : [
				'/',
				'/docs',
				'/docs/keyword',
			]
		}
	},
	// https://svelte.dev/docs/kit/configuration#files
	files: {
		assets: 'assets' // Public dir -> vite publicDir overide
	},
	preprocess: [
		mdsvex({
			extensions: ['.md', '.svx']
		})
	],
	extensions: ['.svelte', '.svx', '.md']
};

export default config;
