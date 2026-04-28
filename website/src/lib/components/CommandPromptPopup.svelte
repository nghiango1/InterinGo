<script>
	import CommandPrompt from './CommandPrompt.svelte';
	let open = true;
</script>

{#if !open}
	<button
		on:click={() => (open = true)}
		class="fixed right-4 bottom-4 z-50 rounded-lg border border-stone-600 bg-stone-800 px-4 py-2 text-sm text-stone-200 shadow-lg hover:bg-stone-700"
	>
		InterinGo REPL
	</button>
{/if}

<div
	class={`fixed right-0 bottom-0 left-0 z-40 transition-transform duration-300 ${
		open ? 'translate-y-0' : 'translate-y-full'
	}`}
>
	<div class="mx-auto w-full max-w-5xl p-4">
		<div
			class="overflow-hidden rounded-xl border border-stone-700 bg-white shadow-2xl shadow-black/40 dark:bg-stone-900"
		>
			<div
				class="flex items-center justify-between border-b border-stone-700 px-4 py-2.5 dark:bg-stone-800"
			>
				<span class="font-mono text-xs text-stone-400"> interingo — repl v0.1 </span>

				<button
					on:click={() => (open = false)}
					class="rounded px-2 py-1 text-xs text-stone-400 hover:bg-stone-700 hover:text-stone-200"
				>
					✕
				</button>
			</div>

			<CommandPrompt hideHeader={true} />
		</div>
	</div>
</div>
