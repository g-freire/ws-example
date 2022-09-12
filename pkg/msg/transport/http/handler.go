package http

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func handleVersion(c *gin.Context) {
	c.JSON(http.StatusOK, "WS API - 2022")
}
