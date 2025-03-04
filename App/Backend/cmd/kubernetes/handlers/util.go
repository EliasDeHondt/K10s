/**********************************/
/* @since 01/01/2025              */
/* @author K10s Open Source Team  */
/**********************************/
package handlers

import (
	"context"
	"github.com/eliasdehondt/K10s/App/Backend/cmd/kubernetes"
	"github.com/gin-gonic/gin"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
	"strconv"
	"time"
)

var c = kubernetes.GetClients()

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

	svc, err := c.GetServices("k10s-namespaces").Get(ctx, "k10s-ingress-service", metav1.GetOptions{})

	if err != nil {
		log.Printf("Failed to get k10s-ingress-service: %v, using default development url", err)
		return "http://localhost:4200"
	}

	if len(svc.Status.LoadBalancer.Ingress) > 0 {

		protocol := "http"
		if svc.Status.LoadBalancer.Ingress[0].Hostname != "" {
			protocol = "https"
		}

		url := protocol + "://" + svc.Status.LoadBalancer.Ingress[0].IP
		if svc.Spec.Ports[0].Port != 80 && svc.Spec.Ports[0].Port != 443 {
			url += ":" + strconv.Itoa(int(svc.Spec.Ports[0].Port))
		}

		return url
	}

	return "http://localhost:4200"
}
