package main

import (
	"context"
	"fmt"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	gormext "github.com/h-varmazyar/gopack/gorm"
	"github.com/h-varmazyar/kiwi/applications/proxy/internal/repositories"
	db2 "github.com/h-varmazyar/kiwi/applications/proxy/pkg/db/PostgreSQL"
	"github.com/h-varmazyar/kiwi/applications/proxy/pkg/entities"
	log2 "github.com/sirupsen/logrus"
	"os"
	"os/signal"
)

const (
	defaultMsg = "خطا در پردازش پیام"
)

type Configs struct {
	AdminChatId      int64
	PublishChannelId int64
	ContentChannelId int64
	TelMtProto       int64
	PowerfulProxy    int64
}

var (
	conf     *Configs
	postRepo *repositories.PostRepository
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	log := new(log2.Logger)
	conf = &Configs{
		AdminChatId:      1689926,
		PublishChannelId: -1002018115271,
		ContentChannelId: -1002132165405,
		TelMtProto:       -1002077953578,
		PowerfulProxy:    -1002077953578,
	}

	dbConf := gormext.Configs{
		DbType:      gormext.PostgreSQL,
		Port:        5433,
		Host:        "localhost",
		Username:    "postgres",
		Password:    "postgres",
		Name:        "proxy",
		IsSSLEnable: false,
	}

	db, err := db2.NewDatabase(ctx, dbConf)
	if err != nil {
		panic(err)
	}

	postRepo, err = repositories.NewPostRepository(ctx, log, db)
	if err != nil {
		panic(err)
	}

	opts := []bot.Option{
		bot.WithDefaultHandler(defaultHandler),
		bot.WithMiddlewares(addLang),
	}

	b, err := bot.New("1050203162:AAG298gFzWIn7ZS_2WAnnM8QglAFPG9tueQ", opts...)
	if err != nil {
		panic(err)
	}

	//conf := &handlers.Configs{
	//	AdminChatId:     1689926,
	//	PublicChannelId: -1001803420363,
	//	RedisDB:         1,
	//	RedisConfigs:    redisConf,
	//}
	//
	//botHandlers, err := handlers.NewHandler(ctx, new(log.Logger), conf)
	//if err != nil {
	//	panic(err)
	//}
	//botHandlers.RegisterHandlers(ctx, b)

	b.Start(ctx)
}

func defaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	proxyLinks := make([]string, 0)
	if update.ChannelPost != nil {
		switch update.ChannelPost.Chat.ID {
		case conf.TelMtProto,
			conf.PowerfulProxy:
			if k := update.ChannelPost.ReplyMarkup.InlineKeyboard; k != nil {
				if len(k) == 0 {
					return
				}
				if len(k[0]) == 0 {
					return
				}
				for _, button := range k[0] {
					proxyLinks = append(proxyLinks, button.URL)
				}
			}

			if len(proxyLinks) > 0 {
				photo, err := postRepo.NewUnused(ctx)
				if err != nil {
					fmt.Println("unused post:", err)
					return
				}
				photoFile := &models.InputFileString{
					Data: photo.FileId,
				}
				proxyBtn := models.InlineKeyboardMarkup{
					InlineKeyboard: make([][]models.InlineKeyboardButton, 0),
				}
				row := make([]models.InlineKeyboardButton, 0)
				proxyLinksText := ""
				for i, link := range proxyLinks {
					proxyLinksText = fmt.Sprintf("%v *[اتصال به پروکسی ✅](%v)*\n", proxyLinksText, link)
					btn := models.InlineKeyboardButton{
						Text: "اتصال ✅",
						URL:  link,
					}
					row = append(row, btn)
					if i%2 == 1 || i == len(proxyLinks)-1 {
						proxyBtn.InlineKeyboard = append(proxyBtn.InlineKeyboard, row)
						row = make([]models.InlineKeyboardButton, 0)
					}
				}

				caption := `
☑️ پروکسی ضد فیلتر و پر سرعت

🔘لطفا پروکسی ها را برای دوستان خود هم ارسال کنید تا استفاده کنند🙏

%v

@kiwi\_proxy
`

				caption = fmt.Sprintf(caption, proxyLinksText)
				_, err = b.SendPhoto(ctx, &bot.SendPhotoParams{
					ChatID:      conf.PublishChannelId,
					Photo:       photoFile,
					Caption:     caption,
					ParseMode:   models.ParseModeMarkdown,
					ReplyMarkup: proxyBtn,
				})
				if err != nil {
					fmt.Println(err)
				}
			}
		case conf.ContentChannelId:
			if len(update.ChannelPost.Photo) > 0 {
				post := &entities.Post{
					FileId: update.ChannelPost.Photo[0].FileID,
				}

				if err := postRepo.Create(ctx, post); err != nil {
					fmt.Println("create post failed:", err)
					return
				}
			}
		}
	}
}

func addLang(next bot.HandlerFunc) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		next(context.WithValue(ctx, "lang", "fa"), b, update)
	}
}
