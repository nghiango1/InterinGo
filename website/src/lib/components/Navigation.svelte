<script lang="ts">
	let { docs }: { docs: string[] } = $props();
	let hide = $state(false);
	function handleHide() {
		hide = !hide;
	}
</script>

{#if !hide}
	<button
		class={'vertical-lr sticky top-16 block rounded-r-lg bg-gray-200 p-1 focus:bg-blue-500 focus:font-bold active:bg-blue-500 md:hidden dark:bg-gray-500 dark:text-white'}
		onclick={handleHide}
	>
		Show menu
	</button>
{:else}
	<div
		class="sticky top-4 block h-full w-[19rem] px-4 before:absolute before:-inset-0 before:-top-4 before:-z-10 before:rounded-b-lg before:bg-white/30 before:object-none before:backdrop-blur-lg md:block dark:text-white before:dark:bg-[#050510]/30"
	>
		<ol class="block list-disc">
			<div class="flex flex-row items-end py-4">
				<h1 class="block flex-1 text-xl font-bold dark:text-white">Documentation</h1>
				<button class="block md:hidden" onclick={handleHide}> Hide menu </button>
			</div>
			<h2 class="my-4 font-bold">Session</h2>
			{#each docs as doc}
				<li
					class="block border-l-2 border-gray-200 py-2 pl-4 hover:border-blue-500 hover:font-bold focus:border-blue-500 focus:font-bold"
				>
					<a href="/">{doc}</a>
				</li>
			{/each}
		</ol>
	</div>
{/if}
