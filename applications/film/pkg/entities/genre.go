package entities

import "gorm.io/gorm"

type Genre struct {
	gorm.Model
	EnName string
	FaName string
}
