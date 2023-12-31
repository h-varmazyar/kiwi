package main

import (
	gormext "github.com/h-varmazyar/gopack/gorm"
	"github.com/h-varmazyar/kiwi/applications/film/internal/handlers"
	redisPkg "github.com/h-varmazyar/kiwi/applications/film/pkg/db/redis"
	"github.com/h-varmazyar/kiwi/applications/film/pkg/imdb"
)

type Configs struct {
	Version  string           `yaml:"version"`
	BotToken string           `yaml:"botToken"`
	Admins   string           `yaml:"admins"`
	DB       gormext.Configs  `yaml:"db"`
	Redis    redisPkg.Configs `yaml:"redis"`
	Handlers handlers.Configs `yaml:"handlers"`
	IMDB     imdb.Configs     `yaml:"imdb"`
	admins   []int64          `yaml:"-"`
}
