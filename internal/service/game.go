package service

import (
	"SurvivalGame/internal/model"
	"SurvivalGame/internal/utils/randUtils"
	"SurvivalGame/internal/utils/response"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
	"time"
)

func InitGame() string {
	ID := randUtils.RandNumString(16)
	mapSize := 30
	initMap := model.Map{
		MapSize: mapSize,
		Stat:    [][]int{},
	}
	initMap.Stat = make([][]int, mapSize)
	for i := 0; i < mapSize; i++ {
		initMap.Stat[i] = make([]int, mapSize)
	}
	initMap.Stat[0][0] = 20
	initMap.Stat[mapSize-1][mapSize-1] = 30
	initRound := model.Round{
		RedStatus: model.PlayerStatus{
			Score: 0,
		},
		BlueStatus: model.PlayerStatus{
			Score: 0,
		},
		Map: initMap,
	}
	model.NowRecord = model.Record{
		ID:         ID,
		Status:     "WAITING",
		RedPlayer:  "EMPTY",
		BluePlayer: "EMPTY",
		RoundDetail: map[int]model.Round{
			0: initRound,
		},
	}
	return ID
}

func JoinGame(c *gin.Context) {
	username := c.Param("username")
	if model.NowRecord.Status == "WAITING" {
		// 红方为空 或 上次活跃时间在3分钟以前 或 红方用户名等于username
		if model.NowRecord.RedPlayer == "EMPTY" ||
			model.LastActiveTime[model.NowRecord.RedPlayer].Before(time.Now().Add(-3*time.Minute)) ||
			model.NowRecord.RedPlayer == username {
			model.NowRecord.RedPlayer = username
			response.Write(c, "RED")
			return
		}
		// 蓝方为空或者上次活跃时间在3分钟以前 或 蓝方用户名等于username
		if model.NowRecord.BluePlayer == "EMPTY" || model.LastActiveTime[model.NowRecord.BluePlayer].Before(time.Now().Add(-3*time.Minute)) {
			model.NowRecord.BluePlayer = username
			response.Write(c, "BLUE")
			return
		}
	}
	// 返回加入失败
	response.Write(c, "FAIL")
	return
}

func CheckStatus(c *gin.Context) {
	UpdateBeginStatus()
	response.Write(c, model.NowRecord.Status)
}

func CheckRedPlayer(c *gin.Context) {
	response.Write(c, model.NowRecord.RedPlayer)
}

func CheckBluePlayer(c *gin.Context) {
	response.Write(c, model.NowRecord.BluePlayer)
}

func GenMapString(nowRoundNum int, record model.Record) string {
	result := ""
	round := record.RoundDetail[nowRoundNum]
	result += strconv.Itoa(round.Map.MapSize) + "\n"
	for i := 0; i < round.Map.MapSize; i++ {
		for j := 0; j < round.Map.MapSize; j++ {
			result += strconv.Itoa(round.Map.Stat[i][j]) + " "
		}
		result += "\n"
	}
	return result
}

func NowMap(c *gin.Context) {
	nowRoundNum := model.NowRecord.RoundsNum
	response.Write(c, GenMapString(nowRoundNum, model.NowRecord))
	return
}

func UpdateBeginStatus() {
	if model.NowRecord.BluePlayer != "EMPTY" && model.NowRecord.RedPlayer != "EMPTY" && model.NowRecord.Status == "WAITING" {
		model.NowRecord.Status = "BLUE"
	}
}

func UpdateStatus(c *gin.Context) {
	username := c.Query("username")
	op := c.Query("op")
	print(op)
	if username == model.NowRecord.RedPlayer && model.NowRecord.Status == "RED" {

	}
	if username == model.NowRecord.BluePlayer && model.NowRecord.Status == "BLUE" {

	}
	response.Write(c, "NOT YOUR TURN")
}

func UpdateMap(op string) {
	UpdateOPToMap(op)
}

func UpdateOPToMap(op string) {
	reader := strings.NewReader(op)
	var n, m int
	fmt.Fscanf(reader, "%d %d", &n, &m)
	for i := 1; i <= n; i++ {

	}
	for j := 1; j <= m; j++ {

	}
}

func FlushMap() {

}
