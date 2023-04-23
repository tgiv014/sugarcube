<script lang="ts">
	import { onDestroy } from 'svelte';
	import GlucoseGraph from '../../components/GlucoseGraph.svelte';
	import { type ErrorResponse, type GlucoseReadingResponse, GlucoseReading } from '../../lib/types';
	import { formatReadingTime } from '$lib/utils';

	let readings: GlucoseReading[] = [];

	let latestReading: GlucoseReading | undefined;
	let timestampstring: string = '';
	let delta: number | undefined;

	$: if (readings.length > 0) {
		latestReading = readings[readings.length - 1];
	}

	$: if (readings.length >= 2) {
		const nextNewestReading = readings[readings.length - 2];
		const timeDelta =
			(latestReading!.timestamp.valueOf() - nextNewestReading.timestamp.valueOf()) / 1000;

		if (Math.round(timeDelta / 60) == 5) {
			delta = latestReading!.value - nextNewestReading.value;
		}
	}

	const fetchReadings = async () => {
		const response = await fetch('/api/readings');
		if (response.status != 200) {
			const err = (await response.json()) as ErrorResponse;
			throw new Error(err.error);
		}
		const responseReadings = (await response.json()) as GlucoseReadingResponse[];

		readings = responseReadings.map((reading) => new GlucoseReading(reading));
	};
	fetchReadings();
	const fetchReadingsInterval = setInterval(fetchReadings, 1000 * 60);

	type baseEvent = {
		Topic: string;
		Payload: any;
	};
	const events = new EventSource('/api/bus');
	events.addEventListener('message', async (ev) => {
		console.log(ev);
		const parsedEvent = JSON.parse(ev.data) as baseEvent;
		console.log(parsedEvent);
		if (parsedEvent.Topic == 'newReading') {
			fetchReadings();
		}
	});

	const updateTimestamp = async () => {
		if (!latestReading) {
			return;
		}
		timestampstring = formatReadingTime(
			(new Date().valueOf() - latestReading.timestamp.valueOf()) / 1000
		);
	};
	updateTimestamp();
	const updateTimestampInterval = setInterval(updateTimestamp, 1000);

	onDestroy(() => {
		clearInterval(fetchReadingsInterval);
		clearInterval(updateTimestampInterval);
	});
</script>

<div class="">
	<div class="flex flex-col md:flex-row">
		<div class="my-2 md:flex-grow">
			<GlucoseGraph data={readings} />
		</div>
		<div class="min-w-fit border-l border-stone-900 px-4 py-4 font-mono">
			{#if latestReading}
				<h1 class="text-8xl font-bold">{latestReading.value}</h1>
				{#if delta}
					<p class="text-right text-2xl">{delta > 0 ? '+' : ''}{delta} / 5min</p>
				{/if}
				<p class="text-right text-2xl italic">
					{timestampstring}
				</p>
			{:else}
				<h1 class="text-8xl font-bold">---</h1>
				<p class="italic">No recent readings...</p>
			{/if}
		</div>
	</div>
	<!-- <div class="flex border-t border-stone-900">
		<p>Oh you just know you're getting widgets here for stats</p>
	</div> -->
</div>
