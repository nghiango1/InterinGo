<script lang="ts">
	import './layout.css';
	import favicon from '$lib/assets/favicon.svg';
	import Header from '$lib/components/Header.svelte';

	let { children } = $props();
</script>

<svelte:head>
	<meta charset="UTF-8" />
	<meta name="viewport" content="width=device-width, initial-scale=1.0" />
	<meta http-equiv="x-ua-compatible" content="ie=edge" />
	<title>InterinGo — Interpreter built in Go</title>
	<link rel="icon" type="image/x-icon" />
	<link rel="stylesheet" href="/public/prism.css" />
	<script src="/public/prism.js"></script>
	<link rel="icon" href={favicon} />
</svelte:head>

<div
	id="app"
	class="min-h-screen dark:bg-stone-950 dark:text-stone-100"
	style="font-family: 'JetBrains Mono', monospace;"
>
	<main>
		<Header />

		{@render children()}
	</main>
</div>
