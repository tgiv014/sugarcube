import { writable } from "svelte/store";
import type { ErrorResponse, Status, Settings } from "./types";

export let status = writable<Status>(undefined)

// Status
export const getStatus = async () => {
    const response = await fetch("/api/status")

    if (response.status != 200) {
        const err = await response.json() as ErrorResponse
        throw new Error(err.error)
    }
    const obj = await response.json() as Status

    status.set(obj)
    return obj
}

export const login = async (loginRequest: any) => {
    const response = await fetch('/api/login', {
        method: 'POST',
        body: JSON.stringify(loginRequest),
        headers: {
            "Content-Type": "application/json",
        },
    })

    if (response.status != 200) {
        const err = await response.json() as ErrorResponse
        throw new Error(err.error)
    }

    await getStatus();
}

// Auth
export const logout = async () => {
    const response = await fetch('/api/logout', {
        method: 'POST',
        headers: {
            "Content-Type": "application/json",
        },
    })

    if (response.status != 200) {
        const err = await response.json() as ErrorResponse
        throw new Error(err.error)
    }

    await getStatus();
}

export const signup = async (password: string) => {
    const response = await fetch('/api/signup', {
        method: 'POST',
        body: JSON.stringify({
            password: password,
        }),
        headers: {
            "Content-Type": "application/json",
        },
    })

    if (response.status != 200) {
        const err = await response.json() as ErrorResponse
        throw new Error(err.error)
    }

    await getStatus();
}

// Settings
export let settings = writable<Settings>(undefined)


export const getSettings = async () => {
    const response = await fetch("/api/settings")

    if (response.status != 200) {
        const err = await response.json() as ErrorResponse
        throw new Error(err.error)
    }
    const obj = await response.json() as Settings

    settings.set(obj)

    return obj
}

export type SettingsUpdate = Partial<Settings>

export const updateSettings = async (update: SettingsUpdate) => {
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

    settings.set(obj)

    await getStatus();

    return obj
}

getStatus().then(async (status) => {
    if (status.sessionValid) {
        await getSettings();
    }
});