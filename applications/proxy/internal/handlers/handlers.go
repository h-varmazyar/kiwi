package handlers

import (
	"context"
	log "github.com/sirupsen/logrus"
)

type Configs struct {
	AdminChatId      int64
	RedisDB          int
	PublishChannelId int64
	ContentChannelId int64
	ProxyChannels    map[int64]struct{}
}

type Handler struct {
	configs *Configs
	log     *log.Logger
}

func NewHandler(ctx context.Context, log *log.Logger, configs *Configs) (*Handler, error) {
	return &Handler{
		configs: configs,
		log:     log,
	}, nil
}
