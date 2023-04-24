import { status } from "./status";
import type { ErrorResponse } from "./types";

// Auth

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

    status.refresh();
}

export const logout = async () => {
    const response = await fetch('/api/logout', {
        method: 'POST',
    })

    if (response.status != 200) {
        const err = await response.json() as ErrorResponse
        throw new Error(err.error)
    }

    status.refresh();
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

    status.refresh();
}

