package middleware

import (
	"goresume/config"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionToken, err := c.Cookie("session")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
			c.Abort()
			return
		}

		sessionKey := config.RedisKeyPrefix + "session:" + sessionToken
		userIDStr, err := config.RedisClient.Get(sessionKey).Result()
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid session"})
			c.Abort()
			return
		}

		userID, _ := strconv.ParseInt(userIDStr, 10, 64)
		c.Set("user_id", userID)
		c.Next()
	}
}
