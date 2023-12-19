package handlers

import (
	"context"
	"fmt"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/h-varmazyar/kiwi/applications/film/internal/repositories"
	"github.com/h-varmazyar/kiwi/applications/film/pkg/entities"
	"github.com/h-varmazyar/kiwi/applications/film/pkg/helpers"
	"strconv"
	"strings"
)

func (h *Handler) addMovie(ctx context.Context, b *bot.Bot, update *models.Message, data []byte) {
	if err := h.userStateRepo.SetState(ctx, update.Chat.ID, repositories.StateAddMovie); err != nil {
		SendError(ctx, b, &ErrParams{
			ChatId:   update.Chat.ID,
			Err:      err,
			Metadata: update,
		})
		return
	}

	movie := &entities.Movie{}
	if err := h.addContentRepo.SetMovie(ctx, update.Chat.ID, movie); err != nil {
		SendError(ctx, b, &ErrParams{
			ChatId: update.Chat.ID,
			Err:    err,
		})
		return
	}
	_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Chat.ID,
		Text:   MsgAddMovieInfo,
	})
}

func (h *Handler) stateAddMovie(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}

	chatId := update.Message.Chat.ID
	data := strings.Split(strings.TrimSpace(update.Message.Text), "\n")
	if len(data) != 10 {
		SendError(ctx, b, &ErrParams{
			ChatId:   chatId,
			Err:      ErrInvalidMovieData,
			Metadata: update,
		})
		return
	}

	movie, err := h.addContentRepo.GetMovie(ctx, chatId)
	if err != nil {
		SendError(ctx, b, &ErrParams{
			ChatId:   chatId,
			Err:      err,
			Metadata: update,
		})
		return
	}

	movie.Title = strings.TrimSpace(data[0])
	movie.FaName = strings.TrimSpace(data[1])
	movie.EnName = strings.TrimSpace(data[2])
	movie.Presentation = strings.TrimSpace(data[3])

	year, err := strconv.Atoi(strings.TrimSpace(data[4]))
	if err != nil {
		SendError(ctx, b, &ErrParams{
			ChatId:   chatId,
			Err:      ErrInvalidYear,
			Metadata: update,
		})
		return
	}
	movie.Year = year

	movie.IMDBLink = strings.TrimSpace(data[5])

	langs := make([]string, 0)
	for _, lang := range strings.Split(strings.TrimSpace(data[6]), " ") {
		langs = append(langs, strings.TrimSpace(lang))
	}
	movie.Languages = langs

	tags := make([]string, 0)
	for _, tag := range strings.Split(strings.TrimSpace(data[7]), " ") {
		tags = append(tags, strings.TrimSpace(tag))
	}
	movie.Tags = tags

	if strings.TrimSpace(data[8]) == "+" {
		movie.HardSubtitle = true
	} else {
		movie.SubtitleLink = strings.TrimSpace(data[9])
	}

	if err = h.addContentRepo.SetMovie(ctx, chatId, movie); err != nil {
		SendError(ctx, b, &ErrParams{
			ChatId:   chatId,
			Err:      err,
			Metadata: update,
		})
		return
	}

	if err = h.userStateRepo.SetState(ctx, chatId, repositories.StateAddMovieBanner); err != nil {
		SendError(ctx, b, &ErrParams{
			ChatId:   chatId,
			Err:      err,
			Metadata: "cannot move to state " + repositories.StateAddMovieBanner,
		})
		return
	}

	_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      chatId,
		Text:        MsgAddMovieBanner,
		ReplyMarkup: h.keyboardCancel(b),
	})
}

