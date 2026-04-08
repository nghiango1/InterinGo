<script lang="ts">
	import Line from '$lib/components/command_prompt/Line.svelte';
	import Control from '$lib/components/command_prompt/Control.svelte';
	import { postEvaluate } from '$lib/controller/repl';
	import { type EvalRequest, type EvalResponseSuccess } from '$lib/server/repl';

	let { command = $bindable('') }: { command?: string } = $props();

	let isEval = $state(false);
	let stick = $state(false);
	let hide = $state(false);
	let wrap = $state(false);

	const STARTED_LINE = '// Let start with help() command';
	let lines: string[] = $state([STARTED_LINE]);

	let replOutput: HTMLElement;

	$effect(() => {
		if (lines.length && replOutput) {
			replOutput.scrollTop = replOutput.scrollHeight;
		}
	});

	function addCommand() {
		lines.push(`>> ${command}`);
	}

	function copyEvalResult(resp: EvalResponseSuccess) {
		lines.push(resp.output);
	}

	async function evaluate() {
		if (!command.trim() || isEval) return;
		isEval = true;

		const req: EvalRequest = { data: command };
		addCommand();

		const [status, resp] = await postEvaluate(req);

		if (status === 200) {
			copyEvalResult(resp as EvalResponseSuccess);
		}

		isEval = false;
	}

	function handleKeydown(e: KeyboardEvent) {
		if (e.key === 'Enter') evaluate();
	}
</script>

<div class={stick ? 'fixed' : 'sticky' + ' top-0 left-0 z-10'}>
	<div
		class="flex flex-col overflow-hidden rounded-xl border border-stone-700 bg-white shadow-2xl shadow-black/40 dark:bg-stone-900"
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
					state={wrap}
					label="wrap"
					toggle={() => {
						wrap = !wrap;
					}}
				/>
				<Control
					state={stick}
					label="stick"
					toggle={() => {
						stick = !stick;
					}}
				/>
				<Control
					state={hide}
					label="hide"
					toggle={() => {
						hide = !hide;
					}}
				/>
			</div>
		</div>

		{#if hide}
			<!-- Minimized: show only last line -->
			<div class="border-b border-stone-700/60 bg-stone-950/50 px-4 py-1.5">
				<Line line={lines[lines.length - 1]} />
			</div>
		{:else}
			<pre
				bind:this={replOutput}
				class={[
					'scrollbar-thin scrollbar-track-transparent scrollbar-thumb-stone-700 not-prose dark:bg-stone-450 m-0 h-52 resize-y overflow-auto px-4 py-3',
					wrap ? 'break-all whitespace-pre-wrap' : 'whitespace-pre'
				].join(' ')}>{#each lines as line}<Line {line} />{/each}</pre>
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
				bind:value={command}
				onkeydown={handleKeydown}
				disabled={isEval}
			/>
			<button
				class={[
					'rounded-lg border px-4 py-1.5 font-mono text-xs transition-all',
					isEval
						? 'cursor-not-allowed border-stone-700 text-stone-600'
						: 'border-stone-600 hover:bg-stone-700 hover:text-stone-200 active:scale-95 dark:bg-stone-800 dark:text-stone-200'
				].join(' ')}
				id="repl-send"
				type="button"
				disabled={isEval}
				onclick={evaluate}
			>
				{isEval ? '...' : 'Run'}
			</button>
		</div>
	</div>

	<div>
		<p class="mt-3 text-[11px] leading-relaxed dark:text-stone-300">
			Session is shared with one backend — variables persist across snippets even after reload.
		</p>
	</div>
</div>
