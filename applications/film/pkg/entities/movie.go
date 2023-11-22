package entities

import "gorm.io/gorm"

// type MovieType string
type Language string

//const (
//	MediaSeries MovieType = "سریال"
//	MediaMovie  MovieType = "سینمایی"
//)

type Movie struct {
	gorm.Model
	//Type         MovieType
	Title        string
	FaName       string
	EnName       string
	Genres       []*Genre
	Presentation string
	IMDB         float32
	Language     Language
	Banner       *Media
	Tags         string
	Year         int
	HardSubtitle bool
	SubtitleLink string
	Videos       []*Media
}
