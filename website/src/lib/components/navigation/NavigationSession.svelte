<script lang="ts">
	import type { DocRecord } from '$lib/type';
	import NavigationItem from './NavigationItem.svelte';

	const {
		docs,
		name,
		handleHide
	}: {
		docs: DocRecord[];
		name: string;
		handleHide: () => void;
	} = $props();
</script>

<ol class="block list-disc">
	<div class="flex flex-row items-end py-4">
		<h1 class="block flex-1 text-xl font-bold dark:text-white">Documentation</h1>
		<button class="block md:hidden" onclick={handleHide}> Hide menu </button>
	</div>
	<h2 class="my-4 font-bold">{name}</h2>
	{#each docs as doc}
		<NavigationItem href={doc.href} name={doc.name} />
	{/each}
</ol>
