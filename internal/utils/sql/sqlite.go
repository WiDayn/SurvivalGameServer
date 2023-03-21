package sql

import (
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var Database *gorm.DB

func init() {
	if tempDatabase, err := gorm.Open(sqlite.Open("SG.db"), &gorm.Config{}); err != nil {
		panic(err)
	} else {
		Database = tempDatabase
	}
}
