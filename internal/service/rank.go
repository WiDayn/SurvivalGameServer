package service

import (
	"SurvivalGame/internal/model"
	"SurvivalGame/internal/utils/sql"
	"math"
)

// SA 1 -> A Win, 0.5 -> Draw, 0 -> B Win
func calculateNewRank(A float64, B float64, SA float64) (float64, float64) {
	var K float64
	if math.Min(A, B) < 2000 {
		K = 30
	} else if math.Min(A, B) < 2400 {
		K = 130 - math.Min(A, B)/20
	} else {
		K = 10
	}

	EA := 1 / (1 + math.Pow(10, (B-A)/400))
	EB := 1 / (1 + math.Pow(10, (A-B)/400))

	RA := A + K*(SA-EA)
	RB := B + K*((1-SA)-EB)

	return RA, RB
}

// UpdateRank winStat 1 -> A Win, 0 -> Draw, -1 -> B Win
func UpdateRank(usernameA string, usernameB string, winStat int) {
	var playerA, playerB model.Player
	sql.Database.First(&playerA, "username = ?", usernameA)
	sql.Database.First(&playerB, "username = ?", usernameB)

	if playerA.Username == "" {
		playerA.Username = usernameA
		playerA.ELO = 800
	}
	if playerB.Username == "" {
		playerB.Username = usernameB
		playerB.ELO = 800
	}
	var SA float64
	if winStat == 1 {
		playerA.Win++
		playerB.Loss++
		SA = 1.0
	}
	if winStat == 0 {
		playerA.Draw++
		playerB.Draw++
		SA = 0.5
	}
	if winStat == -1 {
		playerA.Loss++
		playerB.Win++
		SA = 0
	}

	playerA.ELO, playerB.ELO = calculateNewRank(playerA.ELO, playerB.ELO, SA)

	sql.Database.Model(&model.Player{
		Username: playerA.Username,
	}).Updates(&playerA)
	sql.Database.Model(&model.Player{
		Username: playerB.Username,
	}).Updates(&playerB)
}