func (h *Handler) stateAddMovieBanner(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil {

		return
	}

	chatId := update.Message.Chat.ID
	if update.Message.Photo == nil || len(update.Message.Photo) == 0 {
		SendError(ctx, b, &ErrParams{
			ChatId:   chatId,
			Err:      ErrNoBannerAdded,
			Metadata: update,
		})
		return
	}

	movie, err := h.addContentRepo.GetMovie(ctx, chatId)
	if err != nil {
		SendError(ctx, b, &ErrParams{
			ChatId:   chatId,
			Err:      ErrNoMovieCached,
			Metadata: update,
		})
		return
	}

	movie.Banner = &entities.Media{
		TelegramFileId: update.Message.Photo[0].FileID,
		Type:           entities.MediaTypePicture,
	}

	if err = h.addContentRepo.SetMovie(ctx, chatId, movie); err != nil {
		SendError(ctx, b, &ErrParams{
			ChatId:   chatId,
			Err:      err,
			Metadata: movie,
		})
		return
	}

	if err = h.userStateRepo.DeleteState(ctx, chatId); err != nil {
		SendError(ctx, b, &ErrParams{
			ChatId: chatId,
			Err:    err,
		})
		return
	}

	photo := &models.InputFileString{
		Data: movie.Banner.TelegramFileId,
	}

	_, err = b.SendPhoto(ctx, &bot.SendPhotoParams{
		ChatID:      chatId,
		Caption:     prepareMovieCaptionAllQuality(movie),
		Photo:       photo,
		ParseMode:   models.ParseModeMarkdown,
		ReplyMarkup: h.keyboardMovieSubmit(ctx, b, movie.ID),
	})

	if err != nil {
		SendError(ctx, b, &ErrParams{
			ChatId: chatId,
			Err:    err,
		})
		return
	}
}

func (h *Handler) submitMovieAddition(ctx context.Context, b *bot.Bot, update *models.Message, data []byte) {
	chatId := update.Chat.ID
	movie, err := h.addContentRepo.GetMovie(ctx, chatId)
	if err != nil {
		SendError(ctx, b, &ErrParams{
			ChatId:   chatId,
			Err:      ErrNoMovieCached,
			Metadata: update,
		})
		return
	}

	if err = h.moviesRepo.Create(ctx, movie); err != nil {
		SendError(ctx, b, &ErrParams{
			ChatId:   chatId,
			Metadata: update,
		})
		return
	}

	if err = h.addContentRepo.SetMovie(ctx, chatId, movie); err != nil {
		SendError(ctx, b, &ErrParams{
			ChatId:   chatId,
			Err:      err,
			Metadata: movie,
		})
		return
	}

	photo := &models.InputFileString{
		Data: movie.Banner.TelegramFileId,
	}
	_, _ = b.SendPhoto(ctx, &bot.SendPhotoParams{
		ChatID:      update.Chat.ID,
		Caption:     prepareMovieCaptionAllQuality(movie),
		Photo:       photo,
		ParseMode:   models.ParseModeMarkdown,
		ReplyMarkup: h.keyboardAddMovieVideo(ctx, b, movie.ID),
	})
}

func (h *Handler) cancelMovieAddition(ctx context.Context, b *bot.Bot, update *models.Message, _ []byte) {
	if err := h.addContentRepo.DeleteMovie(ctx, update.Chat.ID); err != nil {
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

func (h *Handler) addMovieVideo(ctx context.Context, b *bot.Bot, update *models.Message, data []byte) {
	if err := h.userStateRepo.SetState(ctx, update.Chat.ID, repositories.StateAddMedia); err != nil {
		SendError(ctx, b, &ErrParams{
			ChatId:   update.Chat.ID,
			Err:      err,
			Metadata: update,
		})
		return
	}

	movieId, err := strconv.Atoi(string(data))
	if err != nil {
		SendError(ctx, b, &ErrParams{
			ChatId:   update.Chat.ID,
			Err:      err,
			Metadata: string(data),
		})
		return
	}

	media := &entities.Media{
		Type:      entities.MediaTypeVideo,
		OwnerID:   uint(movieId),
		OwnerType: "movie",
	}

	if err = h.addContentRepo.SetMedia(ctx, update.Chat.ID, media); err != nil {
		SendError(ctx, b, &ErrParams{
			ChatId:   update.Chat.ID,
			Err:      err,
			Metadata: string(data),
		})
		return
	}

	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.Chat.ID,
		Text:        MsgAddQuality,
		ReplyMarkup: h.keyboardCancel(b),
	})
}

