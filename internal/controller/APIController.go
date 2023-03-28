package controller

import (
	"SurvivalGame/internal/service"
	"github.com/gin-gonic/gin"
)

func APIController(router *gin.Engine) {
	router.GET("/NowAll", service.APINowAll)
	router.GET("/GetRank", service.APIGetRank)
	router.GET("/Record", service.APIRecord)
	router.GET("/RecordList", service.APIRecentRecord)
}
