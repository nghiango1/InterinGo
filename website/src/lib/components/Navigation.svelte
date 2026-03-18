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

<input id="linked" class="peer m-auto hidden" type="checkbox" name="checked" bind:checked={hide} />
<label
	for="linked"
	class="sticky top-16 block rounded-r-lg bg-gray-200 p-1 [writing-mode:vertical-lr] peer-checked:hidden focus:bg-blue-500 focus:font-bold active:bg-blue-500 md:hidden dark:bg-gray-500 dark:text-white"
>
	Show menu
</label>
<div
	class="sticky top-4 hidden h-full w-[19rem] px-4 peer-checked:block before:absolute before:-inset-0 before:-top-4 before:-z-10 before:rounded-b-lg before:bg-white/30 before:object-none before:backdrop-blur-lg md:block dark:text-white before:dark:bg-[#050510]/30"
>
	<NavigationSession {docs} {name} {handleHide} />
</div>
