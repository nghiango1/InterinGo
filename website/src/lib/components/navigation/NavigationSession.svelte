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

<div
	class="sticky top-4 block h-full w-[19rem] px-4 before:absolute before:-inset-0 before:-top-4 before:-z-10 before:rounded-b-lg before:bg-white/30 before:object-none before:backdrop-blur-lg md:block dark:text-white before:dark:bg-[#050510]/30"
>
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
</div>
