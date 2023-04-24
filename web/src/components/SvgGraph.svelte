<script lang="ts">
	import { scaleLinear, scaleTime } from 'd3-scale';
	import { line, curveBasis, curveCatmullRom } from 'd3-shape';
	import type { GlucoseReading } from '$lib/glucose';

	// Bindings
	export let w: number = 0;
	export let h: number = 0;
	export let tEnd: Date;
	export let tStart: Date;
	export let data: GlucoseReading[] | undefined;

	$: yMax = h - 32;
	$: xMax = w - 64;

	// D3 Scale Magic
	$: extents = [tStart, tEnd];
	$: xScale = scaleTime().domain(extents).range([0, xMax]);
	$: yScale = scaleLinear().domain([40, 200]).range([yMax, 0]);
	$: pathLine = line<GlucoseReading>()
		.x((d) => xScale(d.timestamp))
		.y((d) => yScale(d.value))
		.curve(curveCatmullRom);
</script>

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

		{#if data}
			<path d={pathLine(data)} />
			{#each data as d}
				<circle cx={xScale(d.timestamp)} cy={yScale(d.value)} r="4px" />
			{/each}
		{/if}
	</g>
</svg>

<style>
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
