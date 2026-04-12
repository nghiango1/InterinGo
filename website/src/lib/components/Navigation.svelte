<script lang="ts">
	import type { NavigationRecord } from '$lib/type';
	import NavigationSession from './navigation/NavigationSession.svelte';

	let { nav }: { nav: NavigationRecord[] } = $props();
	const priority = ['Getting started', 'Syntax'];
	let sortedNav = $derived.by(() => {
		const result = nav.toSorted((a: NavigationRecord, b: NavigationRecord) => {
			const index_a = priority.findIndex((v) => v == a.name) || 100;
			const index_b = priority.findIndex((v) => v == b.name) || 100;
			return -index_a + index_b;
		});

		for (let i = 0; i < result.length; i++) {
			result[i].docs.sort((a, b) => (a.index || 100) - (b.index || 100));
		}
		return result;
	});

	let hide = $state(false);
	function handleHide() {
		hide = !hide;
	}
</script>

<input id="linked" class="peer m-auto hidden" type="checkbox" name="checked" bind:checked={hide} />
<label
	for="linked"
	class="sticky top-16 block rounded-r-lg bg-gray-200 p-1 [writing-mode:vertical-lr] peer-checked:hidden focus:bg-blue-500 focus:font-bold active:bg-blue-500 md:hidden dark:bg-gray-500 dark:text-white"
>
	Show menu
</label>
<div
	class="sticky top-4 hidden w-76 px-4 peer-checked:block before:absolute before:inset-0 before:-top-4 before:-z-10 before:rounded-b-lg before:bg-white/30 before:object-none before:backdrop-blur-lg md:block dark:text-white before:dark:bg-[#050510]/30"
>
	<div class="flex flex-row items-end py-4">
		<h1 class="block flex-1 text-xl font-bold dark:text-white">Documentation</h1>
		<button class="block md:hidden" onclick={handleHide}> Hide menu </button>
	</div>
	{#each sortedNav as session}
		<NavigationSession docs={session.docs} name={session.name} />
	{/each}
</div>
