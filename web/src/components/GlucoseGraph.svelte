<script lang="ts">
	import { scaleLinear, scaleTime } from 'd3-scale';
	import { line, curveBasis, curveCatmullRom } from 'd3-shape';
	import type { GlucoseReading } from '../lib/types';
	import { tweened } from 'svelte/motion';
	import { cubicOut } from 'svelte/easing';

	// Inputs
	export let data: GlucoseReading[];

	// Bindings
	let w: number = 0;
	let h: number = 0;
	let follow = true;

	$: yMax = h - 32;
	$: xMax = w - 64;

	const baseRange = 3600 * 3;

	// Data range
	let tEnd = new Date();
	const scale = tweened(1, { duration: 400, easing: cubicOut });
	$: rangeSeconds = baseRange * $scale;
	$: tStart = new Date(tEnd.getTime() - rangeSeconds * 1000);

	// D3 Scale Magic
	$: extents = [tStart, tEnd];
	$: xScale = scaleTime().domain(extents).range([0, xMax]);
	$: yScale = scaleLinear().domain([40, 200]).range([yMax, 0]);
	$: pathLine = line<GlucoseReading>()
		.x((d) => xScale(d.timestamp))
		.y((d) => yScale(d.value))
		.curve(curveCatmullRom);

	// When toggling on "follow", snap to now
	$: if (data && follow) {
		tEnd = new Date();
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
		<svg width={w} height={h}>
			<g>
				<g>
					<pattern
						id="lowHatch"
						patternUnits="userSpaceOnUse"
						width="4"
						height="4"
						patternTransform="scale(4 4)"
					>
						<path
							d="M-1,1 l2,-2
						 M0,4 l4,-4
						 M3,5 l2,-2"
							class="low-hash"
						/>
					</pattern>
					<line y1={yScale(60)} y2={yScale(60)} x1="0" x2={xMax} class="low-line" />
					<rect
						x="0"
						y={yScale(60)}
						width={xMax}
						height={yScale(40) - yScale(60)}
						fill="url(#lowHatch)"
					/>
					<text x={xMax + 16} y={yScale(60)}> 60 </text>
				</g>
				<g>
					<pattern
						id="highHatch"
						patternUnits="userSpaceOnUse"
						width="4"
						height="4"
						patternTransform="scale(4 4)"
					>
						<path
							d="M-1,1 l2,-2
					 M0,4 l4,-4
					 M3,5 l2,-2"
							class="high-hash"
						/>
					</pattern>
					<line y1={yScale(180)} y2={yScale(180)} x1="0" x2={xMax} class="high-line" />
					<rect x="0" y="0" width={xMax} height={yScale(180)} fill="url(#highHatch)" />
					<text x={xMax + 16} y={yScale(180)}> 180 </text>
				</g>
				<g>
					<line y1={yScale(100)} y2={yScale(100)} x1="0" x2={xMax} class="goal-line" />
					<text x={xMax + 16} y={yScale(100)}> 100 </text>
				</g>
				<text x={xMax + 16} y={yScale(250)}> 250 </text>

				<!-- transform="translate({xScale( -->
				{#each xScale.ticks(6) as tick, i (tick)}
					<g class="tick" transform="translate({xScale(tick)},0)">
						<line class="grid-line" y1={yMax} y2="0" x1="0" x2="0" />
						<text x={0} y={h - 8} text-anchor="middle">
							{tick.toLocaleTimeString('en', { timeStyle: 'short' })}
							<!-- {i} -->
						</text>
					</g>
				{/each}

				<path d={pathLine(data)} />
				{#each data as d}
					<circle cx={xScale(d.timestamp)} cy={yScale(d.value)} r="4px" />
				{/each}
			</g>
		</svg>
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
		touch-action: pan-x pan-y;
	}
	svg {
		position: absolute;
		top: 0;
		left: 0;
		overflow: auto;
	}
	path {
		stroke: theme('colors.stone.400');
		stroke-width: 4px;
		fill: none;
		stroke-linecap: round;
	}
	circle {
		stroke: theme('colors.stone.900');
		stroke-width: 2px;
		fill: white;
	}
	text {
		fill: theme('colors.stone.900');
	}

	.grid-line {
		stroke: theme('colors.stone.400');
		stroke-width: 1px;
		stroke-dasharray: 4;
	}
	.goal-line {
		stroke: theme('colors.stone.400');
		stroke-width: 1px;
	}
	.low-line {
		fill: none;
		stroke: theme('colors.red.200');
		stroke-width: 2px;
	}
	.low-hash {
		fill: none;
		stroke: theme('colors.red.200');
		stroke-width: 0.5px;
	}
	.high-line {
		fill: none;
		stroke: theme('colors.amber.200');
		stroke-width: 2px;
	}
	.high-hash {
		fill: none;
		stroke: theme('colors.amber.200');
		stroke-width: 0.5px;
	}

	@media (prefers-color-scheme: dark) {
		text {
			fill: theme('colors.stone.100');
		}

		.grid-line {
			stroke: theme('colors.stone.600');
		}
		.goal-line {
			stroke: theme('colors.stone.500');
		}
		.low-line {
			fill: none;
			stroke: theme('colors.red.900');
		}
		.low-hash {
			fill: none;
			stroke: theme('colors.red.900');
		}
		.high-line {
			fill: none;
			stroke: theme('colors.amber.900');
		}
		.high-hash {
			fill: none;
			stroke: theme('colors.amber.900');
		}
	}
</style>
