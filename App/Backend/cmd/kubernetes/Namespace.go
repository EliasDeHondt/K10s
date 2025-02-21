package kubernetes

import corev1 "k8s.io/api/core/v1"

type Namespace struct {
	Name string
}

func NewNamespace(namespace corev1.Namespace) Namespace {
	return Namespace{
		Name: namespace.Name,
	}
}
