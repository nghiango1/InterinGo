<script lang="ts">
	import CodeBlock from './CodeBlock.svelte';

	let {
		setCommand
	}: {
		setCommand: (input: string) => void;
	} = $props();
</script>

<article class="prose max-w-none dark:prose-invert">
	<h1>InterinGo</h1>
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
	<CodeBlock name={'Comparation'} code={'1 > 2'} setCommand={setCommand} />
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
	<CodeBlock name={'Complex calculation'} code={'4 * (4 / 2) * (3 + 2) + 1'} setCommand={setCommand} />
	<CodeBlock name={'Control flow'} code={'if (1 > 2) { return 10 } else { return 3 }'} setCommand={setCommand} />
	<CodeBlock name={'Variable'} code={'let x = 2 * 2 * 2; return x;'} setCommand={setCommand} />
	<q>
		Also, current REPL session is shared with only one backend, so we can reuse our variables even
		the page is reload. With above code, we currently have <code>x = 8</code>
	</q>
	<CodeBlock name={'Function'} code={'let add = fn (x,y) { x + y }; return add(4,x);'} setCommand={setCommand} />
	<p>Error throwing is here too</p>
	<CodeBlock name={'Error'} code={'let x = 2/0'} setCommand={setCommand} />
	<p>In-case use didn't see the input box placeholder just yet, try this too</p>
	<CodeBlock name={'Built-in command'} code={'help()'} setCommand={setCommand} />
</article>
