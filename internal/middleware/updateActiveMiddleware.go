package middleware

import (
	"SurvivalGame/internal/model"
	"SurvivalGame/internal/utils/response"
	"github.com/gin-gonic/gin"
	"time"
)

func UpdateActiveMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		username := c.Query("username")
		if username == "" || username == "EMPTY" {
			response.Write(c, "USERNAME ERROR")
			c.Abort()
		}
		model.LastActiveTime[username] = time.Now()
		c.Next()
	}
}
