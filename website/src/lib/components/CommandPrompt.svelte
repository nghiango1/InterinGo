<script lang="ts">
	import Line from '$lib/components/command_prompt/Line.svelte';
	import Control from '$lib/components/command_prompt/Control.svelte';
	import { postEvaluate } from '$lib/controller/repl';
	import { type EvalRequest, type EvalResponseSuccess } from '$lib/server/repl';
	import { commandPromptState as state, connect } from '$lib/components/CommandPromptState.svelte';

	let { forceNotHide = false }: { forceNotHide?: boolean } = $props();
	// svelte-ignore non_reactive_update
	let replOutput: HTMLElement;

	import { onMount } from 'svelte';
	onMount(() => {
		connect();
	});

	$effect(() => {
		state.hide; // Turn hide on and off should also have scroll effect
		if (state.lines.length && replOutput) {
			replOutput.scrollTop = replOutput.scrollHeight;
		}
	});

	function addCommand() {
		state.lines.push(`>> ${state.command}`);
	}

	function copyEvalResult(resp: EvalResponseSuccess) {
		if (resp.output != null) state.lines.push(resp.output);
	}

	async function evaluate() {
		if (!state.command.trim() || state.isEval) return;
		state.isEval = true;

		const req: EvalRequest = { data: state.command };
		addCommand();

		const [status, resp] = await postEvaluate(req);

		if (status === 200) {
			copyEvalResult(resp as EvalResponseSuccess);
		}

		state.isEval = false;
	}

	function handleKeydown(e: KeyboardEvent) {
		if (e.key === 'Enter') evaluate();
	}
</script>

<div
	class={'flex flex-1 flex-col overflow-hidden rounded-xl border border-stone-700 bg-white shadow-2xl shadow-black/40 dark:bg-stone-900' +
		' ' +
		(state.hide && !forceNotHide ? '' : 'h-full min-h-72')}
>
	<div class="flex items-center gap-3 border-b border-stone-700 px-4 py-2.5 dark:bg-stone-800">
		<div class="flex items-center gap-1.5">
			<div class="h-3 w-3 rounded-full bg-red-500/80"></div>
			<div class="h-3 w-3 rounded-full bg-yellow-500/80"></div>
			<div class="h-3 w-3 rounded-full bg-green-500/80"></div>
		</div>

		<span class="flex-1 truncate font-mono text-xs tracking-wide dark:text-stone-400"
			>interingo — repl v0.1</span
		>

		<div class="flex items-center gap-1">
			<Control
				state={state.wrap}
				label={state.wrap ? 'normal' : 'wrap'}
				toggle={() => {
					state.wrap = !state.wrap;
				}}
			/>
			{#if !forceNotHide}
				<Control
					state={state.hide}
					label={state.hide ? 'expand' : 'hide'}
					toggle={() => {
						state.hide = !state.hide;
					}}
				/>
			{/if}
		</div>
	</div>

	{#if state.hide && !forceNotHide}
		<!-- Minimized: show only last line -->
		<div class="bg-white px-4 py-1.5 dark:bg-stone-900">
			<Line line={state.lines[state.lines.length - 1]} />
		</div>
	{:else}
		<pre
			bind:this={replOutput}
			class={[
				'scrollbar-thin scrollbar-track-transparent scrollbar-thumb-stone-700 not-prose m-0 flex-1 resize-y overflow-auto bg-white px-4 py-3 dark:bg-stone-900 ',
				state.wrap ? 'break-all whitespace-pre-wrap' : 'whitespace-pre'
			].join(' ')}>{#each state.lines as line}<Line {line} />{/each}</pre>
	{/if}

	<!-- Input row -->
	<div class="flex items-center gap-3 border-t border-stone-700 px-4 py-2.5 dark:bg-stone-800">
		<label for="repl-input" class="hidden font-mono text-xs sm:block dark:text-stone-500">
			Custom command:
		</label>
		<input
			id="repl-input"
			class="flex-1 border-b border-stone-700 bg-transparent font-mono text-sm placeholder-stone-600 transition-colors focus:border-stone-400 focus:outline-none dark:text-stone-100"
			type="text"
			name="repl-input"
			autocomplete="off"
			autocorrect="off"
			autocapitalize="off"
			enterkeyhint="go"
			spellcheck={false}
			placeholder="help()"
			bind:value={state.command}
			onkeydown={handleKeydown}
			disabled={state.isEval}
		/>
		<button
			class={[
				'rounded-lg border px-4 py-1.5 font-mono text-xs transition-all',
				state.isEval
					? 'cursor-not-allowed border-stone-700 text-stone-600 dark:text-stone-300'
					: 'border-stone-600 hover:bg-stone-700 hover:text-stone-200 active:scale-95 dark:bg-stone-800 dark:text-stone-200'
			].join(' ')}
			id="repl-send"
			type="button"
			disabled={state.isEval}
			onclick={evaluate}
		>
			{state.isEval ? '...' : 'Run'}
		</button>
	</div>
</div>
