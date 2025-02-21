/**********************************/
/* @since 01/01/2025              */
/* @author K10s Open Source Team  */
/**********************************/
package handlers

import (
	"github.com/eliasdehondt/K10s/App/Backend/cmd/kubernetes"
	"github.com/gin-gonic/gin"
	"strconv"
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
