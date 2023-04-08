<script lang="ts">
	import { draw } from 'svelte/transition';
	import { scaleLinear } from 'd3-scale';
	import { line, curveBasis } from 'd3-shape';
	import { extent } from 'd3';
	import type { GlucoseReading } from '../lib/types';

	export let data: GlucoseReading[];

	let container: HTMLDivElement;
	let svg: SVGElement;
	let w: number = 0;
	let h: number = 0;

	$: extents = extent(data.map((d) => d.ID)) as [number, number];
	$: xScale = scaleLinear().domain(extents).range([0, w]);
	$: yScale = scaleLinear().domain([40, 300]).range([h, 0]);

	$: pathLine = line<GlucoseReading>()
		.x((d) => xScale(d.ID))
		.y((d) => yScale(d.Value))
		.curve(curveBasis);
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
			<!-- transform="translate({xScale( -->
			{#each xScale.ticks(10) as tick, i (tick)}
				<g class="tick" transform="translate({xScale(tick)},0)">
					<line class="gridline" y1={h} y2="0" x1="0" x2="0" />
					<text x={xScale(tick)} y={h}>
						{new Date(tick * 1000).toLocaleTimeString()}
					</text>
				</g>
			{/each}

			<path d={pathLine(data)} />
			{#each data as d}
				<circle cx={xScale(d.ID)} cy={yScale(d.Value)} r="2px" />
			{/each}
		</svg>
	</div>
	<div class="px-4">
		<p>12 hrs</p>
	</div>
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
		stroke-width: 4px;
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
		stroke-width: 2px;
	}
</style>
