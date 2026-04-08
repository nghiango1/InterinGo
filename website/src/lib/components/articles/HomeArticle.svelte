<script lang="ts">
	import CommandPrompt from '../CommandPrompt.svelte';
	import CodeBlock from './CodeBlock.svelte';

	let command = $state('');
	function setCommand(input: string) {
		command = input;
	}

	const snippets = [
		{
			label: 'Comparison',
			code: '1 > 2',
			description: 'Basic comparison operators'
		},
		{
			label: 'Calculation',
			code: '4 * (4 / 2) * (3 + 2) + 1',
			description: 'Nested arithmetic expressions'
		},
		{
			label: 'Control flow',
			code: 'if (1 > 2) { return 10 } else { return 3 }',
			description: 'if / else branching'
		},
		{
			label: 'Variable',
			code: 'let x = 2 * 2 * 2; return x;',
			description: 'let bindings & return'
		},
		{
			label: 'Function',
			code: 'let add = fn (x,y) { x + y }; return add(4,x);',
			description: 'First-class functions'
		},
		{
			label: 'Error',
			code: 'let x = 2/0',
			description: 'Division by zero error'
		},
		{
			label: 'Built-in',
			code: 'help()',
			description: 'Built-in commands'
		}
	];
</script>

<article class="prose max-w-none dark:prose-invert">
	<section class="mx-auto max-w-6xl px-6 py-16">
		<div class="grid grid-cols-1 items-start gap-12 xl:grid-cols-2">
			<div>
				<p class="mb-4 text-xs tracking-[0.2em] uppercase dark:text-stone-600">
					interpreter · built in go
				</p>

				<h1
					class="mb-6 text-5xl leading-none font-bold tracking-tight lg:text-6xl dark:text-stone-100"
				>
					Interin<span class="text-stone-600">Go</span>
				</h1>

				<p
					class="mb-8 max-w-md text-sm leading-relaxed text-stone-400"
					style="font-family: 'Instrument Serif', serif; font-size: 1.1rem;"
				>
					A hand-crafted interpreter language built to challenge advanced compiler and evaluator
					topics — now available in your browser via a live REPL.
				</p>

				<div class="mb-10 flex flex-wrap gap-2">
					{#each ['variables', 'functions', 'closures', 'control flow', 'error handling', 'built-ins'] as feat}
						<span class="rounded border border-stone-800 px-2.5 py-1 text-[11px] text-stone-500">
							{feat}
						</span>
					{/each}
				</div>

				<!-- ── DIVIDER ── -->
				<div class="mx-auto max-w-6xl px-6">
					<div class="border-t border-stone-800"></div>
				</div>

				<!-- ── SNIPPETS ── -->
				<section class="mx-auto max-w-6xl px-6 py-16">
					<div class="mb-8 flex items-end justify-between">
						<div>
							<p class="mb-1 text-[10px] tracking-[0.2em] text-stone-600 uppercase">
								Examples - Try the snippets below in order.
							</p>
							<h2 class="text-xl font-bold dark:text-stone-100">Click &amp; run</h2>
						</div>
						<p class="max-w-xs text-right text-xs leading-relaxed text-stone-600">
							Clicking a snippet loads its code into the REPL input above.
						</p>
					</div>

					<div class="grid grid-cols-1 gap-3 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
						{#each snippets as snippet}
							<CodeBlock
								code={snippet.code}
								name={snippet.label}
								description={snippet.description}
								{setCommand}
							/>
						{/each}
					</div>

					<!-- Inline note about shared state -->
					<div class="mt-6 rounded-lg border border-stone-800/60 bg-stone-900/40 px-4 py-3">
						<p class="text-xs leading-relaxed text-stone-500">
							<span class="font-semibold text-stone-400">Note:</span>
							After running the Variable snippet, <code class="text-emerald-500">x = 8</code> is
							bound in the session. The Function snippet uses it:
							<code class="text-emerald-500">add(4, x)</code>
							→
							<code class="text-emerald-500">12</code>.
						</p>
					</div>
				</section>
			</div>

				<CommandPrompt bind:command />
		</div>
	</section>

	<!-- ── DIVIDER ── -->
	<div class="mx-auto max-w-6xl px-6">
		<div class="border-t border-stone-800"></div>
	</div>

	<!-- ── ABOUT + DOWNLOAD ── -->
	<section class="mx-auto max-w-6xl px-6 py-16">
		<div class="grid grid-cols-1 gap-12 lg:grid-cols-2">
			<!-- About -->
			<div>
				<p class="mb-1 text-[10px] tracking-[0.2em] text-stone-600 uppercase">About</p>
				<h2 class="mb-4 text-xl font-bold text-stone-100">Why InterinGo?</h2>
				<p
					class="mb-4 text-sm leading-relaxed text-stone-400"
					style="font-family: 'Instrument Serif', serif; font-size: 1rem;"
				>
					Building an interpreter from scratch in Go forces you to deeply understand lexing,
					parsing, AST construction, and tree-walking evaluation. This is that journey, made public.
				</p>
				<p
					class="text-sm leading-relaxed text-stone-500"
					style="font-family: 'Instrument Serif', serif; font-size: 1rem;"
				>
					The online REPL is a live preview against a single shared backend. For the full experience
					and offline use, download the binary.
				</p>
			</div>

			<!-- Download CTA -->
			<div class="flex flex-col justify-center rounded-xl border border-stone-800 bg-stone-900 p-8">
				<p class="mb-2 text-[10px] tracking-[0.2em] text-stone-600 uppercase">Offline use</p>
				<h3 class="mb-3 text-lg font-bold text-stone-100">Download the REPL binary</h3>
				<p class="mb-6 text-xs leading-relaxed text-stone-500">
					The online REPL won't be available forever. Download the build and run InterinGo locally
					with full language support.
				</p>
				<a
					href="/download"
					class="flex w-fit items-center gap-2 rounded-lg border border-stone-600 bg-stone-800 px-5 py-2.5 text-sm text-stone-200 transition-all hover:border-stone-400 hover:bg-stone-700 hover:text-white active:scale-95"
				>
					<span>↓</span>
					<span>Download binary</span>
				</a>
			</div>
		</div>
	</section>

	<p>
		"interprester-in-go" or InterinGo (for short) is a new interpreter language by <a
			href="https://www.linkedin.com/in/nghia-ngo-duc">me</a
		> to chalenge my self with more advanged topics. The command-prompt window is the directed way to
		interact with InterinGo REPL, which have been broughted to web so you can try it now without the need
		to download any binary.
	</p>
	<p>
		To make sure you not get lost with how to use the language, here is a sample craft code to try
		out in REPL command prompt. Click on Click and run to copy the code into Command prompt and have
		REPL run the code. The input box will reflect the code being used, while the command prompt will
		show the evaluation result.
	</p>
	<CodeBlock
		name={'Comparation'}
		code={'1 > 2'}
		description={'Basic comparison operators'}
		{setCommand}
	/>
	<p>
		As you see, the output will be <code>false</code> because 1 is less than 2. Also, the command-prompt
		can be annoy in smaller screen, try using Hide checkbox to minimize it, don't be too worry, it still
		show evaluation result in minimized state (or uncheck sticky box eh... wanna use that?).
	</p>
	<h2>Examples</h2>
	<p>
		Of course, these examples isn't all of what InterinGo can do, but this home page is just for
		showing InterinGo REPL to the world, let have some readding in <a href="/docs">document</a> to
		learning more about it. The online online REPL won't be there thought, so let head on and
		downloading REPL build file <a href="https://github.com/nghiango1/hello/releases">here </a>
	</p>
	<CodeBlock
		name={'Complex calculation'}
		code={'4 * (4 / 2) * (3 + 2) + 1'}
		description={'Nested arithmetic expressions'}
		{setCommand}
	/>
	<CodeBlock
		name={'Control flow'}
		code={'if (1 > 2) { return 10 } else { return 3 }'}
		description={'if / else branching'}
		{setCommand}
	/>
	<CodeBlock
		name={'Variable'}
		code={'let x = 2 * 2 * 2; return x;'}
		description={'let bindings & return'}
		{setCommand}
	/>
	<q>
		Also, current REPL session is shared with only one backend, so we can reuse our variables even
		the page is reload. With above code, we currently have <code>x = 8</code>
	</q>
	<CodeBlock
		name={'Function'}
		code={'let add = fn (x,y) { x + y }; return add(4,x);'}
		description={'First-class functions'}
		{setCommand}
	/>
	<p>Error throwing is here too</p>
	<CodeBlock
		name={'Error'}
		code={'let x = 2/0'}
		description={'Division by zero error'}
		{setCommand}
	/>
	<p>In-case use didn't see the input box placeholder just yet, try this too</p>
	<CodeBlock
		name={'Built-in command'}
		code={'help()'}
		description={'Built-in commands'}
		{setCommand}
	/>
</article>
