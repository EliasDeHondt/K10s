/**********************************/
/* @since 01/01/2025              */
/* @author K10s Open Source Team  */
/**********************************/
package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetStatsHandler(ctx *gin.Context) {
	name, ok := ctx.GetQuery("nodeName")
	if ok {
		GetStatsForNode(ctx, name)
	} else {
		GetTotalStats(ctx)
	}
}

func GetTotalStats(ctx *gin.Context) {
	metrics, err := c.GetTotalUsage()
	if err != nil {
		fmt.Printf("GetTotalStats err: %v\n", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "An error has occurred or the request has been timed out."})
		return
	}
	ctx.JSON(http.StatusOK, metrics)
}

func GetStatsForNode(ctx *gin.Context, name string) {
	metrics, err := c.GetUsageForNode(name)
	if err != nil {
		fmt.Printf("GetStatsForNode err: %v\n", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "An error has occurred or the request has been timed out."})
		return
	}
	ctx.JSON(http.StatusOK, metrics)
}
