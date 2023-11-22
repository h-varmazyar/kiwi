package entities

import (
	"encoding/json"
	"github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Episode struct {
	gorm.Model
	Title        string         `gorm:"varchar(100)"`
	Presentation string         `gorm:"varchar(300)"`
	BannerId     uint           `gorm:"check:banner_id > -1"`
	Banner       *Media         `gorm:"foreignKey:BannerId"`
	HardSubtitle bool           `gorm:"default:false"`
	SubtitleLink string         `gorm:"varchar(200)"`
	Number       int            `gorm:"check:banner_id > 0"`
	Videos       []*Media       `gorm:"polymorphic:Owner;polymorphicValue:episode"`
	Tags         pq.StringArray `gorm:"type:text[]"`
	Languages    pq.StringArray `gorm:"type:text[]"`
	VisitCount   uint           `gorm:"default:0"`
	IMDBLink     string         `gorm:"varchar(200)"`
	IMDB         float32
	SeasonId     uint
}

func (e *Episode) Json() string {
	d, err := json.Marshal(e)
	if err != nil {
		log.WithError(err).Errorf("cannot marshal episode: %v", e)
		return ""
	}

	return string(d)
}
