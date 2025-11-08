package models

import "gorm.io/gorm"

type Transmutation struct {
	gorm.Model
	AlchemistID uint
	MaterialID  uint
	Formula     string
	Status      string `gorm:"default:en_proceso"`
	Result      string
}