func prepareMovieCaptionAllQuality(movie *entities.Movie) string {
	text := helpers.EscapeText(fmt.Sprintf("📹 %v\n", movie.Title))
	text += helpers.EscapeText(fmt.Sprintf("🇮🇷 %v\n", movie.FaName))
	text += helpers.EscapeText(fmt.Sprintf("🏴󠁧󠁢󠁥󠁮󠁧󠁿 %v\n", movie.EnName))
	text += helpers.EscapeText(fmt.Sprintf("📝 %v\n", movie.Presentation))
	text += helpers.EscapeText(fmt.Sprintf("📽 %v\n", movie.Year))

	languages := ""
	for _, language := range movie.Languages {
		languages = fmt.Sprintf("%v - %v", languages, language)
	}
	languages = strings.Trim(strings.TrimSpace(languages), "-")
	text += helpers.EscapeText(fmt.Sprintf("🗣 %v\n", languages))

	if movie.HardSubtitle {
		text += helpers.EscapeText(fmt.Sprintf("💬 %v\n", hardSubtitleLabel))
	} else if movie.SubtitleLink != "" {
		text += fmt.Sprintf("💬 [دانلود زیرنویس](%v)\n", hardSubtitleLabel)
	}

	if movie.IMDB != 0 {
		text += helpers.EscapeText(fmt.Sprintf("💯 %v\n", movie.IMDB))
	}

	hashtags := ""
	for _, hashtag := range movie.Tags {
		if !strings.HasPrefix(hashtag, "#") {
			hashtag = fmt.Sprintf("#%v", hashtag)
		}
		hashtags = fmt.Sprintf("%v %v", hashtags, hashtag)
	}

	if hashtags != "" {
		text += helpers.EscapeText(fmt.Sprintf("\n#️⃣ %v\n", hashtags))
	}

	downloads := ""
	for _, video := range movie.Videos {
		download := fmt.Sprintf("⬇️ کیفیت %v: ", video.Quality)
		if video.DownloadUrl != "" {
			download += fmt.Sprintf("-\\ [دانلود](%v)", video.DownloadUrl)
		}
		botLink := fmt.Sprintf("https://t.me/Kiwifilm_bot?start=%v", video.ID)
		download += fmt.Sprintf(" [مشاهده](%v)\n", botLink)
		downloads += download
	}

	if downloads != "" {
		text += fmt.Sprintf("\n%v", downloads)
	}

	return fmt.Sprintf("%v\n\n@kiwifilm", text)
}

func prepareMovieQualityCaption(movie *entities.Movie, media *entities.Media) string {
	text := helpers.EscapeText(fmt.Sprintf("📹 %v\n", movie.Title))
	text = helpers.EscapeText(fmt.Sprintf("%v🇮🇷 %v\n", text, movie.FaName))
	text = helpers.EscapeText(fmt.Sprintf("%v🏴󠁧󠁢󠁥󠁮󠁧󠁿 %v\n", text, movie.EnName))

	if movie.HardSubtitle {
		text = helpers.EscapeText(fmt.Sprintf("%v💬 %v\n", text, hardSubtitleLabel))
	} else if movie.SubtitleLink != "" {
		text = fmt.Sprintf("%v💬 [دانلود زیرنویس](%v)\n", text, hardSubtitleLabel)
	}

	text = helpers.EscapeText(fmt.Sprintf("%v🎞 %v\n", text, media.Quality))

	return fmt.Sprintf("%v\n\n@kiwifilm", text)
}

