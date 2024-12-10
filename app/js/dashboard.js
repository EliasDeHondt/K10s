/**********************************/
/* @since 01/01/2025              */
/* @author K10s Open Source Team  */
/**********************************/

// Load external content
document.addEventListener('DOMContentLoaded', function() {
    loadExternalContent("navigation-bar", "/app/includes/navigation-bar.html");
    loadExternalContent("context-menu", "/app/includes/context-menu.html");
});

//Temp
// Set up the data
const nodes = [
    // Clusters
    { id: "Cluster01", group: "cluster" },

    // Pools
    { id: "Pool01", group: "pool", cluster: "Cluster01" },
    { id: "Pool02", group: "pool", cluster: "Cluster01" },

    // Nodes
    { id: "Node01", group: "node", pool: "Pool01" },
    { id: "Node02", group: "node", pool: "Pool01" },
    { id: "Node03", group: "node", pool: "Pool02" },
    { id: "Node04", group: "node", pool: "Pool02" },

    // Services
    { id: "Service01", group: "service", node: "Node01" },
    { id: "Service02", group: "service", node: "Node01" },
    { id: "Service03", group: "service", node: "Node02" },
    { id: "Service04", group: "service", node: "Node03" },
    { id: "Service05", group: "service", node: "Node04" },

    // Containers
    { id: "Container01", group: "container", service: "Service01" },
    { id: "Container02", group: "container", service: "Service01" },
    { id: "Container03", group: "container", service: "Service02" },
    { id: "Container04", group: "container", service: "Service03" },
    { id: "Container05", group: "container", service: "Service04" },
    { id: "Container06", group: "container", service: "Service05" }

];

const links = [
    { source: "Cluster01", target: "Pool01" },
    { source: "Cluster01", target: "Pool02" },
    { source: "Pool01", target: "Node01" },
    { source: "Pool01", target: "Node02" },
    { source: "Pool02", target: "Node03" },
    { source: "Pool02", target: "Node04" },
    { source: "Node01", target: "Service01" },
    { source: "Node01", target: "Service02" },
    { source: "Node02", target: "Service03" },
    { source: "Node03", target: "Service04" },
    { source: "Node04", target: "Service05" },
    { source: "Service01", target: "Container01" },
    { source: "Service01", target: "Container02" },
    { source: "Service02", target: "Container03" },
    { source: "Service03", target: "Container04" },
    { source: "Service04", target: "Container05" },
    { source: "Service05", target: "Container06" }
];

// Set up dimensions
const width = 1000;
const height = 1000;

// Create the SVG canvas
const svg = d3.select("body")
    .append("svg")
    .attr("width", width)
    .attr("height", height);

// Create a simulation with forces
const simulation = d3.forceSimulation(nodes)
    .force("link", d3.forceLink(links).id(d => d.id).distance(100))
    .force("charge", d3.forceManyBody().strength(-300))
    .force("center", d3.forceCenter(width / 2, height / 2));

// Add links (lines)
const link = svg.append("g")
    .attr("stroke", "#ac8fbd")
    .selectAll("line")
    .data(links)
    .join("line")
    .attr("class", "link");

// Add nodes (circles)
const node = svg.append("g")
    .selectAll("circle")
    .data(nodes)
    .join("circle")
    .attr("class", "node")
    .attr("r", 10)
    .call(drag(simulation));

// Add labels to the nodes
const labels = svg.append("g")
    .selectAll("text")
    .data(nodes)
    .join("text")
    .attr("text-anchor", "middle")
    .attr("dy", 4)
    .text(d => d.id);

// Update positions on each tick
simulation.on("tick", () => {
link
    .attr("x1", d => d.source.x)
    .attr("y1", d => d.source.y)
    .attr("x2", d => d.target.x)
    .attr("y2", d => d.target.y);

node
    .attr("cx", d => d.x)
    .attr("cy", d => d.y);

labels
    .attr("x", d => d.x)
    .attr("y", d => d.y);
});

// Drag behavior
function drag(simulation) {
return d3.drag()
    .on("start", event => {
    if (!event.active) simulation.alphaTarget(0.3).restart();
    event.subject.fx = event.subject.x;
    event.subject.fy = event.subject.y;
    })
    .on("drag", event => {
    event.subject.fx = event.x;
    event.subject.fy = event.y;
    })
    .on("end", event => {
    if (!event.active) simulation.alphaTarget(0);
    event.subject.fx = null;
    event.subject.fy = null;
    });
}