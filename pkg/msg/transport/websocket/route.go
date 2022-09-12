package websocket

import (
	"github.com/gin-gonic/gin"
)

func ApplyRoutes(r *gin.Engine, h Handler) {
	v1 := r.Group("/v1/")
	{
		v1.GET("live", h.GetLastEvent)
	}
}
