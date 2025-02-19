export interface PaginatedResponse {
    Response: Pod[]
    PageToken: string
}

export interface Pod {
    Namespace: string
    Name: string
    ServicesReady: string
    Restarts: number
    Status: string
    IP: string
    Node: string
    Age: string
}