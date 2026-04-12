import tailwindcss from '@tailwindcss/vite';
import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig, loadEnv } from 'vite';

import { compile } from 'mdsvex';
import { readFileSync } from 'fs';
import { sync } from 'glob';
import type { DocInfo } from '$lib/type';

async function getDocs(): Promise<DocInfo[]> {
	let docs: DocInfo[] = [];
	const pattern = 'src/lib/docs/**/*.md';
	try {
		const files = sync(pattern);
		for (let i = 0; i < files.length; i++) {
			const file = readFileSync(files[i], 'utf8');
			const mdx = await compile(file);
			const href = files[i].split('/').pop()?.replace('.md', '');
			if (!href) {
				console.error('[ERROR] Failed to get slug name for', files[i]);
				continue
			}

			const frontmatter = mdx?.data?.fm as Record<string, unknown> | undefined;
			// Get front matter title value
			const title = frontmatter?.title as string | undefined;
			const session = frontmatter?.session as string | undefined;
			const index = frontmatter?.index as number | undefined;

			// Need to cover the whole interface, as we will need to use JSON.stringify
			docs.push({
				slug: href,
				session: session, // Session title
				index: index, // Session position
				title: title
			});
		}
	} catch (err) {
		console.error('Error fetching files: ', err);
	}

	return docs;
}

export default defineConfig(async ({ mode }) => {
	const env = loadEnv(mode, process.cwd(), '')

	return {
		server: {
			port: 3000,
			open: true,
			host: true,
			proxy: {
				'/api': {
					target: env.BACKEND_SERVER_URL,
				},
				'/ws': {
					target: env.BACKEND_SERVER_URL,
					ws: true,
					rewriteWsOrigin: true,
				}
			}
		},
		plugins: [tailwindcss(), sveltekit()],
		define: {
			// Double JSON.stringify:
			//   outer → valid JS string literal for Vite's define substitution
			//   inner → the JSON payload that JSON.parse() reads at runtime
			'import.meta.env.VITE_DOCS': JSON.stringify(JSON.stringify(await getDocs()))
		}
	}
});
