<script lang="ts">
	import ShortInfo from '$lib/components/ShortInfo.svelte';
	
</script>
<footer
	class="prose flex max-w-none flex-col gap-4 border-t-2 border-t-blue-200 pt-2 sm:flex-row-reverse sm:pt-8 dark:prose-invert"
>
	<div class="flex flex-col sm:ml-auto sm:text-right">
		<h2 class="font-bold">Menu</h2>
		<a class="text-blue-500 underline" href="/">home</a>
		<a class="text-blue-500 underline" href="/info">about</a>
		<a class="text-blue-500 underline" href="/404">404</a>
	</div>

	<ShortInfo />
</footer>
