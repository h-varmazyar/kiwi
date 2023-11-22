package entities

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Series struct {
	gorm.Model
	Title            string    `gorm:"varchar(100)"`
	EnName           string    `gorm:"varchar(100)"`
	FaName           string    `gorm:"varchar(100)"`
	Presentation     string    `gorm:"varchar(300)"`
	BannerId         uint      `gorm:"check:banner_id > -1"`
	Banner           *Media    `gorm:"foreignKey:BannerId"`
	Year             string    `gorm:"varchar(15)"`
	Genres           []*Genre  `gorm:"many2many:series_genres;"`
	Seasons          []*Season `gorm:"foreignKey:SeriesId"`
	IMDBLink         string    `gorm:"varchar(200)"`
	IMDB             float32
	PrivateChannelId int64
}

func (s *Series) Json() string {
	d, err := json.Marshal(s)
	if err != nil {
		log.WithError(err).Errorf("cannot marshal series: %v", s)
		return ""
	}

	return string(d)
}
