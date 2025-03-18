import * as d3 from "d3";

export interface Visualization {
    Cluster: ClusterView
    Services: ServiceView[]
}

interface ClusterView {
    Name: string
    Nodes: NodeView[]
    ControlPlaneURL: string
    Timeout: string
    QPS: number
    Burst: number
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

export interface NodeDatum extends d3.SimulationNodeDatum {
    id: string
    controlPlaneURL?: string | ""
    timeout?: string
    qps?: number
    burst?: number
    icon: string
}

const nodeDatumIsEqual = (node1: NodeDatum, node2: NodeDatum) => {
    return node1.id === node2.id
}

export interface LinkDatum {
    source: string | NodeDatum;
    target: string | NodeDatum;
}

const getId = (value: string | NodeDatum): string => {
    return typeof value === 'string' ? value : value.id;
};

const linkDatumIsEqual = (link1: LinkDatum, link2: LinkDatum) => {
    let sourceId1 = getId(link1.source);
    let sourceId2 = getId(link2.source);
    let targetId1 = getId(link1.source);
    let targetId2 = getId(link2.source);
    return sourceId1 === sourceId2 && targetId1 === targetId2
}

export class NodeLinks {
    nodes: NodeDatum[]
    links: LinkDatum[]

    constructor(nodes: NodeDatum[], links: LinkDatum[]) {
        this.nodes = nodes;
        this.links = links;
    }

    isEqual(nodeLink: NodeLinks): boolean {
        if (this.nodes.length != nodeLink.nodes.length || this.links.length != nodeLink.links.length) return false

        console.log("length true")

        const nodesAreEqual = this.nodes.every((node1, index) => {
            const node2 = nodeLink.nodes[index];
            return nodeDatumIsEqual(node1, node2);
        });

        if (!nodesAreEqual) {
            console.log("nodes false")
            return false;
        }

        return this.links.every((link1, index) => {
            const link2 = nodeLink.links[index];
            return linkDatumIsEqual(link1, link2);
        });
    }

}