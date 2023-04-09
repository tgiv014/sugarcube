import { writable } from "svelte/store";
import type { ErrorResponse, StatusResponse } from "./types";

export let status = writable<StatusResponse>(undefined)

export const getStatus = async () => {
    const response = await fetch("/api/status")

    if (response.status != 200) {
        const err = await response.json() as ErrorResponse
        throw new Error(err.error)
    }
    const obj = await response.json() as StatusResponse

    status.set(obj)
}

getStatus(); // Get on load

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

export const getReadings = async (start: Date, end: Date) => {

}