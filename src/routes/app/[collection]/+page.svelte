<script lang="ts">
	import type { PageData } from './$types';

	let { data }: { data: PageData } = $props();
</script>

<div>
	<h1>{data.collection}</h1>
</div>
