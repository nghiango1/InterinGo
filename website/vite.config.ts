import tailwindcss from '@tailwindcss/vite';
import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';

import { compile } from 'mdsvex';
import { readFileSync } from 'fs';
import { sync } from 'glob';

interface DocRecord {
	href: string;
	name?: string;
}

async function getDocs(): Promise<DocRecord[]> {
	let docs: DocRecord[] = [];
	const pattern = 'src/lib/docs/**/*.md';
	try {
		const files = sync(pattern);
		for (let i = 0; i < files.length; i++) {
			const file = readFileSync(files[i], 'utf8');
			const mdx = await compile(file);
			const href = files[i].replace('src/lib', '').replace('.md', '');

			// Get front matter title value
			const name = (mdx?.data?.fm as Record<string, unknown> | undefined)?.title as
				| string
				| undefined;

			docs.push({
				href: href,
				name: name
			});
		}
	} catch (err) {
		console.error('Error fetching files: ', err);
	}

	return docs;
}


export default defineConfig({
	server: {
		port: 3000,
		open: true,
		host: true,
	},
	plugins: [tailwindcss(), sveltekit()],
	define: {
		// Double JSON.stringify:
		//   outer → valid JS string literal for Vite's define substitution
		//   inner → the JSON payload that JSON.parse() reads at runtime
		'import.meta.env.VITE_DOCS': JSON.stringify(
			JSON.stringify(await getDocs())
		),
	},
});
