<script lang="ts">
	import { draw } from 'svelte/transition';
	import { scaleLinear, scaleTime } from 'd3-scale';
	import { line, curveBasis, curveCatmullRom } from 'd3-shape';
	import { extent } from 'd3';
	import type { GlucoseReading } from '../lib/types';

	export let data: GlucoseReading[];

	let container: HTMLDivElement;
	let svg: SVGElement;
	let w: number = 0;
	let h: number = 0;

	$: yMax = h - 18;

	$: extents = extent(data.map((d) => d.timestamp)) as [Date, Date];
	$: xScale = scaleTime().domain(extents).range([0, w]);
	$: yScale = scaleLinear().domain([40, 300]).range([yMax, 0]);

	$: pathLine = line<GlucoseReading>()
		.x((d) => xScale(d.timestamp))
		.y((d) => yScale(d.value))
		.curve(curveCatmullRom);
</script>

<!-- <svg viewBox="0 0 100 100"> <path transition:draw={{ duration: 200 }} d={pathLine(data)} /></svg> -->

<div class="flex h-full w-full flex-col font-mono">
	<div
		bind:this={container}
		bind:clientWidth={w}
		bind:clientHeight={h}
		class="svg-container h-64 w-full"
	>
		<svg bind:this={svg} width={w} height={h}>
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
			<line y1={yScale(60)} y2={yScale(60)} x1="0" x2="100%" class="low-line" />
			<rect
				x="0"
				y={yScale(60)}
				width="100%"
				height={yScale(40) - yScale(60)}
				fill="url(#lowHatch)"
			/>
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
			<line y1={yScale(180)} y2={yScale(180)} x1="0" x2="100%" class="high-line" />
			<rect x="0" y="0" width="100%" height={yScale(180)} fill="url(#highHatch)" />
			<!-- transform="translate({xScale( -->
			{#each xScale.ticks(6) as tick, i (tick)}
				<g class="tick" transform="translate({xScale(tick)},0)">
					<line class="gridline" y1={yMax} y2="0" x1="0" x2="0" />
					<text x={0} y={h} text-anchor="middle">
						{tick.toLocaleTimeString('en', { timeStyle: 'short' })}
						<!-- {i} -->
					</text>
				</g>
			{/each}

			<path d={pathLine(data)} />
			{#each data as d}
				<circle cx={xScale(d.timestamp)} cy={yScale(d.value)} r="4px" />
			{/each}
		</svg>
	</div>
	<div class="px-4" />
</div>

<style>
	.svg-container {
		overflow: auto;
		position: relative;
	}
	svg {
		position: absolute;
		top: 0;
		left: 0;
		overflow: auto;
	}
	path {
		stroke: theme('colors.stone.400');
		stroke-width: 8px;
		fill: none;
		stroke-linecap: round;
	}
	circle {
		stroke: theme('colors.stone.900');
		stroke-width: 2px;
		fill: white;
	}
	.gridline {
		stroke: theme('colors.stone.400');
		stroke-width: 1px;
		stroke-dasharray: 4;
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
</style>
