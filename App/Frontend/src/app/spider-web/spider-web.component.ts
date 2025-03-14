/**********************************/
/* @since 01/01/2025              */
/* @author K10s Open Source Team  */
/**********************************/
import {Component, ElementRef, AfterViewInit, ViewChild, inject} from '@angular/core';
import * as d3 from 'd3';
import {VisualizationService} from "../services/visualization.service";
import {NotificationService} from "../services/notification.service";
import {TranslateService} from "@ngx-translate/core";
import {Visualization} from "../domain/Visualization";
import {vhToPixels, vwToPixels} from "../domain/util";

interface NodeDatum extends d3.SimulationNodeDatum {
    id: string;
    icon: string;
}

interface LinkDatum {
    source: string | NodeDatum;
    target: string | NodeDatum;
}

@Component({
    selector: 'app-spider-web',
    templateUrl: './spider-web.component.html',
    standalone: true,
    styleUrls: ['./spider-web.component.css'],
})
export class SpiderWebComponent implements AfterViewInit {
    visualizationService = inject(VisualizationService);
    @ViewChild('svgContainer', { static: true }) svgRef!: ElementRef<SVGSVGElement>;
    constructor(private notificationService: NotificationService,private translate: TranslateService) {}

    private graphData: { nodes: NodeDatum[]; links: LinkDatum[] } = { nodes: [], links: [] };

    ngAfterViewInit(): void {
        this.visualizationService.getVisualization().subscribe({
            next: (data: Visualization) => {
                this.updateGraphData(data);
                this.createForceDirectedGraph();
            },
            error: () => {
                this.notificationService.showNotification(this.translate.instant('NOTIF.VISUALIZATION.GETERROR'), 'error');
            },
        });

    }

    private updateGraphData(data: Visualization): void {
        const nodes: NodeDatum[] = [];
        const links: LinkDatum[] = [];
        const nodeMap = new Map<string, NodeDatum>();

        const addNode = (id: string, icon: string) => {
            if (!nodeMap.has(id)) {
                const node = { id, icon };
                nodeMap.set(id, node);
                nodes.push(node);
            }
            return nodeMap.get(id)!;
        };

        addNode(data.Cluster.Name, 'dashboard-cluster.svg');

        data.Cluster.Nodes.forEach((node) => {
            addNode(node.Name, 'dashboard-server.svg');
            links.push({ source: data.Cluster.Name, target: node.Name });

            node.Deployments.forEach((deployment) => {
                addNode(deployment.Name, 'dashboard-deployment.svg');
                links.push({ source: node.Name, target: deployment.Name });
            });
        });

        data.Services.forEach((service) => {
            addNode(service.Name, 'dashboard-service.svg');

            service.Deployments.forEach((deployment) => {
                addNode(deployment.Name, 'dashboard-deployment.svg');
                links.push({ source: deployment.Name, target: service.Name });
            });
            service.LoadBalancers.forEach((lb, index) => {
                const lbId = `${service.Name}-lb-${index + 1}`;
                addNode(lbId, 'dashboard-ip.svg');
                links.push({ source: service.Name, target: lbId });
            });
        });

        this.graphData = { nodes, links };
//     this.graphData = {
//         nodes: [
//             { id: 'Supercluster01', icon: 'dashboard-supercluster.svg' },
//             { id: 'Cluster01', icon: 'dashboard-cluster.svg' },
//             { id: 'Cluster02', icon: 'dashboard-cluster.svg' },
//             { id: 'Node001', icon: 'dashboard-server.svg' },
//             { id: 'Node002', icon: 'dashboard-server.svg' },
//             { id: 'Node003', icon: 'dashboard-server.svg' },
//             { id: 'Node004', icon: 'dashboard-server.svg' },
//             { id: 'Node005', icon: 'dashboard-server.svg' },
//             { id: 'Deployment', icon: 'dashboard-deployment.svg' },
//             { id: 'Service', icon: 'dashboard-service.svg' },
//             { id: 'IP', icon: 'dashboard-ip.svg' },
//         ],
//         links: [
//             { source: 'Supercluster01', target: 'Cluster01' },
//             { source: 'Supercluster01', target: 'Cluster02' },
//             { source: 'Cluster01', target: 'Node001' },
//             { source: 'Cluster01', target: 'Node002' },
//             { source: 'Cluster01', target: 'Node003' },
//             { source: 'Cluster02', target: 'Node004' },
//             { source: 'Cluster02', target: 'Node005' },
//             { source: 'Node001', target: 'Deployment' },
//             { source: 'Node002', target: 'Deployment' },
//             { source: 'Node003', target: 'Deployment' },
//             { source: 'Deployment', target: 'Service' },
//             { source: 'Service', target: 'IP' },
//         ],
//     };
    }


