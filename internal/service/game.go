package service

import (
	"SurvivalGame/internal/model"
	"SurvivalGame/internal/utils/randUtils"
	"SurvivalGame/internal/utils/response"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func InitGame() string {
	ID := randUtils.RandNumString(16)
	mapSize := 10
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
			Coin:  5,
		},
		BlueStatus: model.PlayerStatus{
			Score: 0,
			Coin:  5,
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
		Seed: rand.Int63(),
	}
	rand.Seed(model.NowRecord.Seed)
	FlushMap()
	return ID
}

func JoinGame(c *gin.Context) {
	username := c.Query("username")
	if model.NowRecord.Status == "WAITING" {
		// 红方为空 或 上次活跃时间在3分钟以前 或 红方用户名等于username
		if model.NowRecord.RedPlayer == "EMPTY" ||
			model.LastActiveTime[model.NowRecord.RedPlayer].Before(time.Now().Add(-time.Minute)) ||
			model.NowRecord.RedPlayer == username {
			model.NowRecord.RedPlayer = username
			response.Write(c, "RED")
			return
		}
		// 蓝方为空或者上次活跃时间在3分钟以前 或 蓝方用户名等于username
		if model.NowRecord.BluePlayer == "EMPTY" || model.LastActiveTime[model.NowRecord.BluePlayer].Before(time.Now().Add(-time.Minute)) {
			model.NowRecord.BluePlayer = username
			response.Write(c, "BLUE")
			return
		}
	} else {
		if model.NowRecord.RedPlayer == username {
			response.Write(c, "RED")
			return
		}
		if model.NowRecord.BluePlayer == username {
			response.Write(c, "BLUE")
			return
		}
	}
	// 返回加入失败
	response.Write(c, "FAIL")
	return
}

func CheckStatus(c *gin.Context) {
	KickOfflinePlayer()
	detectTimeOut()
	UpdateBeginStatus()
	response.Write(c, model.NowRecord.Status)
}

func KickOfflinePlayer() {
	// 清除三分钟没有汇报过的玩家
	if model.NowRecord.RedPlayer != "EMPTY" &&
		model.LastActiveTime[model.NowRecord.RedPlayer].Before(time.Now().Add(-time.Minute)) {
		model.NowRecord.RedPlayer = "EMPTY"
	}
	if model.NowRecord.BluePlayer != "EMPTY" &&
		model.LastActiveTime[model.NowRecord.BluePlayer].Before(time.Now().Add(-time.Minute)) {
		model.NowRecord.BluePlayer = "EMPTY"
	}
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
			result += strconv.Itoa(round.Map.Stat[i][j])
			if j != round.Map.MapSize-1 {
				result += " "
			}
		}
		if i != round.Map.MapSize-1 {
			result += "\n"
		}
	}
	return result
}

func NowMap(c *gin.Context) {
	nowRoundNum := model.NowRecord.RoundsNum
	response.Write(c, GenMapString(nowRoundNum, model.NowRecord))
	return
}

func NowCoin(c *gin.Context) {
	username := c.Query("username")
	nowRoundNum := model.NowRecord.RoundsNum
	if model.NowRecord.BluePlayer == username {
		response.Write(c, strconv.Itoa(model.NowRecord.RoundDetail[nowRoundNum].BlueStatus.Coin))
	}
	if model.NowRecord.RedPlayer == username {
		response.Write(c, strconv.Itoa(model.NowRecord.RoundDetail[nowRoundNum].RedStatus.Coin))
	}
	response.Write(c, "YOU NOT IN GAME")
}

func GetResult(c *gin.Context) {
	response.Write(c, model.LastResult)
}

func UpdateBeginStatus() {
	if model.NowRecord.BluePlayer != "EMPTY" && model.NowRecord.RedPlayer != "EMPTY" && model.NowRecord.Status == "WAITING" {
		model.NowRecord.StartTime = time.Now()
		model.NowRecord.Status = "BLUE"
	}
}

func UpdateStatus(c *gin.Context) {
	username := c.Query("username")
	op := c.Query("op")
	if username == model.NowRecord.RedPlayer && model.NowRecord.Status == "RED" {
		fmt.Println("RED MOVE:" + op)
		UpdateOPToMap(op)
		detectWin()
		if model.NowRecord.FinalResult != "" {
			// 进入游戏结算
			settleGame()
			return
		}

		var newRound *model.Round

		round := model.NowRecord.RoundDetail[model.NowRecord.RoundsNum]
		newRound = round.DeepCopyRound()
		// 红方返回的地图要存入下一个round
		model.NowRecord.RoundDetail[model.NowRecord.RoundsNum+1] = *newRound
		// 然后让当前回合+1
		model.NowRecord.RoundsNum += 1
		// 红方移动完后要刷新地图和金钱信息(直接写入下一个回合)
		FlushMap()
		FlushMoney()

		model.NowRecord.Status = "BLUE"
	}
	if username == model.NowRecord.BluePlayer && model.NowRecord.Status == "BLUE" {
		fmt.Println("BLUE MOVE:" + op)
		UpdateOPToMap(op)
		detectWin()
		if model.NowRecord.FinalResult != "" {
			// 进入游戏结算
			settleGame()
			return
		}
		model.NowRecord.Status = "RED"
	}
	response.Write(c, "NOT YOUR TURN")
}

