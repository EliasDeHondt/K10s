package kubernetes

type Pod struct {
}

type ClusterStructure struct {
	Name  string
	Nodes []NodeTree
}

type PodTree struct {
	Pod Pod
}
