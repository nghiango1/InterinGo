<script lang="ts">
	import Footer from '$lib/components/Footer.svelte';
	import HomeArticle from '$lib/components/articles/HomeArticle.svelte';
</script>

<main class="flex max-w-3xl flex-col gap-4 p-4 sm:mx-8 md:mx-auto xl:mx-8 xl:max-w-full">
	<div class="flex flex-col gap-4 xl:flex-row xl:gap-8">
		<!-- Right: REPL -->
		<div class="flex-1">
			<HomeArticle/>
		</div>
	</div>
	<Footer />
</main>
