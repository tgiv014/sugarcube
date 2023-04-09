import type { GlucoseReadingResponse } from "./types";

const timeFormat = new Intl.RelativeTimeFormat('en', { style: 'short' });

export const formatReadingTime = (duration: number) => {
    if (duration < 60) {
        return timeFormat.format(-Math.round(duration), "seconds")
    }

    if (duration < 20 * 60) {
        return timeFormat.format(-Math.round(duration / 60), "minutes")
    }

    return ">20 min ago"
}

export const readingDate = (reading: GlucoseReadingResponse) => {
    return new Date(reading.ID * 1000)
}