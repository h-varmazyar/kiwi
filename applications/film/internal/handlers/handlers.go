package handlers

import (
	"context"
	"fmt"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	repositories2 "github.com/h-varmazyar/kiwi/applications/film/internal/repositories"
	"github.com/h-varmazyar/kiwi/applications/film/pkg/db/PostgreSQL"
	redisPkg "github.com/h-varmazyar/kiwi/applications/film/pkg/db/redis"
	log "github.com/sirupsen/logrus"
	"regexp"
)

type Configs struct {
	AdminChatId     int64
	RedisDB         int
	PublicChannelId int64
	RedisConfigs    *redisPkg.Configs
}

type Handler struct {
	configs             *Configs
	log                 *log.Logger
	userStateRepository *repositories2.UserStateRepository
	addContentRepo      *repositories2.AddContentRepository
	mediaRepo           *repositories2.MediaRepository
	seriesRepo          *repositories2.SeriesRepository
	episodeRepo         *repositories2.EpisodeRepository
	seasonRepo          *repositories2.SeasonRepository
	userStateRepo       *repositories2.UserStateRepository
}

func NewHandler(ctx context.Context, log *log.Logger, configs *Configs, db *db.DB) (*Handler, error) {
	configs.RedisConfigs.DB = configs.RedisDB
	redisClient := redisPkg.NewClient(configs.RedisConfigs)

	userStateRepo := repositories2.NewUserState(redisClient)
	addContentRepo := repositories2.NewAddContentRepository(redisClient)
	//genreRepo, err := repositories2.NewGenreRepository(ctx, log, db)
	//if err != nil {
	//	return nil, err
	//}
	mediaRepo, err := repositories2.NewMediaRepository(ctx, log, db)
	if err != nil {
		return nil, err
	}
	episodeRepo, err := repositories2.NewEpisodeRepository(ctx, log, db)
	if err != nil {
		return nil, err
	}
	seasonRepo, err := repositories2.NewSeasonRepository(ctx, log, db)
	if err != nil {
		return nil, err
	}
	seriesRepo, err := repositories2.NewSeriesRepository(ctx, log, db)
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
	case repositories2.StateAddSeries:
	case repositories2.StateAddEpisode:
		h.stateAddEpisode(ctx, b, update)
	case repositories2.StateAddEpisodeBanner:
		h.stateAddEpisodeBanner(ctx, b, update)
	case repositories2.StateAddEpisodeVideo:
		h.stateAddMedia(ctx, b, update)
	case repositories2.StateAddBanner:
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
