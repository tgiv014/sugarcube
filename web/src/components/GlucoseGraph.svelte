<script lang="ts">
	import { glucoseReadings, liveReadings } from '$lib/glucose';
	import { tweened } from 'svelte/motion';
	import { cubicOut } from 'svelte/easing';
	import SvgGraph from './SvgGraph.svelte';

	// Bindings
	let w: number = 0;
	let h: number = 0;
	let follow = true;

	$: xMax = w - 64;
	const baseRangeMillis = 3600 * 3 * 1000; // Default range is 3 hours

	// Data range
	const glucoseMin = tweened(40, { duration: 100, easing: cubicOut });
	const glucoseMax = tweened(250, { duration: 100, easing: cubicOut });
	const endMillis = tweened(new Date().getTime(), { duration: 200, easing: cubicOut });
	const scale = tweened(1, { duration: 100, easing: cubicOut });
	$: rangeMillis = baseRangeMillis * $scale;
	$: startMillis = $endMillis - rangeMillis;

	// Snap to now when initially enabling follow or when liveReadings are updated
	$: if ($liveReadings && follow) {
		endMillis.set(new Date().getTime());
	}

	let timer: NodeJS.Timeout;
	$: if (startMillis && endMillis) {
		clearTimeout(timer);
		timer = setTimeout(() => {
			glucoseReadings.get(new Date(startMillis - rangeMillis), new Date($endMillis + rangeMillis));
		}, 100);
	}

	$: if ($glucoseReadings && $glucoseReadings.length > 0) {
		// Use data min and max
		glucoseMin.set(
			Math.min(
				...$glucoseReadings.map((d) => {
					const t = d.timestamp.getTime();
					if (t >= $endMillis || t <= startMillis) {
						return 40;
					}
					return d.value - 30;
				})
			)
		);
		glucoseMax.set(
			Math.max(
				...$glucoseReadings.map((d) => {
					const t = d.timestamp.getTime();
					if (t >= $endMillis || t <= startMillis) {
						return 200;
					}
					return d.value + 30;
				})
			)
		);
	}

	// Drag behavior!
	let dragging = false;
	let dragX: number; // X coordinate of the start of the drag
	let dragtEnd: number; // Where was tEnd when drag began?
	let dragNow: number; // The current time when drag began
	const gestureHandler = {
		mousedown: (e: MouseEvent) => {
			dragX = e.x;
			dragtEnd = $endMillis;
			dragNow = new Date().getTime();
			dragging = true;
		},
		mousemove: (e: MouseEvent) => {
			if (!dragging) {
				return;
			}
			let delta = e.x - dragX;
			let tEndNumber = dragtEnd - delta * (rangeMillis / xMax);
			if (tEndNumber > dragNow) {
				tEndNumber = dragNow;
				follow = true;
			} else {
				follow = false;
			}
			endMillis.set(new Date(tEndNumber).getTime());
		},
		mouseup: (e: MouseEvent) => {
			dragging = false;
		},
		wheel: (e: WheelEvent) => {
			if (e.deltaY != 0) {
				let delta = e.deltaY / 200;
				let newScale = $scale * Math.pow(2, delta);
				if (newScale > 5) {
					newScale = 5;
				}
				if (newScale < 0.2) {
					newScale = 0.2;
				}

				scale.set(newScale);
			}
			if (e.deltaX != 0) {
				let delta = e.deltaX;
				let tEndNumber = $endMillis + delta * (rangeMillis / xMax);
				if (tEndNumber > new Date().getTime()) {
					tEndNumber = new Date().getTime();
					follow = true;
				} else {
					follow = false;
				}
				endMillis.set(new Date(tEndNumber).getTime());
			}
		}
	};
</script>

<div class="flex h-full w-full flex-col gap-2 font-mono">
	<div
		bind:clientWidth={w}
		bind:clientHeight={h}
		class="svg-container h-80 w-full"
		on:pointerdown|preventDefault={gestureHandler.mousedown}
		on:pointermove={gestureHandler.mousemove}
		on:pointerup={gestureHandler.mouseup}
		on:wheel|preventDefault={gestureHandler.wheel}
	>
		<SvgGraph
			{w}
			{h}
			{startMillis}
			endMillis={$endMillis}
			glucoseMin={$glucoseMin}
			glucoseMax={$glucoseMax}
			data={$glucoseReadings}
		/>
	</div>
	<div class="px-4">
		<label class="relative inline-flex cursor-pointer items-center">
			<input bind:checked={follow} type="checkbox" class="peer sr-only" />
			<div
				class="peer h-6 w-11 rounded-full bg-gray-200 after:absolute after:left-[2px] after:top-0.5 after:h-5 after:w-5 after:rounded-full after:border after:border-gray-300 after:bg-white after:transition-all after:content-[''] peer-checked:bg-blue-600 peer-checked:after:translate-x-full peer-checked:after:border-white peer-focus:ring-4 peer-focus:ring-blue-300 dark:border-gray-600 dark:bg-gray-700 dark:peer-focus:ring-blue-800"
			/>
			<span class="ml-3 text-sm font-medium text-gray-900 dark:text-gray-300"
				>Follow New Readings</span
			>
		</label>
	</div>
</div>

<style>
	.svg-container {
		overflow: auto;
		position: relative;
		touch-action: none;
	}
</style>
