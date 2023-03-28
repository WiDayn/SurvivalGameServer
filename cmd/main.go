package main

import (
	"SurvivalGame/internal/server/http"
	"SurvivalGame/internal/service"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().Unix())
	service.InitGame()
	service.SaveGameDetail()
	http.StartHTTP()
}
