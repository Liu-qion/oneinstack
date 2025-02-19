package ssh

import (
	"net/http"
	"oneinstack/internal/services/ssh"
	"oneinstack/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

func OpenSSH(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no token"})
		c.Abort()
		return
	}
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token format"})
		c.Abort()
		return
	}
	_, err := utils.ValidateJWT(parts[1])
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		c.Abort()
		return
	}
	ssh.OpenWebShell(c)
}
