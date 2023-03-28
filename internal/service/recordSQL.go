package service

import (
	"SurvivalGame/internal/model"
	"SurvivalGame/internal/utils/sql"
)

func InsertRecord(record model.Record) {
	sql.Database.Create(&model.RecordSQL{
		ID:          record.ID,
		RoundsNum:   record.RoundsNum,
		StartTime:   record.StartTime,
		FinalResult: record.FinalResult,
		Explain:     record.Explain,
		RedPlayer:   record.RedPlayer,
		BluePlayer:  record.BluePlayer,
		Seed:        record.Seed,
	})
}
