package middleware

import (
	"net/http"
	"oneinstack/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

var ipTokenMap = make(map[string]string)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
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
		claims, err := utils.ValidateJWT(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}
		// 获取当前 IP
		ip := c.ClientIP()
		// 获取上一次的 IP
		lastIP, exists := ipTokenMap[claims.Username]
		if exists && lastIP != ip {
			// 如果当前 IP 和上一次的 IP 不一样，需要重新验证
			claims, err = utils.ValidateJWT(parts[1])
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
				c.Abort()
				return
			}
		}
		// 将当前 IP 存储到 map 中
		ipTokenMap[claims.Username] = ip
		c.Set("username", claims.Username)
		c.Next()
	}
}
