import { writable } from "svelte/store";
import type { ErrorResponse } from "./types";

export type Status = {
    setupCompleted: boolean
    sessionValid: boolean
}

function createStatus() {
    const { subscribe, set, update } = writable<Status | undefined>(undefined)

    const status = {
        subscribe,
        refresh: async () => {
            const response = await fetch("/api/status")

            if (response.status != 200) {
                const err = await response.json() as ErrorResponse
                throw new Error(err.error)
            }
            const obj = await response.json() as Status

            set(obj)
            return obj
        }
    }

    status.refresh().then((status) => {
        console.log(status)
    })
    return status
}

export const status = createStatus();