func (h *Handler) selectMovie(ctx context.Context, b *bot.Bot, update *models.Message, data []byte) {
	chatId := update.Chat.ID
	movieId, err := strconv.Atoi(string(data))
	if err != nil {
		SendError(ctx, b, &ErrParams{
			ChatId:   chatId,
			Err:      err,
			Metadata: update,
		})
		return
	}

	movie, err := h.moviesRepo.Return(ctx, uint(movieId))
	if err != nil {
		SendError(ctx, b, &ErrParams{
			ChatId:   chatId,
			Err:      err,
			Metadata: update,
		})
		return
	}

	genres := ""
	for _, genre := range movie.Genres {
		genres = fmt.Sprintf("%v | %v", genres, genre.FaName)
	}

	genres = strings.Trim(strings.TrimSpace(genres), "|")

	text := fmt.Sprintf(MsgSeriesInfo, movie.Title, movie.FaName, movie.EnName, movie.Presentation, movie.Year, movie.IMDB, genres)

	if movie.Banner == nil {
		SendError(ctx, b, &ErrParams{
			ChatId: chatId,
		})
		return
	}

	photo := &models.InputFileString{
		Data: movie.Banner.TelegramFileId,
	}
	_, _ = b.SendPhoto(ctx, &bot.SendPhotoParams{
		ChatID:      update.Chat.ID,
		Photo:       photo,
		Caption:     text,
		ReplyMarkup: h.keyboardMovieInfo(ctx, b, h.isAdmin(update.Chat.ID), movie.ID),
	})
}

func (h *Handler) sendMovieToPublicChannel(ctx context.Context, b *bot.Bot, update *models.Message, data []byte) {
	chatId := update.Chat.ID

	if h.configs.PublicChannelId < 0 {
		movieId, err := strconv.Atoi(string(data))
		if err != nil {
			SendError(ctx, b, &ErrParams{
				ChatId:   chatId,
				Err:      err,
				Metadata: update,
			})
			return
		}

		movie, err := h.moviesRepo.Return(ctx, uint(movieId))
		if err != nil {
			SendError(ctx, b, &ErrParams{
				ChatId:   chatId,
				Err:      err,
				Metadata: update,
			})
			return
		}

		photo := &models.InputFileString{
			Data: movie.Banner.TelegramFileId,
		}
		_, err = b.SendPhoto(ctx, &bot.SendPhotoParams{
			ChatID:    h.configs.PublicChannelId,
			Caption:   prepareMovieCaptionAllQuality(movie),
			Photo:     photo,
			ParseMode: models.ParseModeMarkdown,
		})
		if err != nil {
			SendError(ctx, b, &ErrParams{
				ChatId: chatId,
				Err:    err,
			})
		}
	} else {
		_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Chat.ID,
			Text:   "کانال عمومی یافت نشد",
		})
	}
}

func (h *Handler) showMovie(ctx context.Context, b *bot.Bot, update *models.Message, data []byte) {
	chatId := update.Chat.ID

	movieId, err := strconv.Atoi(string(data))
	if err != nil {
		SendError(ctx, b, &ErrParams{
			ChatId:   chatId,
			Err:      err,
			Metadata: update,
		})
		return
	}

	movie, err := h.moviesRepo.Return(ctx, uint(movieId))
	if err != nil {
		SendError(ctx, b, &ErrParams{
			ChatId:   chatId,
			Err:      err,
			Metadata: update,
		})
		return
	}

	for _, video := range movie.Videos {
		videoFile := &models.InputFileString{
			Data: video.TelegramFileId,
		}
		_, err = b.SendVideo(ctx, &bot.SendVideoParams{
			ChatID:    chatId,
			Caption:   prepareMovieQualityCaption(movie, video),
			Video:     videoFile,
			ParseMode: models.ParseModeMarkdown,
		})
		if err != nil {
			SendError(ctx, b, &ErrParams{
				ChatId: chatId,
				Err:    err,
			})
		}
	}
}
