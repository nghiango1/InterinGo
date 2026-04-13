<script lang="ts">
	import Footer from '$lib/components/Footer.svelte';
	import Navigation from '$lib/components/Navigation.svelte';

	let { data, children } = $props();
</script>

<div class="grid px-8 md:grid-cols-(--docs-grid-cols)">
	<div class="fixed top-0 left-0 z-20 h-screen overflow-y-auto md:static md:h-auto">
		<Navigation nav={data.navSessions}/>
	</div>
	<div class="overflow-auto">
		<article class="prose my-6 w-full max-w-none dark:prose-invert">
			{@render children()}
		</article>
		<Footer />
	</div>
</div>
