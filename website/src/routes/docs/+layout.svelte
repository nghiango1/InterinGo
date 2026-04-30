<script lang="ts">
	import Footer from '$lib/components/Footer.svelte';
	import Navigation from '$lib/components/Navigation.svelte';
	import { afterNavigate } from '$app/navigation';
	import { onMount } from 'svelte';

	onMount(() => {
		afterNavigate(() => {
			primeSetup();
		});
	});
	function primeSetup() {
		// @ts-ignore
		if (window.Prism) {
			// @ts-ignore
			let prism = window.Prism;
			prism.languages.iig = {
				comment: /\/\/.*/,
				number: /\b\d+\b/,
				keyword: /\b(if|else|let|return|fn)\b/,
				boolean: /\b(true|false)\b/,
				operator: /[+\-*/=<>!]+/,
				punctuation: /[{}[\];(),.:]/
			};
			prism.highlightAll();
		}
	}

	let { data, children } = $props();
</script>

<svelte:head>
	<link rel="stylesheet" href="/public/prism.css" />
	<script src="/public/prism.js" onload={primeSetup}></script>
</svelte:head>

<div class="mx-auto max-w-6xl w-full">
	<div class="grid xl:grid-cols-(--docs-grid-cols) xl:w-full">
		<div class="fixed top-0 left-0 z-20 h-screen overflow-y-auto xl:static xl:h-auto">
			<Navigation nav={data.navSessions} />
		</div>

		<div class="mx-6 overflow-auto xl:w-full">
			<article class="my-6 prose max-w-none prose-stone dark:prose-invert ">
				{@render children()}
			</article>
		</div>
	</div>
</div>
