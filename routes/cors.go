package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// allowRemoteCors 补充允许 PATCH（mingfu_go_common 默认 CORS 未包含 PATCH）。
func allowRemoteCors(c *gin.Context) {
	origin := c.GetHeader("Origin")
	if origin != "" {
		c.Header("Access-Control-Allow-Origin", origin)
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token,Authorization,Token")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PATCH, OPTIONS")
	}
	if c.Request.Method == http.MethodOptions {
		c.AbortWithStatus(http.StatusNoContent)
		return
	}
	c.Next()
}
