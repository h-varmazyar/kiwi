package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/h-varmazyar/kiwi/applications/proxy/configs"
	"github.com/h-varmazyar/kiwi/applications/proxy/internal/handlers"
	"github.com/h-varmazyar/kiwi/applications/proxy/internal/repositories"
	db2 "github.com/h-varmazyar/kiwi/applications/proxy/pkg/db/PostgreSQL"
	log2 "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"strconv"
	"strings"
)

var (
	conf *Configs
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	log := new(log2.Logger)
	var err error
	if conf, err = prepareConfigs(ctx, log); err != nil {
		panic(err)
	}

	var db *db2.DB
	if db, err = db2.NewDatabase(ctx, conf.DB); err != nil {
		panic(err)
	}

	var postRepo *repositories.PostRepository
	if postRepo, err = repositories.NewPostRepository(ctx, log, db); err != nil {
		panic(err)
	}

	var b *bot.Bot
	if b, err = prepareBot(ctx, log, postRepo); err != nil {
		panic(err)
	}
	b.Start(ctx)
}

func prepareConfigs(_ context.Context, log *log2.Logger) (*Configs, error) {
	log.Infof("reding configs...")

	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		log.Warnf("failed to read from env: %v", err)
		viper.AddConfigPath("./configs")  //path for docker compose configs
		viper.AddConfigPath("../configs") //path for local configs
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		if err = viper.ReadInConfig(); err != nil {
			log.Warnf("failed to read from yaml: %v", err)
			localErr := viper.ReadConfig(bytes.NewBuffer(configs.DefaultConfig))
			if localErr != nil {
				log.WithError(localErr).Error("read from default configs failed")
				return nil, localErr
			}
		}
	}

	conf := new(Configs)
	if err := viper.Unmarshal(conf); err != nil {
		log.Errorf("faeiled unmarshal")
		return nil, err
	}

	conf.adminsId = make([]int64, 0)
	for _, s := range strings.Split(conf.Admins, ",") {
		id, err := strconv.ParseInt(strings.TrimSpace(s), 10, 64)
		if err != nil {
			return nil, err
		}
		conf.adminsId = append(conf.adminsId, id)
	}

	return conf, nil
}

func prepareBot(ctx context.Context, log *log2.Logger, postRepo *repositories.PostRepository) (*bot.Bot, error) {
	opts := []bot.Option{
		bot.WithDefaultHandler(defaultHandler),
		bot.WithMiddlewares(addLang, checkAdmin),
	}
	b, err := bot.New(conf.BotToken, opts...)
	if err != nil {
		log.WithError(err).Error("failed to create new bot")
		return nil, err
	}
	botHandlers, err := handlers.NewHandler(ctx, log, conf.Handlers, postRepo)
	if err != nil {
		log.WithError(err).Error("failed to create bot handlers")
		return nil, err
	}
	botHandlers.RegisterHandlers(ctx, b)

	return b, nil
}

func addLang(next bot.HandlerFunc) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		next(context.WithValue(ctx, "lang", "fa"), b, update)
	}
}

func checkAdmin(next bot.HandlerFunc) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		if update.Message != nil {
			for _, admin := range conf.adminsId {
				if update.Message.Chat.ID == admin {
					next(ctx, b, update)
					return
				}
			}
		}
	}
}

func defaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	fmt.Println("msg:", update.Message)
	fmt.Println("rpl:", update.Message.ReplyMarkup)
	fmt.Println("inq:", update.InlineQuery)
	fmt.Println("cap:", update.Message.CaptionEntities)
}
