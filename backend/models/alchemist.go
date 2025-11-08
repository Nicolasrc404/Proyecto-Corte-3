package models

import "gorm.io/gorm"

type Alchemist struct {
	gorm.Model
	PersonID  uint   `gorm:"not null;uniqueIndex"`
	Person    Person `gorm:"constraint:OnDelete:CASCADE"`
	Specialty string `gorm:"size:255"`
	Rank      string `gorm:"size:100"` // aprendiz, estatal, maestro...
}
