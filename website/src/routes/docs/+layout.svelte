<script lang="ts">
	import Footer from '$lib/components/Footer.svelte';
	import Navigation from '$lib/components/Navigation.svelte';

	let { data, children } = $props();
</script>

<div class="grid px-8 md:grid-cols-(--docs-grid-cols)">
	<div class="fixed top-0 bottom-0 left-0 z-20 min-w-8 overflow-y-auto md:static">
		<Navigation nav={data.navSessions}/>
	</div>
	<div class="overflow-auto">
		<article class="prose mb-6 w-full max-w-none dark:prose-invert">
			{@render children()}
		</article>
		<Footer />
	</div>
</div>
