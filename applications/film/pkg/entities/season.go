package entities

import "gorm.io/gorm"

type Season struct {
	gorm.Model
	SeriesId     uint
	Title        string `gorm:"varchar(100)"`
	Presentation string `gorm:"varchar(300)"`
	BannerId     uint   `gorm:"check:banner_id > -1"`
	Banner       *Media `gorm:"foreignKey:BannerId"`
	Year         int    `gorm:"varchar(15)"`
	IMDBLink     string `gorm:"varchar(200)"`
	IMDB         float32
	Number       int
	Episodes     []*Episode `gorm:"foreignKey:SeasonId"`
}
