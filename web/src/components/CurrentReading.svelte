<script lang="ts">
	import { onDestroy } from 'svelte';
	import type { GlucoseReading } from '$lib/glucose';
	import { formatReadingTime } from '$lib/utils';
	import { liveReadings } from '$lib/glucose';

	let latestReading: GlucoseReading | undefined;
	let timestampstring: string = '';
	let delta: number | undefined = undefined;

	// Computed values
	$: if ($liveReadings) {
		if ($liveReadings.length >= 1) {
			latestReading = $liveReadings[$liveReadings.length - 1];
		}
		if ($liveReadings.length >= 2) {
			const nextNewestReading = $liveReadings[$liveReadings.length - 2];
			const timeDelta =
				(latestReading!.timestamp.valueOf() - nextNewestReading.timestamp.valueOf()) / 1000;

			if (Math.round(timeDelta / 60) == 5) {
				delta = latestReading!.value - nextNewestReading.value;
			}

			updateTimestamp();
		}
	}

	const updateTimestamp = async () => {
		if (!latestReading) {
			return;
		}
		timestampstring = formatReadingTime(
			(new Date().valueOf() - latestReading.timestamp.valueOf()) / 1000
		);
	};
	const updateTimestampInterval = setInterval(updateTimestamp, 1000);

	onDestroy(() => {
		clearInterval(updateTimestampInterval);
	});
</script>

<div class="flex min-w-fit items-center border-l border-stone-900 px-4 py-4 font-mono md:flex-col">
	{#if latestReading}
		<h1 class="text-8xl font-bold">{latestReading.value}</h1>
		<div class="flex flex-grow flex-col">
			{#if delta !== undefined}
				<p class="text-right text-2xl">{delta > 0 ? '+' : ''}{delta} / 5min</p>
			{/if}
			<p class="text-right text-2xl italic">
				{timestampstring}
			</p>
		</div>
	{:else}
		<h1 class="text-8xl font-bold">---</h1>
		<p class="italic">No recent readings...</p>
	{/if}
</div>
