<script lang="ts">
	import { commandPromptState } from '../CommandPromptState.svelte';

	let {
		id,
		name,
		code,
		description,
		setCommand
	}: {
		id: string;
		name: string;
		code: string;
		description: string;
		setCommand: (input: string) => void;
	} = $props();
</script>

<button
	{id}
	onclick={() => setCommand(code)}
	class={'xl:scroll-mt-10 group flex flex-col overflow-hidden rounded-xl border border-stone-800 text-left transition-all hover:border-stone-600 hover:bg-stone-800/80 hover:text-stone-300 active:scale-[0.98] dark:bg-stone-900' +
		' ' +
		(commandPromptState.hide ? 'scroll-mt-40' : 'scroll-mt-96')}
>
	<!-- scroll-mt help with extra space for fixed element, this case we have the command prompt, other than Header -->
	<div
		class="flex items-center justify-between border-b border-stone-800 px-3 py-2 transition-colors group-hover:border-stone-700 dark:bg-stone-800"
	>
		<span
			class="text-[11px] font-semibold tracking-widest uppercase group-hover:text-stone-300 dark:text-stone-400"
		>
			{name}
		</span>
		<span class="text-[10px] text-stone-500 transition-colors group-hover:text-stone-300">
			run ↵
		</span>
	</div>
	<div class="flex-1 px-3 py-3">
		<pre class="not-prose block leading-relaxed dark:text-stone-300">{code}</pre>
	</div>
	<div class="border-t border-stone-800/60 px-3 py-2 group-hover:border-stone-700">
		<span
			class="text-xs text-stone-600 transition-colors group-hover:text-stone-300 dark:text-stone-500"
			style="font-family: 'Instrument Serif', serif;"
		>
			{description}
		</span>
	</div>
</button>
