package service

import (
	"SurvivalGame/internal/model"
	"SurvivalGame/internal/utils/sql"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"os"
)

func APINowAll(c *gin.Context) {
	c.JSONP(200, gin.H{"Status": model.NowRecord})
}

func APIGetRank(c *gin.Context) {
	var players []model.Player
	sql.Database.Order("ELO desc").Find(&players)
	c.JSONP(200, gin.H{"Rank": players})
}

func APIRecord(c *gin.Context) {
	id := c.Query("id")
	fd, err := os.Open("./saveGame/" + id + ".json")
	if err != nil {
		c.JSONP(200, gin.H{"Err": "Record no exist!"})
		c.Abort()
		return
	}
	data, err := io.ReadAll(fd)
	if err != nil {
		log.Fatal(err)
	}
	var record model.Record
	_ = json.Unmarshal(data, &record)
	c.JSONP(200, gin.H{"Status": record})
}

func APIRecentRecord(c *gin.Context) {
	var records []model.RecordSQL
	sql.Database.Order("start_time desc").Find(&records)
	c.JSONP(200, gin.H{"List": records})
}