    private createForceDirectedGraph(): void {
        const width = vwToPixels(87.5);
        const height = vhToPixels(70);

        const svg = d3
            .select(this.svgRef.nativeElement)
            .attr('width', width)
            .attr('height', height);

        this.graphData.nodes.forEach((node: NodeDatum) => {
            if (node.icon === 'dashboard-cluster.svg') {
                node.x = width / 2;
                node.y = 80;
                node.fx = width / 2;
                node.fy = 80;
            } else if (node.icon === 'dashboard-server.svg') {
                node.y = height * 0.25;
            }else if (node.icon === 'dashboard-deployment.svg') {
                node.y = height * 0.5;
            } else if (node.icon === 'dashboard-service.svg') {
                node.y = height * 0.75;
            } else if (node.icon === 'dashboard-ip.svg') {
                node.y = height * 0.6;
            }
        });

        const simulation = d3
            .forceSimulation<NodeDatum>(this.graphData.nodes)
            .force('link', d3.forceLink<NodeDatum, LinkDatum>(this.graphData.links).id((d) => d.id).distance(150))
            .force('charge', d3.forceManyBody().strength(-800))
            .force('x', d3.forceX(width / 2).strength(0.1))
            .force('y', d3.forceY<NodeDatum>((d) => {
                    if (d.icon === 'dashboard-cluster.svg') return 80;
                    if (d.icon === 'dashboard-server.svg') return height * 0.35;
                    if (d.icon === 'dashboard-deployment.svg') return height * 0.5;
                    if (d.icon === 'dashboard-service.svg') return height * 0.75;
                    if (d.icon === 'dashboard-ip.svg') return height * 0.8;
                    return height / 2;
                }).strength(0.5)
            );

        const link = svg
            .selectAll('.link')
            .data(this.graphData.links)
            .enter()
            .append('line')
            .attr('stroke', '#aaa')
            .attr('stroke-width', 2);

        const node = svg
            .selectAll('.node')
            .data(this.graphData.nodes)
            .enter()
            .append('g')
            .attr('class', 'node') as d3.Selection<SVGGElement, NodeDatum, SVGSVGElement, unknown>;

        const images = node
            .append('image')
            .attr('href', (d) => `/assets/svg/${d.icon}`)
            .attr('width', 80)
            .attr('height', 80)
            .attr('x', -40)
            .attr('y', -40)
            .on('error', () => {
                this.notificationService.showNotification(this.translate.instant('NOTIF.VISUALIZATION.GETERROR'), 'error');
            });

        const tooltip = svg.append('g')
            .attr('class', 'tooltip')
            .style('display', 'none');

        tooltip.append('rect')
            .attr('width', 100)
            .attr('height', 20)
            .attr('fill', '#333')
            .attr('opacity', 0.8);

        tooltip.append('text')
            .attr('x', 50)
            .attr('y', 15)
            .attr('text-anchor', 'middle')
            .attr('fill', '#fff')
            .attr('font-size', '12px');

        images
            .on('mouseover', function (event, d) {
                tooltip.style('display', null);
                tooltip.select('text').text(d.id);
                const [x, y] = d3.pointer(event, svg.node());
                tooltip.attr('transform', `translate(${x + 10},${y - 10})`);
            })
            .on('mouseout', function () {
                tooltip.style('display', 'none');
            });

        node
            .append('text')
            .text((d) => d.id)
            .attr('text-anchor', 'middle')
            .attr('dy', 60)
            .attr('font-size', '12px')
            .attr('fill', 'var(--text)');

        const dragHandler = d3
            .drag<SVGGElement, NodeDatum>()
            .on('start', (event, d) => {
                if (!event.active) simulation.alphaTarget(0.3).restart();
                d.fx = event.x;
                d.fy = event.y;
            })
            .on('drag', (event, d) => {
                d.fx = event.x;
                d.fy = event.y;
            });

        node.call(dragHandler);

        simulation.on('tick', () => {
            node.each((d) => {
                d.x = Math.max(40, Math.min(width - 40, d.x!));
                d.y = Math.max(40, Math.min(height - 40, d.y!));
            });

            link
                .attr('x1', (d) => (d.source as unknown as NodeDatum).x!)
                .attr('y1', (d) => (d.source as unknown as NodeDatum).y!)
                .attr('x2', (d) => (d.target as unknown as NodeDatum).x!)
                .attr('y2', (d) => (d.target as unknown as NodeDatum).y!);

            node.attr('transform', (d) => `translate(${d.x!},${d.y!})`);
        });
    }
}