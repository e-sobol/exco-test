package fetcher

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func ExecuteHandler(ctx *gin.Context) {
	var urlFetchRequest UrlFetchRequest
	err := ReadBody(ctx, &urlFetchRequest)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	results := handleFetchUrlRequest(urlFetchRequest)

	ctx.JSON(http.StatusOK, UrlFetchResponse{Results: results})
}
