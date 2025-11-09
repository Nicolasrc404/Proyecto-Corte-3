package models

import "gorm.io/gorm"

type Alchemist struct {
	gorm.Model
	Name      string `gorm:"not null"`
	Age       int    `gorm:"not null"`
	Specialty string `gorm:"size:255"`
	Rank      string `gorm:"size:100"`
}
