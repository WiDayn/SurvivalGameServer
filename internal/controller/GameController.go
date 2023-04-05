package controller

import (
	"SurvivalGame/internal/middleware"
	"SurvivalGame/internal/service"
	"github.com/gin-gonic/gin"
)

func GameController(router *gin.Engine) {
	router.GET("/JoinGame", middleware.UpdateActiveMiddleware(), service.JoinGame)
	router.GET("/Status", middleware.UpdateActiveMiddleware(), service.CheckStatus)
	router.GET("/NowMap", middleware.UpdateActiveMiddleware(), service.NowMap)
	router.GET("/NowCoin", middleware.UpdateActiveMiddleware(), service.NowCoin)
	router.GET("/CheckRedPlayer", middleware.UpdateActiveMiddleware(), service.CheckRedPlayer)
	router.GET("/CheckBluePlayer", middleware.UpdateActiveMiddleware(), service.CheckBluePlayer)
	router.GET("/UpdateStatus", middleware.UpdateActiveMiddleware(), service.UpdateStatus)
	router.GET("/GetResult", middleware.UpdateActiveMiddleware(), service.GetResult)
}
