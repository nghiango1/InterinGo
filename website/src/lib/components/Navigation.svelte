<script lang="ts">
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import type { NavigationRecord } from '$lib/type';
	let currentPath = $derived(page.url.pathname);

	let { nav }: { nav: NavigationRecord[] } = $props();
	const priority = ['Getting started', 'Syntax'];

	let sortedNav = $derived.by(() => {
		const result = nav.toSorted((a: NavigationRecord, b: NavigationRecord) => {
			const index_a = priority.findIndex((v) => v == a.label) || 100;
			const index_b = priority.findIndex((v) => v == b.label) || 100;
			return -index_a + index_b;
		});

		for (let i = 0; i < result.length; i++) {
			result[i].items.sort((a, b) => (a.index || 100) - (b.index || 100));
			for (let j = 0; j < result[i].items.length; j++) {
				result[i].items[j].href = resolve(`/docs/[slug]`, { slug: result[i].items[j].slug });
			}
		}
		return result;
	});

	let menuOpen = $state(false);
</script>

<button
	class="sticky top-16 mb-6 items-center gap-2 rounded-r-lg border-x border-r border-stone-600 bg-stone-100 p-1 px-3 py-1.5 [writing-mode:vertical-lr] hover:bg-stone-200 md:hidden dark:bg-stone-900 dark:text-stone-400 dark:hover:bg-stone-800"
	onclick={() => (menuOpen = true)}>Show menu</button
>
<div>
	<aside
		class="fixed top-0 left-0 z-40 h-full w-68 shrink-0 overflow-y-auto border-stone-800 pt-17 pb-8 transition-transform duration-200 not-md:bg-stone-50 md:sticky md:top-13 md:translate-x-0 md:pt-8 dark:bg-stone-950"
		class:-translate-x-full={!menuOpen}
		class:translate-x-0={menuOpen}
	>
		<!-- mobile close button -->
		<button
			class="absolute top-4 right-4 rounded border border-stone-700 px-2.5 py-1 text-stone-500 transition-colors hover:text-stone-300 md:hidden dark:text-stone-300"
			onclick={() => (menuOpen = false)}>✕ close</button
		>

		{#each sortedNav as group}
			<div class="mb-7">
				<span
					class="mb-2 block px-5 tracking-[0.2em] text-stone-400 uppercase not-dark:text-stone-600"
				>
					{group.label}
				</span>
				{#each group.items as item}
					<a
						href={item.href}
						class="block border-transparent px-5 py-1.5 text-stone-600 transition-all hover:border-stone-600 hover:bg-stone-200 dark:text-stone-200 dark:hover:bg-stone-600 dark:hover:text-stone-300"
						class:border-blue-500={currentPath === item.href}
						class:bg-blue-500={currentPath === item.href}
						class:text-stone-100={currentPath === item.href}>{item.label}</a
					>
				{/each}
			</div>
		{/each}
	</aside>
</div>
