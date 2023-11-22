package main

import (
	"context"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	gormext "github.com/h-varmazyar/gopack/gorm"
	"github.com/h-varmazyar/kiwi/applications/film/internal/handlers"
	dbPkg "github.com/h-varmazyar/kiwi/applications/film/pkg/db/PostgreSQL"
	redisPkg "github.com/h-varmazyar/kiwi/applications/film/pkg/db/redis"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
)

const (
	defaultMsg = "خطا در پردازش پیام"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		bot.WithDefaultHandler(defaultHandler),
		bot.WithMiddlewares(addLang),
	}

	b, err := bot.New("1972400420:AAHtOEYNhbaNz-Xm27MEP1BYa0cB2QG7APw", opts...)
	if err != nil {
		panic(err)
	}

	redisConf := &redisPkg.Configs{
		Host:     "localhost",
		Port:     6379,
		Password: "",
	}
	conf := &handlers.Configs{
		AdminChatId:     1689926,
		PublicChannelId: -1001803420363,
		RedisDB:         1,
		RedisConfigs:    redisConf,
	}

	dbConf := gormext.Configs{
		DbType:      gormext.PostgreSQL,
		Port:        5433,
		Host:        "localhost",
		Username:    "postgres",
		Password:    "postgres",
		Name:        "kiwi",
		IsSSLEnable: false,
	}

	db, err := dbPkg.NewDatabase(ctx, dbConf)
	if err != nil {
		panic(err)
	}

	botHandlers, err := handlers.NewHandler(ctx, new(log.Logger), conf, db)
	if err != nil {
		panic(err)
	}
	botHandlers.RegisterHandlers(ctx, b)

	b.Start(ctx)
}

func defaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   defaultMsg,
	})
	if err != nil {
		log.WithError(err).Error("failed to send default message")
	}
}

func addLang(next bot.HandlerFunc) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		next(context.WithValue(ctx, "lang", "fa"), b, update)
	}
}
