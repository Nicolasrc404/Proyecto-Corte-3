package models

import "gorm.io/gorm"

type Audit struct {
	gorm.Model
	Action    string
	Entity    string
	EntityID  uint
	UserEmail string
}
