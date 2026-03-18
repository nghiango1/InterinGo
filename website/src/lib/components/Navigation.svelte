<script lang="ts">
	import type { DocRecord } from '$lib/type';
	import NavigationSession from './navigation/NavigationSession.svelte';

	let docs: DocRecord[] = $state([]);
	const raw = import.meta.env.VITE_DOCS;
	if (!raw) {
		console.warn('[docs] VITE_DOCS_MANIFEST not set – did you run the build?');
	}
	try {
		docs = JSON.parse(raw);
	} catch {
		console.error('[docs] Failed to parse VITE_DOCS_MANIFEST');
	}

	const name = 'Getting Started';
	let hide = $state(false);
	function handleHide() {
		hide = !hide;
	}
</script>

{#if !hide}
	<button
		class={'sticky top-16 block rounded-r-lg bg-gray-200 p-1 [writing-mode:vertical-lr] focus:bg-blue-500 focus:font-bold active:bg-blue-500 md:hidden dark:bg-gray-500 dark:text-white'}
		onclick={handleHide}
	>
		Show menu
	</button>
{:else}
	<NavigationSession {docs} {name} {handleHide} />
{/if}
