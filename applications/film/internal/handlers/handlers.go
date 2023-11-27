package handlers

import (
	"context"
	"fmt"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	repositories "github.com/h-varmazyar/kiwi/applications/film/internal/repositories"
	"github.com/h-varmazyar/kiwi/applications/film/pkg/db/PostgreSQL"
	redisPkg "github.com/h-varmazyar/kiwi/applications/film/pkg/db/redis"
	log "github.com/sirupsen/logrus"
	"regexp"
)

type Configs struct {
	AdminChatId     int64             `yaml:"adminChatId"`
	RedisDB         int               `yaml:"redisDB"`
	PublicChannelId int64             `yaml:"publicChannelId"`
	RedisConfigs    *redisPkg.Configs `yaml:"-"`
}

type Handler struct {
	configs             Configs
	log                 *log.Logger
	userStateRepository *repositories.UserStateRepository
	addContentRepo      *repositories.AddContentRepository
	mediaRepo           *repositories.MediaRepository
	seriesRepo          *repositories.SeriesRepository
	episodeRepo         *repositories.EpisodeRepository
	seasonRepo          *repositories.SeasonRepository
	userStateRepo       *repositories.UserStateRepository
}

func NewHandler(ctx context.Context, log *log.Logger, configs Configs, db *db.DB) (*Handler, error) {
	configs.RedisConfigs.DB = configs.RedisDB
	redisClient := redisPkg.NewClient(configs.RedisConfigs)

	userStateRepo := repositories.NewUserState(redisClient)
	addContentRepo := repositories.NewAddContentRepository(redisClient)

	mediaRepo, err := repositories.NewMediaRepository(ctx, log, db)
	if err != nil {
		return nil, err
	}
	episodeRepo, err := repositories.NewEpisodeRepository(ctx, log, db)
	if err != nil {
		return nil, err
	}
	seasonRepo, err := repositories.NewSeasonRepository(ctx, log, db)
	if err != nil {
		return nil, err
	}
	seriesRepo, err := repositories.NewSeriesRepository(ctx, log, db)
	if err != nil {
		return nil, err
	}

	return &Handler{
		configs:             configs,
		log:                 log,
		userStateRepository: userStateRepo,
		addContentRepo:      addContentRepo,
		mediaRepo:           mediaRepo,
		seriesRepo:          seriesRepo,
		episodeRepo:         episodeRepo,
		seasonRepo:          seasonRepo,
		userStateRepo:       userStateRepo,
	}, nil
}

func (h *Handler) RegisterHandlers(_ context.Context, b *bot.Bot) {
	startRegex := regexp.MustCompile(`^/start [a-zA-Z0-9]*`)
	b.RegisterHandlerRegexp(bot.HandlerTypeMessageText, startRegex, h.startCmd)
	b.RegisterHandler(bot.HandlerTypeMessageText, "", bot.MatchTypeContains, h.handleMessage)
}

func (h *Handler) handleMessage(ctx context.Context, b *bot.Bot, update *models.Update) {
	userId := update.Message.Chat.ID

	if userId != h.configs.AdminChatId {
		return
	}

	state, err := h.userStateRepository.GetState(ctx, userId)
	if err != nil {
		fmt.Println(err)
		_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "unhandled state",
		})
		return
	}
	switch state {
	case repositories.StateAddSeries:
	case repositories.StateAddEpisode:
		h.stateAddEpisode(ctx, b, update)
	case repositories.StateAddEpisodeBanner:
		h.stateAddEpisodeBanner(ctx, b, update)
	case repositories.StateAddEpisodeVideo:
		h.stateAddMedia(ctx, b, update)
	case repositories.StateAddBanner:
		//h.addHandler.AddBannerState(ctx, b, update)
	default:
		h.search(ctx, b, update)
	}
}

func (h *Handler) search(ctx context.Context, b *bot.Bot, update *models.Update) {
	series, err := h.seriesRepo.Search(ctx, update.Message.Text)
	if err != nil {
		SendError(ctx, b, &ErrParams{
			ChatId:   update.Message.Chat.ID,
			Err:      err,
			Metadata: update,
		})
		return
	}

	_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.Message.Chat.ID,
		Text:        MsgSelectSeriesBetweenSearch,
		ReplyMarkup: h.keyboardSearchedSeries(ctx, b, series),
	})
}

func (h *Handler) cancelAddition(ctx context.Context, b *bot.Bot, update *models.Message, data []byte) {
	if err := h.userStateRepo.DeleteState(ctx, update.Chat.ID); err != nil {
		SendError(ctx, b, &ErrParams{
			ChatId:   update.Chat.ID,
			Err:      err,
			Metadata: update,
		})
		return
	}

	SendSuccess(ctx, b, &SuccessParams{
		ChatId:   update.Chat.ID,
		IsSilent: false,
	})
}
