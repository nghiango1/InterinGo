<script lang="ts">
	let { state, label, toggle }: { state: boolean; label: string; toggle: () => void } = $props();
</script>

<button
	onclick={toggle}
	class={[
		'rounded border px-2 py-0.5 font-mono text-[11px] transition-all',
		state
			? 'bg-stone-600 font-semibold text-stone-100'
			: 'text-stone-500 hover:bg-stone-700/60 hover:text-stone-300'
	].join(' ')}
>
	{label}
</button>
