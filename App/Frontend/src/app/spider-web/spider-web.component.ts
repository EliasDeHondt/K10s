/**********************************/
/* @since 01/01/2025              */
/* @author K10s Open Source Team  */
/**********************************/
import { Component, OnInit, ViewChild, ElementRef } from '@angular/core';

interface Node {
    id: number;
    x: number;
    y: number;
    radius: number;
}

@Component({
    selector: 'app-spider-web',
    templateUrl: './spider-web.component.html',
    standalone: true,
    styleUrls: ['./spider-web.component.css']
})

export class SpiderWebComponent implements OnInit {
    @ViewChild('canvas', { static: true }) canvasRef!: ElementRef;
    private ctx!: CanvasRenderingContext2D;
    private nodes: Node[] = [];
    private connections: { from: number, to: number }[] = [];
    private isDragging = false;
    private draggedNode: Node | null = null;
    private offsetX = 0;
    private offsetY = 0;

    ngOnInit(): void {
        this.initializeCanvas();
        this.createNodes();
        this.drawWeb();
        this.addEventListeners();
    }

    initializeCanvas() {
        const canvas = this.canvasRef.nativeElement;
        this.ctx = canvas.getContext('2d')!;
        canvas.width = window.innerWidth;
        canvas.height = window.innerHeight;
    }

    createNodes() { //todo echt ophalen
        for (let i = 0; i < 10; i++) {
        this.nodes.push({
            id: i,
            x: Math.random() * window.innerWidth,
            y: Math.random() * window.innerHeight,
            radius: 40
        });
        }

        for (let i = 0; i < this.nodes.length; i++) {
            for (let j = i + 1; j < this.nodes.length; j++) {
                if (Math.random() > 0.5) {
                this.connections.push({ from: i, to: j });
                }
            }
        }
    }

    drawWeb() {
        this.ctx.clearRect(0, 0, window.innerWidth, window.innerHeight);

        const rootStyles = getComputedStyle(document.documentElement);
        const primary = rootStyles.getPropertyValue('--primary').trim();
        const statusBlue = rootStyles.getPropertyValue('--status-blue').trim();
        const statusGrey = rootStyles.getPropertyValue('--status-grey').trim();

        this.connections.forEach(connection => {
        const fromNode = this.nodes[connection.from];
        const toNode = this.nodes[connection.to];
        this.ctx.beginPath();
        this.ctx.moveTo(fromNode.x, fromNode.y);
        this.ctx.lineTo(toNode.x, toNode.y);
        this.ctx.strokeStyle = statusGrey;
        this.ctx.lineWidth = 1;
        this.ctx.stroke();
        });

        this.nodes.forEach(node => {
        this.ctx.beginPath();
        this.ctx.arc(node.x, node.y, node.radius, 0, Math.PI * 2);
        this.ctx.fillStyle = primary;
        this.ctx.fill();
        this.ctx.strokeStyle = statusBlue;
        this.ctx.stroke();
        });
    }

    addEventListeners() {
        const canvas = this.canvasRef.nativeElement;

        canvas.addEventListener('mousedown', (event: MouseEvent) => {
        const mouseX = event.clientX;
        const mouseY = event.clientY;

        this.draggedNode = this.nodes.find(node =>
            Math.sqrt(Math.pow(node.x - mouseX, 2) + Math.pow(node.y - mouseY, 2)) < node.radius
        ) || null;

        if (this.draggedNode) {
            this.isDragging = true;
            this.offsetX = mouseX - this.draggedNode.x;
            this.offsetY = mouseY - this.draggedNode.y;
        }
        });

        canvas.addEventListener('mousemove', (event: MouseEvent) => {
        if (this.isDragging && this.draggedNode) {
            this.draggedNode.x = event.clientX - this.offsetX;
            this.draggedNode.y = event.clientY - this.offsetY;
            this.drawWeb();
        }
        });

        canvas.addEventListener('mouseup', () => {
        this.isDragging = false;
        this.draggedNode = null;
        });
    }
}