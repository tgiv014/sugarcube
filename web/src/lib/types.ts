export type Status = {
    setupCompleted: boolean
    sessionValid: boolean
}

export type ErrorResponse = {
    error: string
}

export type GlucoseReadingResponse = {
    ID: number
    CreatedAt: string
    UpdatedAt: string
    DeletedAt: string
    Value: number
}

export class GlucoseReading {
    timestamp: Date
    value: number

    constructor(res: GlucoseReadingResponse) {
        this.timestamp = new Date(res.ID * 1000)
        this.value = res.Value
    }
}

export type Settings = {
    dexcomUsername: string
    dexcomPassword?: string
}