import { writable } from "svelte/store";
import { status } from "./status";
import type { ErrorResponse } from "./types";

export type Settings = {
    dexcomUsername: string
    dexcomPassword?: string
}

export type SettingsUpdate = Partial<Settings>

function createSettings() {
    const { subscribe, set, update } = writable<Settings>({
        dexcomUsername: ''
    })

    return {
        subscribe,
        get: async () => {
            const response = await fetch("/api/settings")

            if (response.status != 200) {
                const err = await response.json() as ErrorResponse
                throw new Error(err.error)
            }
            const obj = await response.json() as Settings

            set(obj)
            return obj
        },
        update: async (update: SettingsUpdate) => {
            const response = await fetch('/api/settings', {
                method: 'PATCH',
                body: JSON.stringify(update),
                headers: {
                    "Content-Type": "application/json",
                },
            })

            if (response.status != 200) {
                const err = await response.json() as ErrorResponse
                throw new Error(err.error)
            }

            const obj = await response.json() as Settings

            set(obj)

            return obj
        }
    }
}

export const settings = createSettings();

status.subscribe((s) => {
    if (!s) {
        return
    }
    if (s.sessionValid) {
        settings.get()
    }
})