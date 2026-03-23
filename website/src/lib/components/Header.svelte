<script lang="ts">
	import { resolve } from '$app/paths';
</script>

<header
	class="my-4 prose flex max-w-none min-w-80 flex-row gap-4 overflow-hidden border-b-2 border-blue-200 px-4 pb-4 text-xl text-blue-500 sm:mx-4 sm:px-0 dark:prose-invert"
>
	<a
		class="m:text-blue-500 hidden min-w-fit grow font-bold sm:block sm:underline"
		href={resolve('/')}
	>
		Welcome to InterinGo language
	</a>
	<a
		class="min-w-fit grow overflow-hidden text-blue-500 sm:hidden sm:underline"
		href={resolve('/')}
	>
		Home
	</a>
	<a class="shrink text-blue-500 sm:underline" href={resolve('/docs')}> Docs </a>
	<a class="shrink text-blue-500 sm:underline" href="https://github.com/nghiango1/InterinGo">
		Source
	</a>
</header>
