<script lang="ts">
	let { line }: { line: string } = $props();

	type LineType = 'command' | 'output' | 'error' | 'comment';

	function classify(line: string): LineType {
		if (line.startsWith('>> ')) return 'command';
		if (line.startsWith('ERROR') || line.startsWith('error')) return 'error';
		if (line.startsWith('//') || line.startsWith('#')) return 'comment';
		return 'output';
	}

	const type = $derived(classify(line));
</script>

<span
	class={[
		'block font-mono text-sm leading-relaxed',
		type === 'command' && 'dark:text-stone-100',
		type === 'output' && 'text-emerald-400',
		type === 'error' && 'text-red-400',
		type === 'comment' && 'text-stone-600 dark:text-stone-400 italic'
	]
		.filter(Boolean)
		.join(' ')}
>
	{#if type === 'command'}
		<span class="mr-2 text-stone-600 select-none">&gt;&gt;</span>{line.slice(3)}
	{:else}
		{line}
	{/if}
</span>
