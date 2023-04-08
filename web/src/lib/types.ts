export type StatusResponse = {
    setupCompleted: boolean
    sessionValid: boolean
}

export type ErrorResponse = {
    error: string
}

export type GlucoseReading = {
    ID: number
    CreatedAt: string
    UpdatedAt: string
    DeletedAt: string
    Value: number
    Timestamp: string
}