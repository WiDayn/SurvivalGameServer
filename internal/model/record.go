package model

import (
	"SurvivalGame/internal/utils/sql"
	"time"
)

var NowRecord Record
var LastActiveTime = map[string]time.Time{}
var LastResult string

type Map struct {
	MapSize int     // 地图大小
	Stat    [][]int // 0 - 空气, 1 - 红方, 2- 蓝方, 3 - 子弹, 4 - 食物, 5- 墙, 6 - 河流
}

type PlayerStatus struct {
	Score int
	Coin  int
	Kill  int
}

type Round struct {
	RedStatus  PlayerStatus // 红方玩家的情况
	BlueStatus PlayerStatus // 蓝方玩家的情况
	Map        Map          // 地图状态
}

type Record struct {
	ID          string        // 游戏的ID
	RoundsNum   int           // 回合数
	Status      string        // 当前状态
	StartTime   time.Time     // 开始游戏的时间
	FinalResult string        // 最终结果
	Explain     string        // 对于结果的解释
	RedPlayer   string        // 红方玩家ID
	BluePlayer  string        // 蓝色玩家ID
	RoundDetail map[int]Round // 每一回合的细节
	Seed        int64         // 地图种子
}

type RecordSQL struct {
	ID          string    // 游戏的ID
	RoundsNum   int       // 回合数
	StartTime   time.Time // 开始游戏的时间
	FinalResult string    // 最终结果
	Explain     string    // 对于结果的解释
	RedPlayer   string    // 红方玩家ID
	BluePlayer  string    // 蓝色玩家ID
	Seed        int64     // 地图种子
}

func (r *Round) DeepCopyRound() *Round {
	// 深拷贝 Map 结构体
	newMap := Map{
		MapSize: r.Map.MapSize,
		Stat:    make([][]int, len(r.Map.Stat)),
	}
	for i := range r.Map.Stat {
		newMap.Stat[i] = make([]int, len(r.Map.Stat[i]))
		copy(newMap.Stat[i], r.Map.Stat[i])
	}

	// 深拷贝 PlayerStatus 结构体
	newRedStatus := PlayerStatus{
		Score: r.RedStatus.Score,
		Coin:  r.RedStatus.Coin,
		Kill:  r.RedStatus.Kill,
	}
	newBlueStatus := PlayerStatus{
		Score: r.BlueStatus.Score,
		Coin:  r.BlueStatus.Coin,
		Kill:  r.BlueStatus.Kill,
	}

	// 构造新的 Round 结构体
	newRound := Round{
		RedStatus:  newRedStatus,
		BlueStatus: newBlueStatus,
		Map:        newMap,
	}

	return &newRound
}

func init() {
	if err := sql.Database.AutoMigrate(&RecordSQL{}); err != nil {
	}
}
