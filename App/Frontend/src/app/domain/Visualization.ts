export interface Visualization {
    cluster: ClusterView
    services: ServiceView[]
}

interface ClusterView {
    name: string
    nodes: NodeView[]
}

interface NodeView {
    name: string
    namespace: string
    deployments: DeploymentView[]
}

interface ServiceView {
    name: string
    deployments: DeploymentView[]
    loadBalancers: LoadBalancer[]
}

interface LoadBalancer {
    hostName: string
    IP: string
}

interface DeploymentView {
    name: string
}

