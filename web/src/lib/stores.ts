import { writable } from "svelte/store";
import type { StatusResponse } from "./types";

export let status = writable<StatusResponse>(undefined)

export const getStatus = async () => {
    fetch("/api/status")
        .then(res => res.json())
        .then(data => {
            status.set(data)
        })
}

getStatus();