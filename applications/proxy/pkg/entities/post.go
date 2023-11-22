package entities

import (
	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	FileId string
	Used   bool
}
