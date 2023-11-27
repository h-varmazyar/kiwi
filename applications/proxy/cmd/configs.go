package main

import (
	gormext "github.com/h-varmazyar/gopack/gorm"
	"github.com/h-varmazyar/kiwi/applications/proxy/internal/handlers"
)

type Configs struct {
	Version  string           `yaml:"version"`
	BotToken string           `yaml:"botToken"`
	Admins   string           `yaml:"admins"`
	DB       gormext.Configs  `yaml:"db"`
	Handlers handlers.Configs `yaml:"handlers"`
	adminsId []int64          `yaml:"-"`
}
