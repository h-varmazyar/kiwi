package entities

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type MediaType string
type MediaQuality string

const (
	MediaTypePicture MediaType = "picture"
	MediaTypeVideo   MediaType = "video"
)

const (
	Quality480HQ      MediaQuality = "480p HQ"
	Quality720HQ      MediaQuality = "720p HQ"
	Quality1080HQ     MediaQuality = "1080p HQ"
	Quality480BluRay  MediaQuality = "480p BluRay"
	Quality720BluRay  MediaQuality = "720p BluRay"
	Quality1080BluRay MediaQuality = "1080p BluRay"
	Quality480WebDL   MediaQuality = "480p WEB-DL"
	Quality720WebDL   MediaQuality = "720p WEB-DL"
	Quality1080WebDL  MediaQuality = "1080p WEB-DL"
)

type Media struct {
	gorm.Model
	DownloadUrl    string       `gorm:"varchar(200)"`
	TelegramFileId string       `gorm:"varchar(100)"`
	Quality        MediaQuality `gorm:"varchar(50)"`
	Type           MediaType    `gorm:"varchar(50)"`
	OwnerType      string       `gorm:"not null"`
	OwnerID        uint
}

func (e *Media) Json() string {
	d, err := json.Marshal(e)
	if err != nil {
		log.WithError(err).Errorf("cannot marshal media: %v", e)
		return ""
	}

	return string(d)
}
