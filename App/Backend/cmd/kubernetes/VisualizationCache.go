/**********************************/
/* @since 01/01/2025              */
/* @author K10s Open Source Team  */
/**********************************/
package kubernetes

import (
	"context"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"log"
	"sync"
	"time"
)

var (
	VisualizationReady  sync.WaitGroup
	CachedVisualization *Visualization
)

func CreateVisualization(client IClient) *Visualization {
	defer VisualizationReady.Done()
	visualization := VisualizeCluster(client)
	go watchNodes(client, visualization)
	go watchDeployments(client, visualization)
	go watchServices(client, visualization)
	CachedVisualization = visualization
	return visualization
}

func watchNodes(client IClient, visualization *Visualization) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	watcher, err := client.GetNodes().Watch(ctx, metav1.ListOptions{})
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Stop()

	for event := range watcher.ResultChan() {
		node, ok := event.Object.(*corev1.Node)
		if !ok {
			log.Printf("Unexpected event object type: %T\n", event.Object)
			continue
		}

		switch event.Type {
		case watch.Added:
			visualization.AddNode(node, client)
			break
		case watch.Modified:
			break
		case watch.Deleted:
			visualization.DeleteNode(node)
		case watch.Error:
			log.Printf("Error event: %v\n", event.Object)
			break
		}
	}
}

func watchServices(client IClient, visualization *Visualization) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	watcher, err := client.GetServices("").Watch(ctx, metav1.ListOptions{})
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Stop()

	for event := range watcher.ResultChan() {
		service, ok := event.Object.(*corev1.Service)
		if !ok {
			log.Printf("Unexpected event object type: %T\n", event.Object)
			continue
		}

		switch event.Type {
		case watch.Added:
			visualization.AddService(service, client)
			break
		case watch.Modified:
			break
		case watch.Deleted:
			visualization.DeleteService(service)
			break
		case watch.Error:
			log.Printf("Error event: %v\n", event.Object)
			break
		}
	}
}

func watchDeployments(client IClient, visualization *Visualization) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	watcher, err := client.GetDeployments("").Watch(ctx, metav1.ListOptions{})
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Stop()

	for event := range watcher.ResultChan() {
		deployment, ok := event.Object.(*appsv1.Deployment)
		if !ok {
			log.Printf("Unexpected event object type: %T\n", event.Object)
			continue
		}

		switch event.Type {
		case watch.Added:
			visualization.AddDeployment(deployment, client)
			break
		case watch.Modified:
			break
		case watch.Deleted:
			visualization.DeleteDeployment(deployment)
			break
		case watch.Error:
			log.Printf("Error event: %v\n", event.Object)
			break
		}
	}
}
