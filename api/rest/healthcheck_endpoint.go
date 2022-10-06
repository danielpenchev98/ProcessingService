package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Endpoint used for availability check for the services
func CheckHealth(c *gin.Context) {
	c.Status(http.StatusOK)
}
