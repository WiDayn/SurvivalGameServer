package service

import (
	"SurvivalGame/internal/model"
	"fmt"
	"strings"
)

var cost = []int{0, 2, 5, 20, 2, 5, 20}

func jue(x int) int {
	if x > 0 {
		return x
	} else {
		return -x
	}
}

func notRuleMove(id int, x int, y int, tox int, toy int, se int) int {
	id -= se * 3
	if id == 1 || id == 3 { //buBing+zhiMinZhe
		if jue(tox-x)+jue(toy-y) > 1 {
			return 1
		} else {
			return 0
		}
	} else { //qiBing
		if (jue(tox-x) == 1 && jue(toy-y) == 2) || (jue(tox-x) == 2 && jue(toy-y) == 1) {
			return 0
		} else {
			return 1
		}
	}
}

func TurnLoss(Loser string, Explain string) {
	if Loser == "RED" {
		model.NowRecord.FinalResult = "BLUE WIN"
		model.NowRecord.Explain = Explain
	} else {
		model.NowRecord.FinalResult = "RED WIN"
		model.NowRecord.Explain = Explain
	}
}

func UpdateOPToMap(op string) {
	var changed [233][233]int
	var se, coin, opLength, shadi int
	nowRound := model.NowRecord.RoundsNum
	dt := model.NowRecord.RoundDetail[nowRound]
	mapSize := dt.Map.MapSize
	mp := dt.Map.Stat
	reader := strings.NewReader(op)
	if model.NowRecord.Status == "RED" {
		se = 1
		coin = dt.RedStatus.Coin
	} else {
		se = 0
		coin = dt.BlueStatus.Coin
	}
	_, err := fmt.Fscanf(reader, "%d", &opLength)
	if err != nil {
		return
	}
	for i := 1; i <= opLength; i++ {
		var opt, x, y, tox, toy, id int
		_, err := fmt.Fscanf(reader, "%d", &opt)
		if err != nil {
			return
		}
		if opt == 1 {
			_, err := fmt.Fscanf(reader, "%d%d%d", &x, &y, &id)
			if err != nil {
				return
			}
			if mp[x][y]/10 != 2+se {
				TurnLoss(model.NowRecord.Status, "ERROR : This position can't produce soldiers!")
				return
			}
			if changed[x][y] == 1 {
				TurnLoss(model.NowRecord.Status, "ERROR : You can't move a soldier twice in one round!")
				return
			}
			changed[x][y] = 1
			coin -= cost[id]
			if coin < 0 {
				TurnLoss(model.NowRecord.Status, "ERROR : You don't have enough coin to buy this soldier!")
				return
			}
			mp[x][y] += id
			if model.NowRecord.Status == "RED" {
				dt.RedStatus.Score += 2
			} else {
				dt.BlueStatus.Score += 2
			}
		} else {
			_, err := fmt.Fscanf(reader, "%d%d%d%d", &x, &y, &tox, &toy)
			if err != nil {
				return
			}
			if mp[x][y]%10 < se*3+1 || (se+1)*3 < mp[x][y]%10 {
				TurnLoss(model.NowRecord.Status, "ERROR : There are no soldiers you can move in this position!")
				return
			}
			if se*3+1 <= mp[tox][toy]%10 && mp[tox][toy]%10 <= (se+1)*3 && mp[x][y]%10 != 3 && mp[x][y]%10 != 6 {
				TurnLoss(model.NowRecord.Status, "ERROR : You can't attack your own soldiers!")
				return
			}
			if tox < 0 || mapSize < tox || toy < 0 || mapSize < toy {
				TurnLoss(model.NowRecord.Status, "ERROR : soldiers can't be moved outside the map!")
				return
			}
			if notRuleMove(mp[x][y]%10, x, y, tox, toy, se) == 1 {
				TurnLoss(model.NowRecord.Status, "ERROR : This movement does not obey the rules!")
				return
			}
			if changed[x][y] == 1 {
				TurnLoss(model.NowRecord.Status, "ERROR : You can't move a soldier twice in one round!")
				return
			}
			changed[tox][toy] = 1
			if x == tox && y == toy {
				if mp[x][y]%10 != 3 && mp[x][y]%10 != 6 {
					TurnLoss(model.NowRecord.Status, "ERROR : This soldier is not a ZhiMinZhe, but Operation 2 did not change his position")
					return
				} else {
					if mp[x][y]%10 == 3 {
						mp[x][y] = 20
					} else {
						mp[x][y] = 30
					}
				}
				continue
			}
			if (se^1)*3+1 <= mp[tox][toy]%10 && mp[tox][toy]%10 <= ((se^1)+1)*3 {
				shadi++
			}
			if mp[tox][toy]/10 == 3-se {
				mp[tox][toy] = mp[x][y] % 10
			} else {
				mp[tox][toy] = mp[tox][toy]/10*10 + mp[x][y]%10
			}
			mp[x][y] = mp[x][y] / 10 * 10
		}
	}
	dt.Map.Stat = mp
	if model.NowRecord.Status == "RED" {
		dt.RedStatus.Coin = coin
		dt.RedStatus.Score += shadi * 5
		dt.RedStatus.Kill += shadi
	} else {
		dt.BlueStatus.Coin = coin
		dt.BlueStatus.Score += shadi * 5
		dt.BlueStatus.Kill += shadi
	}
	model.NowRecord.RoundDetail[nowRound] = dt
}
