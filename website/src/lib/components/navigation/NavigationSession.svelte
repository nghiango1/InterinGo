<script lang="ts">
	import { resolve } from '$app/paths';
	import type { DocInfo } from '$lib/type';
	import NavigationItem from './NavigationItem.svelte';

	const {
		docs,
		name
	}: {
		docs: DocInfo[];
		name: string;
	} = $props();
</script>

<ol class="block list-disc">
	<h2 class="my-4 font-bold">{name}</h2>
	{#each docs as doc}
		<NavigationItem href={resolve(`/docs/${doc.slug}`)} name={doc.title || doc.slug} />
	{/each}
</ol>
