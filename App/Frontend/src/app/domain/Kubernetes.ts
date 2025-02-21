/**********************************/
/* @since 01/01/2025              */
/* @author K10s Open Source Team  */
/**********************************/

export interface PaginatedResponse {
    Response: Pod[] | Node[] | Service[] | Deployment[] | ConfigMap[] | Secret[]
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

export interface Node {
    Name: string
    Status: string
    Role: string
    Version: string
    PodsAmount: number
    NodeAge: string
    IP: string
}

export interface Service {
    Namespace: string
    Name: string
    Type: string
    ClusterIp: string
    ExternalIp: string[]
    Ports: number[]
    Age: string
}

export interface Deployment {
    Namespace: string
    Name: string
    Ready: string
    Updated: boolean
    Available: boolean
    Age: string
}

export interface ConfigMap {
    Namespace: string
    Name: string
    Data: { [key: string]: string}
    Age: string
}

export interface Secret {
    Namespace: string
    Name: string
    Type: string
    Data: { [key: string]: string }
    Age: string
}