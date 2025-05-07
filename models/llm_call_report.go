package models

import "time"

type LLMCallReport struct {
	ID        uint `gorm:"primaryKey"`
	Entry     string
	Res       string
	Resp      string
	Amount    float64
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
