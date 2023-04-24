import { writable } from 'svelte/store';
import type { ErrorResponse } from './types';

export type GlucoseReadingResponse = {
	ID: number;
	CreatedAt: string;
	UpdatedAt: string;
	DeletedAt: string;
	Value: number;
};

export class GlucoseReading {
	timestamp: Date;
	value: number;

	constructor(res: GlucoseReadingResponse) {
		this.timestamp = new Date(res.ID * 1000);
		this.value = res.Value;
	}
}

export const fetchReadings = async (start?: Date, end?: Date) => {
	const params = new URLSearchParams();
	if (start) {
		console.log(start.toISOString());
		params.append('start', start.toISOString());
	}
	if (end) {
		console.log(end.toISOString());
		params.append('end', end.toISOString());
	}
	const response = await fetch('/api/readings?' + params);
	if (response.status != 200) {
		const err = (await response.json()) as ErrorResponse;
		throw new Error(err.error);
	}
	const responseReadings = (await response.json()) as GlucoseReadingResponse[];

	const readings = responseReadings.map((reading) => new GlucoseReading(reading));

	return readings;
};

function createReadings() {
	const { subscribe, set, update } = writable<GlucoseReading[]>([]);

	const readings = {
		subscribe,
		get: async (start: Date, end: Date) => {
			const readings = await fetchReadings(start, end);
			set(readings);
		}
	};

	return readings;
}

export const glucoseReadings = createReadings();

type baseEvent = {
	Topic: string;
	Payload: any;
};

function createLiveReadings() {
	const { subscribe, set, update } = writable<GlucoseReading[]>([]);

	const readings = {
		subscribe,
		get: async () => {
			const readings = await fetchReadings();
			set(readings);
		}
	};
	const events = new EventSource('/api/bus');
	events.addEventListener('message', async (ev) => {
		console.log(ev);
		const parsedEvent = JSON.parse(ev.data) as baseEvent;
		console.log(parsedEvent);
		if (parsedEvent.Topic == 'newReading') {
			readings.get();
		}
	});

	readings.get();
	return readings;
}

export const liveReadings = createLiveReadings();
