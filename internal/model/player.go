package model

import (
	"SurvivalGame/internal/utils/sql"
	"gorm.io/gorm"
)

type Player struct {
	gorm.Model
	Username string
	Win      int
	Draw     int
	Loss     int
	ELO      float64
}

func init() {
	if err := sql.Database.AutoMigrate(&Player{}); err != nil {
	}
}
