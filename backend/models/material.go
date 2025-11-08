package models

import "gorm.io/gorm"

type Material struct {
	gorm.Model
	Name     string
	Category string
	Quantity float64
}
