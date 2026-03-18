<script lang="ts">
	let { href, name }: { href: undefined | null | string; name: string } = $props();
</script>

<li
	class="block border-l-2 border-gray-200 py-2 pl-4 hover:border-blue-500 hover:font-bold focus:border-blue-500 focus:font-bold"
>
	<a {href}>{name}</a>
</li>
