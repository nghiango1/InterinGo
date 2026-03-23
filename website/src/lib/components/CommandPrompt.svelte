<script lang="ts">
	import Line from './command_prompt/Line.svelte';
	import { postEvaluate } from '$lib/controller/repl';
	import { type EvalRequest, type EvalResponseSuccess } from '$lib/server/repl';

	let {
		command
	} : {
		command : string
	} = $props()

	let isEval = $state(false);
	let stick = $state(false);
	let hide = $state(false);
	// For even more freeform use tailwind break-all 
	let wrap = $state(false);

	const STARTED_LINE = 'Let start with help() command';
	let lines: string[] = $state([STARTED_LINE]);

	function updateStick() {
		stick != stick;
	}

	// Bind to the command prompt output element
	// svelte-ignore non_reactive_update
	let replOutput: HTMLElement;

	// Which we have to use effect until I found better solution
	$effect(() => {
		// Access lines.length to create a reactive dependency
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
		isEval = true;
		const req: EvalRequest = {
			data: command
		};
		addCommand();

		const [status, resp] = await postEvaluate(req);

		if (status == 200) {
			copyEvalResult((resp as EvalResponseSuccess));
			// command = '';
		}
		isEval = false;
	}
</script>

<input type="checkbox" id="hiddenstick" class="peer/hiddenstick hidden" bind:checked={stick} />
<div class={(stick ? 'sticky' : 'top-4') + 'block xl:w-5/12'}>
	<div
		class="sticky top-4 flex flex-col gap-1 before:absolute before:-inset-4 before:-top-4 before:-z-10 before:rounded-b-lg before:bg-white/30 before:object-none before:backdrop-blur-sm after:absolute after:-inset-2 after:-top-2 after:-z-10 after:rounded-b-lg after:bg-white/30 after:object-none after:blur-sm dark:text-[#d1d5db] before:dark:bg-[#050510]/30 after:dark:bg-[#050510]/30"
	>
		<div class="overflow-hidden rounded-t-lg border-2 border-blue-500 dark:border-white">
			<div class="flex gap-2 border-b-2 bg-blue-200 p-1 dark:bg-[#090d1a]">
				<h2 class="m-auto block flex-1 overflow-clip whitespace-nowrap">Command-prompt window</h2>
				<input id="wrap" class="peer/wrap m-auto" type="checkbox" name="wrap" bind:checked={wrap} />
				<label
					for="wrap"
					class="m-auto flex flex-row gap-2 rounded-lg object-none peer-checked:bg-gray-200 peer-checked/wrap:font-bold"
				>
					Wrap
				</label>
				<input
					id="stick"
					class="peer/stick m-auto"
					type="checkbox"
					name="stick"
					oninput={updateStick}
					checked={true}
				/>
				<label
					for="stick"
					class="m-auto flex flex-row gap-2 rounded-lg peer-checked:bg-gray-200 peer-checked/stick:font-bold"
				>
					Stick
				</label>
				<input id="hide" class="m-auto" type="checkbox" name="hide" bind:checked={hide} />
				<label
					for="hide"
					class="m-auto flex flex-row gap-2 rounded-lg peer-checked:bg-gray-200 peer-checked/hide:font-bold"
				>
					Hide
				</label>
			</div>
			{#if hide}
				<pre id="repl-result" class="block max-h-6 overflow-auto px-2"><Line line={lines[lines.length -1]} /></pre>
			{:else}
				<pre
					id="repl-output"
					bind:this={replOutput}
					class={'scrollbar flex h-56 resize-y overflow-auto rounded-lg p-2 whitespace-pre outline-blue-200 ' +
						(wrap ? 'whitespace-pre-wrap' : '')}>{#each lines as line}<Line {line} />{/each}</pre>
			{/if}
		</div>
		<form class="flex h-fit flex-row gap-4 outline-blue-200 dark:outline-white" method="POST">
			<label for="repl-input" class="hidden sm:my-auto sm:block">Custom command:</label>
			<input
				class="my-auto flex-1 border-b border-gray-500 bg-transparent pl-2 font-mono focus:outline-none"
				type="text"
				name="repl-input"
				aria-autocomplete="both"
				aria-labelledby="docsearch-label"
				id="repl-input"
				autocomplete="off"
				autocorrect="off"
				autocapitalize="off"
				enterkeyhint="go"
				spellcheck="false"
				placeholder="help()"
				bind:value={command}
			/>
			<button
				class="activate:dark:bg-blue-900 my-auto block rounded-lg border-2 bg-blue-200 p-1 active:bg-blue-200 dark:bg-[#090d1a]"
				id="repl-send"
				type="submit"
				disabled={isEval}
				onclick={evaluate}
			>
				Run
			</button>
		</form>
	</div>
</div>
