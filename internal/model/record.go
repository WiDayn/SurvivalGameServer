package model

import (
	"time"
)

var NowRecord Record
var LastActiveTime = map[string]time.Time{}

type Map struct {
	MapSize int     // 地图大小
	Stat    [][]int // 0 - 空气, 1 - 红方, 2- 蓝方, 3 - 子弹, 4 - 食物, 5- 墙, 6 - 河流
}

type PlayerStatus struct {
	Score int
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
	FinalResult int           // 最终结果
	RedPlayer   string        // 红方玩家ID
	BluePlayer  string        // 蓝色玩家ID
	RoundDetail map[int]Round // 每一回合的细节
}
