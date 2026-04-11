<script lang="ts">
	import CommandPrompt from '../CommandPrompt.svelte';
	import CodeBlock from './CodeBlock.svelte';
	import { commandPromptState as cpState } from '$lib/components/CommandPromptState.svelte';
	import { onMount } from 'svelte';

	function setCommand(input: string) {
		cpState.command = input;
	}

	const snippets = [
		{
			label: 'Comparison',
			code: '1 > 2',
			description:
				'Basic comparison operators: As you see, the output will be "false" because 1 is less than 2'
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
			code: 'let add = fn (x,y) { return x + y };',
			description: 'First-class functions'
		},
		{
			label: 'Function call',
			code: 'add(4,x);',
			description:
				'Calling add function with x = 8 from last code block, add(4, x) = add(4, 8) = 12 is expected'
		},
		{
			label: 'Error',
			code: 'let x = 2/0',
			description: 'Division by zero error'
		},
		{
			label: 'Built-in',
			code: 'help()',
			description:
				'Built-in commands. In-case you didnt see the input box placeholder just yet, `help()` will list all built-in command'
		}
	];

	let intersecting = $state(false);
	let container: HTMLElement;

	onMount(() => {
		if (typeof IntersectionObserver !== 'undefined') {
			const rootMargin = `0px 0px 0px 0px`;

			const observer = new IntersectionObserver(
				(entries) => {
					intersecting = entries[0].isIntersecting;
				},
				{
					rootMargin
				}
			);

			observer.observe(container);
			return () => observer.unobserve(container);
		}

		function handler() {
			const bcr = container.getBoundingClientRect();
			intersecting =
				bcr.bottom > 0 &&
				bcr.right > 0 &&
				bcr.top < window.innerHeight &&
				bcr.left < window.innerWidth;

			if (intersecting) {
				window.removeEventListener('scroll', handler);
			}
		}

		window.addEventListener('scroll', handler);
		return () => window.removeEventListener('scroll', handler);
	});
</script>

<article class="mx-auto prose max-w-6xl px-6 py-16 dark:prose-invert">
	<section class="">
		<div class={'grid grid-cols-1 items-start gap-12 xl:grid-cols-2'}>
			<div class="w-full">
				<span class="mb-4 text-xs tracking-[0.2em] uppercase dark:text-stone-600">
					interpreter · built in go
				</span>

				<h1
					class="mb-6 text-5xl leading-none font-bold tracking-tight lg:text-6xl dark:text-stone-100"
				>
					Interin<span class="text-stone-600">Go</span>
				</h1>

				<p
					class="mb-8 max-w-md text-sm leading-relaxed text-stone-600 dark:text-stone-400"
					style="font-family: 'Instrument Serif', serif; font-size: 1.1rem;"
				>
					A hand-crafted interpreter language built to challenge advanced compiler and evaluator
					topics — now available in your browser via a live REPL.
				</p>

				<div class="mb-10 flex flex-wrap gap-2">
					{#each ['variables', 'functions', 'closures', 'control flow', 'error handling', 'built-ins'] as feat}
						<span
							class="rounded border border-stone-800 px-2.5 py-1 text-[11px] dark:text-stone-500"
						>
							{feat}
						</span>
					{/each}
				</div>

				<div class={'my-2 flex h-96 xl:hidden'} bind:this={container}>
					{#if intersecting}
						<CommandPrompt forceNotHide={true} />
					{/if}
				</div>
				<div class="fixed top-0 right-[10%] my-2 h-0 w-[80dvw] xl:hidden">
					{#if !intersecting}
						<CommandPrompt />
					{/if}
				</div>

				<div class="mx-auto max-w-6xl">
					<div class="border-t border-stone-800"></div>
				</div>

				<section class="mx-auto max-w-6xl py-16">
					<div class="mb-8 flex items-end justify-between">
						<div>
							<span class="mb-4 text-xs tracking-[0.2em] uppercase dark:text-stone-600">
								Examples
							</span>
							<h2 class="not-prose text-xl font-bold dark:text-stone-100">Click &amp; run</h2>
						</div>
					</div>
					<p
						class="mb-8 max-w-md text-sm leading-relaxed text-stone-600 dark:text-stone-400"
						style="font-family: 'Instrument Serif', serif; font-size: 1.1rem;"
					>
						Clicking a snippet will loads its code into the REPL input above. Try the snippets below
						in order.
					</p>

					<div class="grid grid-cols-1 gap-3">
						{#each snippets as snippet}
							<CodeBlock
								code={snippet.code}
								name={snippet.label}
								description={snippet.description}
								{setCommand}
							/>
						{/each}
					</div>
				</section>
			</div>

			<div class={'h-full not-xl:hidden pb-16 pt-2'}>
				<div class={'sticky top-0 left-0 z-10 h-[80dvh]'}>
					<CommandPrompt forceNotHide={true} />
				</div>
			</div>
		</div>
	</section>

	<div class="mx-auto max-w-6xl px-6">
		<div class="border-t border-stone-800"></div>
	</div>

	<section class="mx-auto max-w-6xl px-6 py-16">
		<div class="mb-8 flex items-end justify-between">
			<div>
				<span class="mb-4 text-xs tracking-[0.2em] uppercase dark:text-stone-600"> About </span>
				<h2 class="not-prose text-xl font-bold dark:text-stone-100">Why InterinGo?</h2>
			</div>
		</div>
		<div class="grid grid-cols-1 gap-12 lg:grid-cols-2">
			<div>
				<p
					class="mb-8 max-w-md text-sm leading-relaxed text-stone-600 dark:text-stone-400"
					style="font-family: 'Instrument Serif', serif; font-size: 1.1rem;"
				>
					It hard! To chalenge my self and prove that I can handle more advanged topics. Building an
					interpreter from scratch in Go forces me to deeply understand lexing, parsing, AST
					construction, and tree-walking evaluation. Also this was a course at my University.
				</p>
			</div>

			<div
				class="flex h-fit flex-col justify-center rounded-xl border border-stone-800 bg-stone-900 p-8"
			>
				<span class="mb-4 text-xs tracking-[0.2em] text-stone-400 uppercase">Offline use</span>
				<h3 class="not-prose mb-3 text-lg font-bold text-stone-100">Download the REPL binary</h3>
				<a
					href="/download"
					class="flex w-fit items-center gap-2 rounded-lg border border-stone-600 bg-stone-800 px-5 py-2.5 text-sm text-stone-200 transition-all hover:border-stone-400 hover:bg-stone-700 hover:text-white active:scale-95"
				>
					<span>↓</span>
					<span>Download binary</span>
				</a>
			</div>
		</div>
		<p
			class="mb-8 text-sm leading-relaxed text-stone-600 dark:text-stone-400"
			style="font-family: 'Instrument Serif', serif; font-size: 1.1rem;"
		>
			Of course, these examples isn't all of what InterinGo can do, but this home page is just for
			showing InterinGo REPL to the world, let have some readding in <a href="/docs">document</a>.
		</p>
	</section>
</article>
