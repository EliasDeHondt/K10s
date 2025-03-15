/**********************************/
/* @since 01/01/2025              */
/* @author K10s Open Source Team  */
/**********************************/
package handlers

import (
	"context"
	"github.com/eliasdehondt/K10s/App/Backend/cmd/kubernetes"
	"github.com/gorilla/websocket"
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
	CachedVisualization *kubernetes.Visualization
)

func CreateVisualization(client kubernetes.IClient) *kubernetes.Visualization {
	defer VisualizationReady.Done()
	visualization := kubernetes.VisualizeCluster(client)
	go watchNodes(client, visualization)
	go watchDeployments(client, visualization)
	go watchServices(client, visualization)
	CachedVisualization = visualization
	return visualization
}

func watchNodes(client kubernetes.IClient, visualization *kubernetes.Visualization) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	watcher, err := client.GetNodes().Watch(ctx, metav1.ListOptions{})
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		watcher.Stop()
		log.Println("Visualization watcher stopped")
	}()

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
		sendVisualizations()
	}
}

func watchServices(client kubernetes.IClient, visualization *kubernetes.Visualization) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	watcher, err := client.GetServices("").Watch(ctx, metav1.ListOptions{})
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		watcher.Stop()
		log.Println("Visualization watcher stopped")
	}()

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
		sendVisualizations()
	}
}

func watchDeployments(client kubernetes.IClient, visualization *kubernetes.Visualization) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	watcher, err := client.GetDeployments("").Watch(ctx, metav1.ListOptions{})
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		watcher.Stop()
		log.Println("Visualization watcher stopped")
	}()

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
		sendVisualizations()
	}
}

func sendVisualizations() {
	conns := GetVisualizationConns()
	for _, conn := range conns {
		err := conn.WriteJSON(CachedVisualization)
		if err != nil {
			log.Println("Error writing visualization stats:", err)
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Println("WebSocket connection closed by client.")
				kubernetes.CloseConn(conn)
			}
		}
	}
}
