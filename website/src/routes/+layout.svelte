<script lang="ts">
	import './layout.css';
	import Header from '$lib/components/Header.svelte';

	let { children } = $props();
</script>

<svelte:head>
	<meta charset="UTF-8" />
	<meta name="viewport" content="width=device-width, initial-scale=1.0" />
	<meta http-equiv="x-ua-compatible" content="ie=edge" />
	<title>InterinGo — Interpreter built in Go</title>
	<link rel="apple-touch-icon" sizes="180x180" href="/apple-touch-icon.png" />
	<link rel="icon" type="image/png" sizes="32x32" href="/favicon-32x32.png" />
	<link rel="icon" type="image/png" sizes="16x16" href="/favicon-16x16.png" />
	<link rel="manifest" href="/site.webmanifest" />
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
