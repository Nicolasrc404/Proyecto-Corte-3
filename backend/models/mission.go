package models

import "gorm.io/gorm"

type Mission struct {
	gorm.Model
	Title       string
	Description string
	Difficulty  string
	Status      string `gorm:"default:pendiente"`
	AssignedTo  uint   // Alchemist ID
}
