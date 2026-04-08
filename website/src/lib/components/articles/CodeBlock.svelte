<script lang="ts">
	let {
		name,
		code,
		description,
		setCommand
	}: {
		name: string;
		code: string;
		description: string;
		setCommand: (input: string) => void;
	} = $props();
</script>

<button
	onclick={() => setCommand(code)}
	class="group flex flex-col overflow-hidden rounded-xl border border-stone-800 dark:bg-stone-900 text-left transition-all hover:border-stone-600 hover:bg-stone-800/80 active:scale-[0.98]"
>
	<!-- Snippet header -->
	<div
		class="flex items-center justify-between border-b border-stone-800 dark:bg-stone-800 px-3 py-2 transition-colors group-hover:border-stone-700"
	>
		<span
			class="text-[11px] font-semibold tracking-widest dark:text-stone-400 uppercase group-hover:text-stone-300"
		>
			{name}
		</span>
		<span class="text-[10px] text-stone-700 transition-colors group-hover:text-stone-500">
			run ↵
		</span>
	</div>
	<!-- Code -->
	<div class="flex-1 px-3 py-3">
		<code class="block text-xs leading-relaxed dark:text-stone-300">
			{code}
		</code>
	</div>
	<!-- Description -->
	<div class="border-t border-stone-800/60 px-3 py-2">
		<span
			class="text-[11px] text-stone-600 transition-colors group-hover:text-stone-500"
			style="font-family: 'Instrument Serif', serif;"
		>
			{description}
		</span>
	</div>
</button>