func detectWin() {
	if model.NowRecord.FinalResult != "" {
		return
	}
	nowRound := model.NowRecord.RoundsNum
	df := model.NowRecord.RoundDetail[nowRound]
	for i := 0; i < df.Map.MapSize; i++ {
		for j := 0; j < df.Map.MapSize; j++ {
			// 检测到对方还有Base
			if model.NowRecord.Status == "RED" && (df.Map.Stat[i][j]/10)%10 == 2 {
				return
			}
			if model.NowRecord.Status == "BLUE" && (df.Map.Stat[i][j]/10)%10 == 3 {
				return
			}
		}
	}
	if model.NowRecord.Status == "RED" {
		model.NowRecord.FinalResult = "RED WIN"
		model.NowRecord.Explain = "CLEAR ALL ENEMY BASE"
	} else {
		model.NowRecord.FinalResult = "BLUE WIN"
		model.NowRecord.Explain = "CLEAR ALL ENEMY BASE"
	}
}

func settleGame() {
	model.LastResult = model.NowRecord.FinalResult
	if model.NowRecord.FinalResult == "BLUE WIN" {
		UpdateRank(model.NowRecord.BluePlayer, model.NowRecord.RedPlayer, 1)
		SaveGameDetail()
	}
	if model.NowRecord.FinalResult == "RED WIN" {
		UpdateRank(model.NowRecord.BluePlayer, model.NowRecord.RedPlayer, -1)
		SaveGameDetail()
	}
	InsertRecord(model.NowRecord)
	InitGame()
}

func FlushMoney() {
	nowRound := model.NowRecord.RoundsNum
	rd := model.NowRecord.RoundDetail[nowRound]
	for i := 0; i < rd.Map.MapSize; i++ {
		for j := 0; j < rd.Map.MapSize; j++ {
			if (rd.Map.Stat[i][j]/10)%10 == 1 {
				if rd.Map.Stat[i][j]%10 != 0 && rd.Map.Stat[i][j]%10 <= 3 {
					rd.BlueStatus.Coin += 1
					rd.BlueStatus.Score += 1
				}
				if rd.Map.Stat[i][j]%10 > 3 {
					rd.RedStatus.Coin += 1
					rd.RedStatus.Score += 1
				}
			}
		}
	}
	model.NowRecord.RoundDetail[nowRound] = rd
}

func FlushMap() {
	nowRound := model.NowRecord.RoundsNum
	mp := model.NowRecord.RoundDetail[nowRound].Map
	existMine := 0
	for i := 0; i < mp.MapSize; i++ {
		for j := 0; j < mp.MapSize; j++ {
			if (mp.Stat[i][j]/10)%10 == 1 {
				leftTime := mp.Stat[i][j] / 100
				leftTime--
				if leftTime > 0 {
					// 更新剩余时间
					mp.Stat[i][j] = (mp.Stat[i][j] % 10) + 10 + leftTime*100
					existMine++
				} else {
					// 变为空地
					mp.Stat[i][j] = (mp.Stat[i][j] % 10) + 0
				}
			}
		}
	}

	needMine := rand.Int()%10 + 2
	// 场上金矿不够，补充
	for existMine < needMine {
		i := rand.Int() % mp.MapSize
		j := rand.Int() % mp.MapSize
		// 找空地
		for (mp.Stat[i][j]/10)%10 != 0 {
			i = rand.Int() % mp.MapSize
			j = rand.Int() % mp.MapSize
		}
		mp.Stat[i][j] = (mp.Stat[i][j] % 10) + 10 + (rand.Int()%50+20)*100
		existMine++
	}
}

func SaveGameDetail() {
	marshal, _ := json.Marshal(model.NowRecord)

	newFile, _ := os.Create("./saveGame/" + model.NowRecord.ID + ".json")
	_, err := newFile.Write(marshal)
	if err != nil {
		log.Printf(err.Error())
		return
	}
	err = newFile.Close()
	if err != nil {
		log.Printf(err.Error())
		return
	}
}

func detectTimeOut() {
	if model.NowRecord.Status == "WAITING" {
		return
	}
	if model.LastActiveTime[model.NowRecord.RedPlayer].Before(time.Now().Add(-30 * time.Second)) {
		fmt.Println("RED TIME OUT")
		model.NowRecord.FinalResult = "BLUE WIN"
		model.NowRecord.Explain = "RED TIME OUT"
	}
	if model.LastActiveTime[model.NowRecord.BluePlayer].Before(time.Now().Add(-30 * time.Second)) {
		fmt.Println("BLUE TIME OUT")
		model.NowRecord.FinalResult = "RED WIN"
		model.NowRecord.Explain = "BLUE TIME OUT"
	}
	if model.NowRecord.FinalResult != "" {
		// 进入游戏结算
		settleGame()
		return
	}
}
