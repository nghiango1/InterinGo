<script lang="ts">
	import Footer from '$lib/components/Footer.svelte';
</script>

<main class="flex max-w-3xl flex-col gap-4 p-4 sm:mx-8 md:mx-auto xl:mx-8 xl:max-w-full">
	<div class="flex flex-col gap-4 xl:flex-row xl:gap-8">
		<input type="checkbox" id="hiddenstick" class="peer/hiddenstick hidden" checked="true" />
		<div class="block peer-checked/hiddenstick:sticky peer-checked/hiddenstick:top-4 xl:w-5/12">
			<div
				class="sticky top-4 flex flex-col gap-1 before:absolute before:-inset-4 before:-top-4 before:-z-10 before:rounded-b-lg before:bg-white/30 before:object-none before:backdrop-blur-sm after:absolute after:-inset-2 after:-top-2 after:-z-10 after:rounded-b-lg after:bg-white/30 after:object-none after:blur-sm dark:text-[#d1d5db] before:dark:bg-[#050510]/30 after:dark:bg-[#050510]/30"
			>
				<div class="overflow-hidden rounded-t-lg border-2 border-blue-500 dark:border-white">
					<input id="hiddenwrap" class="peer/wrap hidden" type="checkbox" name="wrap" />
					<input id="hiddenhide" class="peer/hide hidden" type="checkbox" name="hide" />
					<div class="flex gap-2 border-b-2 bg-blue-200 p-1 dark:bg-[#090d1a]">
						<h2 class="m-auto block flex-1 overflow-clip whitespace-nowrap">
							Command-prompt window
						</h2>
						<input
							id="wrap"
							class="peer/wrap m-auto"
							type="checkbox"
							name="wrap"
							oninput={"updateWrap(); scrollBottom();"}
						/>
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
							name="hide"
							oninput={"updateStick();"}
							checked="true"
						/>
						<label
							for="stick"
							class="m-auto flex flex-row gap-2 rounded-lg peer-checked:bg-gray-200 peer-checked/stick:font-bold"
						>
							Stick
						</label>
						<input
							id="hide"
							class="peer/hide m-auto"
							type="checkbox"
							name="hide"
							oninput={"updateHide(); scrollBottom();"}
						/>
						<label
							for="hide"
							class="m-auto flex flex-row gap-2 rounded-lg peer-checked:bg-gray-200 peer-checked/hide:font-bold"
						>
							Hide
						</label>
					</div>
					<pre
						id="repl-output"
						class="scrollbar flex h-56 resize-y overflow-auto rounded-lg p-2 whitespace-pre outline-blue-200 peer-checked/hide:hidden peer-checked/wrap:whitespace-pre-wrap"
						hx-on:htmx:after-swap="scrollBottom()"> Let start with help() command {'\n'} </pre>
					<pre
						id="repl-result"
						class="hidden max-h-6 overflow-auto px-2 peer-checked/hide:block"
						hx-on:htmx:after-swap="copyEvalResult();scrollBottom();"></pre>
				</div>
				<form
					class="flex h-fit flex-row gap-4 outline-blue-200 dark:outline-white"
					hx-post="/api/evaluate"
					hx-target="#repl-result"
					hx-swap="innerHTML"
					hx-on:htmx:config-request="addCommand()"
					hx-ext="json-enc"
				>
					<p class="hidden sm:my-auto sm:block">Custom command:</p>
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
					/>
					<button
						class="activate:dark:bg-blue-900 my-auto block rounded-lg border-2 bg-blue-200 p-1 active:bg-blue-200 dark:bg-[#090d1a]"
						id="repl-send"
						type="Summit"
						click="addCommand()"
					>
						Run
					</button>
				</form>
			</div>
		</div>
		<div class="flex-1">@exampleArticalFragment()</div>
	</div>
</main>
<Footer />
