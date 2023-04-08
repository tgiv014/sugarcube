<script lang="ts">
	import { onDestroy } from 'svelte';
	import GlucoseGraph from '../../components/GlucoseGraph.svelte';
	import type { ErrorResponse, GlucoseReading } from '../../lib/types';

	type ReadingAndCalculatedData = {
		Reading: GlucoseReading;
		Delta: number;
	};
	let readings: GlucoseReading[] = [];

	let latestReading: ReadingAndCalculatedData | undefined;

	$: if (readings.length > 0) {
		latestReading = {
			Reading: readings[readings.length - 1],
			Delta: 0
		};
	}

	const fetchReadings = async () => {
		const response = await fetch('/api/readings');
		if (response.status != 200) {
			const err = (await response.json()) as ErrorResponse;
			throw new Error(err.error);
		}
		readings = (await response.json()) as GlucoseReading[];
	};
	fetchReadings();
	const interval = setInterval(fetchReadings, 1000 * 60);

	onDestroy(() => {
		clearInterval(interval);
	});
</script>

<div class="my-4 flex">
	<GlucoseGraph data={readings} />
	<div class="mx-4 flex-grow font-mono">
		<h1 class="text-8xl font-bold">{latestReading ? latestReading.Reading.Value : '---'}</h1>
		<p class="text-right text-2xl">+10 / 5min</p>
		<p class="text-right text-2xl italic">5 min ago</p>
	</div>
</div>
