<script lang="ts">
	let {
		name,
		code,
		setCommand
	}: {
		name: string;
		code: string;
		setCommand: (input: string) => void;
	} = $props();

	function handleClick() {
		setCommand(code);
	}
</script>

<div class="flex flex-col">
	<div class="flex flex-row">
		<p class="my-auto block flex-1">Snippet - <code>{name}</code></p>
		<button class="my-auto block rounded-lg underline" {name} onclick={handleClick}>
			Copy and Run
		</button>
	</div>
	<pre id={name}>{code}</pre>
</div>
