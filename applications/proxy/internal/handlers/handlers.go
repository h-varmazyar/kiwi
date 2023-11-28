package handlers

import (
	"context"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/h-varmazyar/kiwi/applications/proxy/internal/repositories"
	log "github.com/sirupsen/logrus"
	"regexp"
	"strings"
)

var (
	startRegex = regexp.MustCompile(`^/start [a-zA-Z0-9]*`)

	isValidContent = func(update *models.Update) bool {
		return update.Message.Photo != nil && len(update.Message.Photo) > 0
	}

	isValidProxy = func(update *models.Update) bool {
		msgWithProxy := strings.Contains(update.Message.Text, "https://t.me/proxy?server")
		callbackWithProxy := update.Message != nil && update.Message.ReplyMarkup.InlineKeyboard != nil
		return msgWithProxy || callbackWithProxy
	}
)

type Configs struct {
	PublishChannelId int64 `yaml:"publishChannelId"`
}

type Handler struct {
	configs  Configs
	log      *log.Logger
	postRepo *repositories.PostRepository
}

func NewHandler(_ context.Context, log *log.Logger, configs Configs, postRepo *repositories.PostRepository) (*Handler, error) {
	return &Handler{
		configs:  configs,
		log:      log,
		postRepo: postRepo,
	}, nil
}

func (h *Handler) RegisterHandlers(_ context.Context, b *bot.Bot) {
	b.RegisterHandlerRegexp(bot.HandlerTypeMessageText, startRegex, h.handleStartCmd)
	b.RegisterHandlerMatchFunc(isValidContent, h.handleMedia)
	b.RegisterHandlerMatchFunc(isValidProxy, h.handleProxy)
}

func (h *Handler) handleStartCmd(ctx context.Context, bot *bot.Bot, update *models.Update) {

}
