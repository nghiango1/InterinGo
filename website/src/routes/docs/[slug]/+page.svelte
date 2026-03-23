<script lang="ts">
	import type { PageProps } from './$types';
	let { data }: PageProps = $props();
</script>

{#key data.content}
    <data.content />
{/key}
