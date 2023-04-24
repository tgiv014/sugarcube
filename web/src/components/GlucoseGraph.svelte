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

	const baseRange = 3600 * 3;

	// Data range
	let tEnd = new Date();
	const scale = tweened(1, { duration: 400, easing: cubicOut });
	$: rangeSeconds = baseRange * $scale;

	$: tStart = new Date(tEnd.getTime() - rangeSeconds * 1000);

	// When toggling on "follow", snap to now
	$: if ($liveReadings && follow) {
		tEnd = new Date();
	}
	let timer: NodeJS.Timeout;
	$: if (tStart && tEnd) {
		clearTimeout(timer);
		timer = setTimeout(() => {
			console.log(tStart, tEnd);
			glucoseReadings.get(
				new Date(tStart.getTime() - rangeSeconds * 1000),
				new Date(tEnd.getTime() + rangeSeconds * 1000)
			);
		}, 500);
	}

	// Drag behavior!
	let dragging = false;
	let dragX: number; // X coordinate of the start of the drag
	let dragtEnd: number; // Where was tEnd when drag began?
	let dragNow: number; // The current time when drag began
	const gestureHandler = {
		mousedown: (e: MouseEvent) => {
			dragX = e.x;
			dragtEnd = tEnd.getTime();
			dragNow = new Date().getTime();
			dragging = true;
		},
		mousemove: (e: MouseEvent) => {
			if (!dragging) {
				return;
			}
			let delta = e.x - dragX;
			let tEndNumber = dragtEnd - 1000 * delta * (rangeSeconds / xMax);
			if (tEndNumber > dragNow) {
				tEndNumber = dragNow;
				follow = true;
			} else {
				follow = false;
			}
			tEnd = new Date(tEndNumber);
		},
		mouseup: (e: MouseEvent) => {
			dragging = false;
		},
		wheel: (e: WheelEvent) => {
			if (e.deltaY == 0) {
				return;
			}
			let delta = e.deltaY / 200;
			let newScale = $scale * Math.pow(2, delta);
			// scale = newScale;
			scale.set(newScale);
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
		on:wheel={gestureHandler.wheel}
	>
		<SvgGraph {w} {h} {tStart} {tEnd} data={$glucoseReadings} />
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
