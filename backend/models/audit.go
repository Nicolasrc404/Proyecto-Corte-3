package models

import "gorm.io/gorm"

type Audit struct {
	gorm.Model
	Entity    string
	EntityID  uint
	Action    string
	Timestamp string
}
