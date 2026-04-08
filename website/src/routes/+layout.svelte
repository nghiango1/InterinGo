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
	<title>InterinGo language</title>
	<link rel="icon" type="image/x-icon" />
	<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/7.0.1/css/all.min.css" integrity="sha512-2SwdPD6INVrV/lHTZbO2nodKhrnDdJK9/kg2XD1r9uGqPo1cUbujc+IYdlYdEErWNu69gVcYgdxlmVmzTWnetw==" crossorigin="anonymous" referrerpolicy="no-referrer" />
	<link rel="stylesheet" href="/public/prism.css" />
	<script src="/public/htmx.min.js"></script>
	<script src="/public/custom.js"></script>
	<script src="/public/prism.js"></script>
	<link rel="icon" href={favicon} />
</svelte:head>

<div id="app" class="scrollbar min-w-80 dark:bg-[#050510]">
	<main>
		<Header />

		{@render children()}
	</main>
</div>
