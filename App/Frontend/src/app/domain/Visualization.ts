export interface Visualization {
    Cluster: ClusterView
    Services: ServiceView[]
}

interface ClusterView {
    Name: string
    Nodes: NodeView[]
}

interface NodeView {
    Name: string
    Namespace: string
    Deployments: DeploymentView[]
}

interface ServiceView {
    Name: string
    Deployments: DeploymentView[]
    LoadBalancers: LoadBalancer[]
}

interface LoadBalancer {
    HostName: string
    IP: string
}

interface DeploymentView {
    Name: string
}

