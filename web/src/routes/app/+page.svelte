<script lang="ts">
	import GlucoseGraph from '../../components/GlucoseGraph.svelte';
	import CurrentReading from '../../components/CurrentReading.svelte';
</script>

<div class="mx-auto max-w-6xl">
	<div class="flex flex-col md:flex-row">
		<div class="my-2 md:flex-grow">
			<GlucoseGraph />
		</div>
		<CurrentReading />
	</div>
</div>
