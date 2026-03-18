<script lang="ts">
	import type { Component } from 'svelte';
	import Footer from './Footer.svelte';
	import Navigation from './Navigation.svelte';

	let {
		SvxArticle,
		docs
	}: {
		SvxArticle: Component;
		docs: string[];
	} = $props();
</script>

<div class="grid px-8 md:grid-cols-(--docs-grid-cols)">
	<div class="fixed top-0 bottom-0 left-0 z-20 min-w-[2rem] overflow-y-auto md:static">
		<Navigation {docs} />
	</div>
	<div class="overflow-auto">
		<article class="prose mb-6 w-full max-w-none dark:prose-invert">
			<SvxArticle />
		</article>
		<Footer />
	</div>
</div>
