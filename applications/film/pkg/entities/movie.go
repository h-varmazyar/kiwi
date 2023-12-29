package entities

import (
	"encoding/json"
	"github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Language string

type Movie struct {
	gorm.Model
	Title        string         `gorm:"varchar(100)"`
	FaName       string         `gorm:"varchar(100)"`
	EnName       string         `gorm:"varchar(100)"`
	Presentation string         `gorm:"varchar(300)"`
	BannerId     uint           `gorm:"check:banner_id > -1"`
	Banner       *Media         `gorm:"foreignKey:BannerId"`
	Year         int            `gorm:"varchar(15)"`
	HardSubtitle bool           `gorm:"default:false"`
	SubtitleLink string         `gorm:"varchar(200)"`
	VisitCount   uint           `gorm:"default:0"`
	Languages    pq.StringArray `gorm:"type:text[]"`
	Tags         pq.StringArray `gorm:"type:text[]"`
	Genres       []*Genre       `gorm:"many2many:series_genres;"`
	Videos       []*Media       `gorm:"polymorphic:Owner;polymorphicValue:movie"`
	ImdbId       string         `gorm:"varchar(50)"`
	IMDB         float32
}

func (e *Movie) Json() string {
	d, err := json.Marshal(e)
	if err != nil {
		log.WithError(err).Errorf("cannot marshal movie: %v", e)
		return ""
	}

	return string(d)
}
