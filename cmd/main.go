package main

import (
	"SurvivalGame/internal/model"
	"SurvivalGame/internal/server/http"
	"SurvivalGame/internal/service"
	"fmt"
)

func main() {
	player := model.Player{
		Username: "123",
	}
	fmt.Println(player)
	service.InitGame()
	http.StartHTTP()
}
