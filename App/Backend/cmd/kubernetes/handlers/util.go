/**********************************/
/* @since 01/01/2025              */
/* @author K10s Open Source Team  */
/**********************************/
package handlers

import (
	"context"
	"github.com/eliasdehondt/K10s/App/Backend/cmd/kubernetes"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
	"net/http"
	"strconv"
	"time"
)

var C = kubernetes.GetClients()

type PaginatedResponse[T any] struct {
	Response  T
	PageToken string
}

func GetPageSizeAndPageToken(ctx *gin.Context) (int, string) {
	pageSizeString, _ := ctx.GetQuery("pageSize")
	pageToken, _ := ctx.GetQuery("pageToken")
	pageSize, err := strconv.Atoi(pageSizeString)
	if err != nil {
		pageSize = 20
	}
	return pageSize, pageToken
}

func GetFrontendIP() string {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	svc, err := C.GetServices("k10s-namespaces").Get(ctx, "k10s-ingress-service", metav1.GetOptions{})

	if err != nil {
		log.Printf("Failed to get k10s-ingress-service: %v, using default development url", err)
		return "http://localhost:4200"
	}

	if len(svc.Status.LoadBalancer.Ingress) > 0 {

		protocol := "http"
		if svc.Status.LoadBalancer.Ingress[0].Hostname != "" {
			protocol = "https"
		}

		// For local docker kubernetes testing
		ingressIP := svc.Status.LoadBalancer.Ingress[0].IP
		if ingressIP != "" {
			ingressIP = "localhost"
		}
		url := protocol + "://" + ingressIP

		if svc.Spec.Ports[0].Port != 80 && svc.Spec.Ports[0].Port != 443 {
			url += ":" + strconv.Itoa(int(svc.Spec.Ports[0].Port))
		}
		log.Println(url)
		return url
	}

	return "http://localhost:4200"
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		return origin == GetFrontendIP()
	},
}
