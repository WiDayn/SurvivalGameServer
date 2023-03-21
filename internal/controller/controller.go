package controller

import "github.com/gin-gonic/gin"

func InitController(router *gin.Engine) {
	GameController(router)
	APIController(router)
}
