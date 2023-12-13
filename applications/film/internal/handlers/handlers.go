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
	"strings"
)

var (
	startRegex = regexp.MustCompile(`^/start[ a-zA-Z0-9]*`)
	addRegex   = regexp.MustCompile(`^/add`)
)

type Configs struct {
	RedisDB         int                `yaml:"redisDB"`
	PublicChannelId int64              `yaml:"publicChannelId"`
	RedisConfigs    *redisPkg.Configs  `yaml:"-"`
	Admins          map[int64]struct{} `yaml:"-"`
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
	moviesRepo          *repositories.MoviesRepository
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
	moviesRepo, err := repositories.NewMoviesRepository(ctx, log, db)
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
		moviesRepo:          moviesRepo,
		userStateRepo:       userStateRepo,
	}, nil
}

func (h *Handler) RegisterHandlers(_ context.Context, b *bot.Bot) {
	b.RegisterHandlerRegexp(bot.HandlerTypeMessageText, startRegex, h.startCmd)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/add", bot.MatchTypeExact, h.addCmd)
	b.RegisterHandlerMatchFunc(h.isInState, h.handleStates)
	b.RegisterHandlerMatchFunc(h.isSearch, h.search)
}

func (h *Handler) isInState(update *models.Update) bool {
	if update.Message == nil {
		return false
	}

	userId := update.Message.Chat.ID
	if !h.isAdmin(userId) {
		return false
	}

	state, err := h.userStateRepository.GetState(context.Background(), userId)
	if err != nil {
		h.log.WithError(err).Error("failed to fetch user state")
		return false
	}
	if state == "" {
		return false
	}
	return true
}

func (h *Handler) isSearch(update *models.Update) bool {
	if update.Message == nil {
		return false
	}

	if strings.HasPrefix(update.Message.Text, "/") {
		return false
	}
	return true
}

func (h *Handler) handleStates(ctx context.Context, b *bot.Bot, update *models.Update) {
	userId := update.Message.Chat.ID

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
	case repositories.StateAddMovie:
		h.stateAddMovie(ctx, b, update)
	case repositories.StateAddMovieBanner:
		h.stateAddMovieBanner(ctx, b, update)
	case repositories.StateAddEpisode:
		h.stateAddEpisode(ctx, b, update)
	case repositories.StateAddEpisodeBanner:
		h.stateAddEpisodeBanner(ctx, b, update)
	case repositories.StateAddMedia:
		h.stateAddMedia(ctx, b, update)
	default:
		h.search(ctx, b, update)
	}
}

func (h *Handler) search(ctx context.Context, b *bot.Bot, update *models.Update) {
	fmt.Println(update.Message.Text)
	movies, _ := h.moviesRepo.Search(ctx, update.Message.Text)
	series, _ := h.seriesRepo.Search(ctx, update.Message.Text)

	searchResults := make([]*SearchResult, 0)
	for _, movie := range movies {
		sr := &SearchResult{
			FaName: movie.FaName,
			EnName: movie.EnName,
			Id:     movie.ID,
			Type:   SearchResultMovie,
		}
		searchResults = append(searchResults, sr)
	}
	for _, s := range series {
		sr := &SearchResult{
			FaName: s.FaName,
			EnName: s.EnName,
			Id:     s.ID,
			Type:   SearchResultSeries,
		}
		searchResults = append(searchResults, sr)
	}

	if len(searchResults) == 0 {
		_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   MsgNoSearchResultFound,
		})
	} else {
		_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:      update.Message.Chat.ID,
			Text:        MsgSelectBetweenSearch,
			ReplyMarkup: h.keyboardSearchedSeries(ctx, b, searchResults),
		})
	}
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
