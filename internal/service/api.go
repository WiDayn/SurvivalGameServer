package service

import (
	"SurvivalGame/internal/model"
	"SurvivalGame/internal/utils/sql"
	"github.com/gin-gonic/gin"
)

func APINowAll(c *gin.Context) {
	c.JSONP(200, gin.H{"Status": model.NowRecord})
}

func APIGetRank(c *gin.Context) {
	var players []model.Player
	sql.Database.Order("ELO desc").Find(&players)
	c.JSONP(200, gin.H{"Rank": players})
}
