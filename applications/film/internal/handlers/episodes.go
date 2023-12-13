package handlers

import (
	"context"
	"fmt"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/h-varmazyar/kiwi/applications/film/internal/repositories"
	"github.com/h-varmazyar/kiwi/applications/film/pkg/entities"
	"strconv"
	"strings"
)

const (
	hardSubtitleLabel = "زیر نویس چسبیده دارد"
)

func (h *Handler) showEpisodes(ctx context.Context, b *bot.Bot, update *models.Message, data []byte) {
	chatId := update.Chat.ID
	seasonId, err := strconv.Atoi(string(data))
	if err != nil {
		SendError(ctx, b, &ErrParams{
			ChatId:   chatId,
			Err:      err,
			Metadata: update,
		})
		return
	}

	season, err := h.seasonRepo.Return(ctx, uint(seasonId))
	if err != nil {
		SendError(ctx, b, &ErrParams{
			ChatId:   chatId,
			Err:      err,
			Metadata: update,
		})
		return
	}

	episodes, err := h.episodeRepo.List(ctx, uint(seasonId))
	if err != nil {
		SendError(ctx, b, &ErrParams{
			ChatId:   chatId,
			Err:      err,
			Metadata: update,
		})
		return
	}

	text := fmt.Sprintf(MsgSeasonInfo, season.Title, season.Presentation, season.Year, season.IMDB)

	_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.Chat.ID,
		Text:        text,
		ReplyMarkup: h.keyboardEpisodeList(ctx, b, episodes),
	})
}

func (h *Handler) showEpisode(ctx context.Context, b *bot.Bot, update *models.Message, data []byte) {
	chatId := update.Chat.ID
	episodeId, err := strconv.Atoi(string(data))
	if err != nil {
		SendError(ctx, b, &ErrParams{
			ChatId:   chatId,
			Err:      err,
			Metadata: update,
		})
		return
	}

	episode, err := h.episodeRepo.Return(ctx, uint(episodeId))
	if err != nil {
		SendError(ctx, b, &ErrParams{
			ChatId:   chatId,
			Err:      err,
			Metadata: update,
		})
		return
	}

	photo := &models.InputFileString{
		Data: episode.Banner.TelegramFileId,
	}
	_, _ = b.SendPhoto(ctx, &bot.SendPhotoParams{
		ChatID:      update.Chat.ID,
		Caption:     prepareEpisodeCaptionAllQuality(episode),
		Photo:       photo,
		ParseMode:   models.ParseModeMarkdown,
		ReplyMarkup: h.keyboardEpisodeAction(ctx, b, episode.ID),
	})
}

func (h *Handler) addEpisode(ctx context.Context, b *bot.Bot, update *models.Message, data []byte) {
	if err := h.userStateRepo.SetState(ctx, update.Chat.ID, repositories.StateAddEpisode); err != nil {
		SendError(ctx, b, &ErrParams{
			ChatId:   update.Chat.ID,
			Err:      err,
			Metadata: update,
		})
		return
	}

	seasonId, err := strconv.Atoi(string(data))
	if err != nil {
		SendError(ctx, b, &ErrParams{
			ChatId: update.Chat.ID,
			Err:    err,
		})
		return
	}

	episode := &entities.Episode{SeasonId: uint(seasonId)}
	if err := h.addContentRepo.SetEpisode(ctx, update.Chat.ID, episode); err != nil {
		SendError(ctx, b, &ErrParams{
			ChatId: update.Chat.ID,
			Err:    err,
		})
		return
	}
	_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Chat.ID,
		Text:   MsgAddEpisodeInfo,
	})
}

func (h *Handler) stateAddEpisode(ctx context.Context, b *bot.Bot, update *models.Update) {
	chatId := update.Message.Chat.ID
	data := strings.Split(strings.TrimSpace(update.Message.Text), "\n")
	if len(data) != 7 {
		SendError(ctx, b, &ErrParams{
			ChatId:   chatId,
			Err:      ErrInvalidEpisodeData,
			Metadata: update,
		})
		return
	}

	episode, err := h.addContentRepo.GetEpisode(ctx, chatId)
	if err != nil {
		SendError(ctx, b, &ErrParams{
			ChatId:   chatId,
			Err:      err,
			Metadata: update,
		})
		return
	}

	episode.Title = strings.TrimSpace(data[0])
	episode.Presentation = strings.TrimSpace(data[1])
	episode.IMDBLink = strings.TrimSpace(data[5])

	if strings.TrimSpace(data[3]) == "+" {
		episode.HardSubtitle = true
	} else {
		episode.SubtitleLink = strings.TrimSpace(data[4])
	}

	tags := make([]string, 0)
	for _, tag := range strings.Split(strings.TrimSpace(data[6]), " ") {
		tags = append(tags, tag)
	}
	episode.Tags = tags

	langs := make([]string, 0)
	for _, tag := range strings.Split(strings.TrimSpace(data[2]), " ") {
		langs = append(langs, tag)
	}
	episode.Languages = langs

	if err = h.addContentRepo.SetEpisode(ctx, chatId, episode); err != nil {
		SendError(ctx, b, &ErrParams{
			ChatId:   chatId,
			Err:      err,
			Metadata: update,
		})
		return
	}

	if err = h.userStateRepo.SetState(ctx, chatId, repositories.StateAddEpisodeBanner); err != nil {
		SendError(ctx, b, &ErrParams{
			ChatId:   chatId,
			Err:      err,
			Metadata: "cannot move to state " + repositories.StateAddEpisodeBanner,
		})
		return
	}

	_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      chatId,
		Text:        MsgAddEpisodeBanner,
		ReplyMarkup: h.keyboardCancel(b),
	})
}

func (h *Handler) stateAddEpisodeBanner(ctx context.Context, b *bot.Bot, update *models.Update) {
	chatId := update.Message.Chat.ID
	if update.Message.Photo == nil || len(update.Message.Photo) == 0 {
		SendError(ctx, b, &ErrParams{
			ChatId:   chatId,
			Err:      ErrNoBannerAdded,
			Metadata: update,
		})
		return
	}

	episode, err := h.addContentRepo.GetEpisode(ctx, chatId)
	if err != nil {
		SendError(ctx, b, &ErrParams{
			ChatId:   chatId,
			Err:      ErrNoEpisodeCached,
			Metadata: update,
		})
		return
	}

	episode.Banner = &entities.Media{
		TelegramFileId: update.Message.Photo[0].FileID,
		Type:           entities.MediaTypePicture,
	}

	if err = h.addContentRepo.SetEpisode(ctx, chatId, episode); err != nil {
		SendError(ctx, b, &ErrParams{
			ChatId:   chatId,
			Err:      err,
			Metadata: episode,
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
		Data: episode.Banner.TelegramFileId,
	}

	_, _ = b.SendPhoto(ctx, &bot.SendPhotoParams{
		ChatID:      update.Message.Chat.ID,
		Caption:     prepareEpisodeCaptionAllQuality(episode),
		Photo:       photo,
		ParseMode:   models.ParseModeMarkdown,
		ReplyMarkup: h.keyboardEpisodeSubmit(ctx, b, episode.ID),
	})
}

func (h *Handler) submitEpisodeAddition(ctx context.Context, b *bot.Bot, update *models.Message, data []byte) {
	chatId := update.Chat.ID
	episode, err := h.addContentRepo.GetEpisode(ctx, chatId)
	if err != nil {
		SendError(ctx, b, &ErrParams{
			ChatId:   chatId,
			Err:      ErrNoEpisodeCached,
			Metadata: update,
		})
		return
	}

	if err = h.episodeRepo.Create(ctx, episode); err != nil {
		SendError(ctx, b, &ErrParams{
			ChatId:   chatId,
			Metadata: update,
		})
		return
	}

	if err = h.addContentRepo.SetEpisode(ctx, chatId, episode); err != nil {
		SendError(ctx, b, &ErrParams{
			ChatId:   chatId,
			Err:      err,
			Metadata: episode,
		})
		return
	}

	photo := &models.InputFileString{
		Data: episode.Banner.TelegramFileId,
	}
	_, _ = b.SendPhoto(ctx, &bot.SendPhotoParams{
		ChatID:      update.Chat.ID,
		Caption:     prepareEpisodeCaptionAllQuality(episode),
		Photo:       photo,
		ParseMode:   models.ParseModeMarkdown,
		ReplyMarkup: h.keyboardAddEpisodeVideo(ctx, b, episode.ID),
	})
}

func (h *Handler) cancelEpisodeAddition(ctx context.Context, b *bot.Bot, update *models.Message, data []byte) {
	if err := h.addContentRepo.DeleteEpisode(ctx, update.Chat.ID); err != nil {
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

func (h *Handler) addEpisodeVideo(ctx context.Context, b *bot.Bot, update *models.Message, data []byte) {
	if err := h.userStateRepo.SetState(ctx, update.Chat.ID, repositories.StateAddMedia); err != nil {
		SendError(ctx, b, &ErrParams{
			ChatId:   update.Chat.ID,
			Err:      err,
			Metadata: update,
		})
		return
	}

	fmt.Println("set state add video")

	episodeId, err := strconv.Atoi(string(data))
	if err != nil {
		SendError(ctx, b, &ErrParams{
			ChatId:   update.Chat.ID,
			Err:      err,
			Metadata: string(data),
		})
		return
	}

	fmt.Println("episode id:", episodeId)

	media := &entities.Media{
		Type:      entities.MediaTypeVideo,
		OwnerID:   uint(episodeId),
		OwnerType: "episode",
	}

	if err = h.addContentRepo.SetMedia(ctx, update.Chat.ID, media); err != nil {
		SendError(ctx, b, &ErrParams{
			ChatId:   update.Chat.ID,
			Err:      err,
			Metadata: string(data),
		})
		return
	}

	fmt.Println("media set")

	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.Chat.ID,
		Text:        MsgAddQuality,
		ReplyMarkup: h.keyboardCancel(b),
	})
	fmt.Println(err)
}

func (h *Handler) sendEpisodeToPrivateChannel(ctx context.Context, b *bot.Bot, update *models.Message, data []byte) {
	chatId := update.Chat.ID
	episodeId, err := strconv.Atoi(string(data))
	if err != nil {
		SendError(ctx, b, &ErrParams{
			ChatId:   chatId,
			Err:      err,
			Metadata: update,
		})
		return
	}

	episode, err := h.episodeRepo.Return(ctx, uint(episodeId))
	if err != nil {
		SendError(ctx, b, &ErrParams{
			ChatId:   chatId,
			Err:      err,
			Metadata: update,
		})
		return
	}

	fmt.Println("episode found")

	season, err := h.seasonRepo.Return(ctx, episode.SeasonId)
	if err != nil {
		SendError(ctx, b, &ErrParams{
			ChatId:   chatId,
			Err:      err,
			Metadata: update,
		})
		return
	}

	fmt.Println("season found")

	series, err := h.seriesRepo.Return(ctx, season.SeriesId)
	if err != nil {
		SendError(ctx, b, &ErrParams{
			ChatId:   chatId,
			Err:      err,
			Metadata: update,
		})
		return
	}

	fmt.Println("series found")

	if series.PrivateChannelId < 0 {
		for _, v := range episode.Videos {
			video := &models.InputFileString{
				Data: v.TelegramFileId,
			}
			_, _ = b.SendVideo(ctx, &bot.SendVideoParams{
				ChatID:    series.PrivateChannelId,
				Video:     video,
				Caption:   fmt.Sprintf("%v\n\n@kiwifilm", prepareEpisodeCaptionForWatch(episode, v.Quality)),
				ParseMode: models.ParseModeMarkdown,
			})
		}
	} else {
		_, err = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Chat.ID,
			Text:   "کانال خصوصی برای این سریال وجود ندارد",
		})
	}

}

func prepareEpisodeCaptionAllQuality(episode *entities.Episode) string {
	text := fmt.Sprintf("📹 %v\n", episode.Title)
	text = fmt.Sprintf("%v📝 %v\n", text, episode.Presentation)

	languages := ""
	for _, language := range episode.Languages {
		languages = fmt.Sprintf("%v \\- %v", languages, language)
	}
	languages = strings.Trim(strings.TrimSpace(languages), "\\-")
	text = fmt.Sprintf("%v🗣 %v\n", text, languages)

	if episode.HardSubtitle {
		text = fmt.Sprintf("%v💬 %v\n", text, hardSubtitleLabel)
	} else if episode.SubtitleLink != "" {
		text = fmt.Sprintf("%v💬 [دانلود زیرنویس](%v)\n", text, hardSubtitleLabel)
	}

	if episode.IMDB != 0 {
		text = fmt.Sprintf("%v💯 %v\n", text, episode.IMDB)
	}

	hashtags := ""
	for _, hashtag := range episode.Tags {
		hashtags = fmt.Sprintf("%v \\%v", hashtags, hashtag)
	}

	if hashtags != "" {
		text = fmt.Sprintf("%v\n\\#️⃣ %v\n", text, hashtags)
	}

	downloads := ""
	for _, video := range episode.Videos {
		botLink := fmt.Sprintf("https://t.me/Kiwifilm_bot?start=%v", video.ID)
		downloads = fmt.Sprintf("%v⬇️ کیفیت %v: [دانلود](%v) \\- [مشاهده](%v)\n", downloads, video.Quality, video.DownloadUrl, botLink)
	}

	if downloads != "" {
		text = fmt.Sprintf("%v\n%v", text, downloads)
	}

	return text
}

func prepareEpisodeCaptionForWatch(episode *entities.Episode, quality entities.MediaQuality) string {
	text := fmt.Sprintf("📹 %v\n", episode.Title)
	text = fmt.Sprintf("%v📝 %v\n", text, episode.Presentation)

	languages := ""
	for _, language := range episode.Languages {
		languages = fmt.Sprintf("%v \\- %v", languages, language)
	}
	languages = strings.Trim(strings.TrimSpace(languages), "\\-")
	text = fmt.Sprintf("%v🗣 %v\n", text, languages)

	if episode.HardSubtitle {
		text = fmt.Sprintf("%v💬 %v\n", text, hardSubtitleLabel)
	} else if episode.SubtitleLink != "" {
		text = fmt.Sprintf("%v💬 [دانلود زیرنویس](%v)\n", text, hardSubtitleLabel)
	}

	if episode.IMDB != 0 {
		text = fmt.Sprintf("%v💯 %v\n", text, episode.IMDB)
	}

	text = fmt.Sprintf("%v\n⬇️ کیفیت %v\n", text, quality)

	hashtags := ""
	for _, hashtag := range episode.Tags {
		hashtags = fmt.Sprintf("%v \\%v", hashtags, hashtag)
	}

	if hashtags != "" {
		text = fmt.Sprintf("%v\n\\#️⃣ %v\n", text, hashtags)
	}

	return text
}

func (h *Handler) sendEpisodeToPublicChannel(ctx context.Context, b *bot.Bot, update *models.Message, data []byte) {
	chatId := update.Chat.ID

	if h.configs.PublicChannelId < 0 {
		episodeId, err := strconv.Atoi(string(data))
		if err != nil {
			SendError(ctx, b, &ErrParams{
				ChatId:   chatId,
				Err:      err,
				Metadata: update,
			})
			return
		}

		fmt.Println("episode parsed")

		episode, err := h.episodeRepo.Return(ctx, uint(episodeId))
		if err != nil {
			SendError(ctx, b, &ErrParams{
				ChatId:   chatId,
				Err:      err,
				Metadata: update,
			})
			return
		}

		fmt.Println("episode found")

		photo := &models.InputFileString{
			Data: episode.Banner.TelegramFileId,
		}
		_, err = b.SendPhoto(ctx, &bot.SendPhotoParams{
			ChatID:    h.configs.PublicChannelId,
			Caption:   fmt.Sprintf("%v\n\n@kiwifilm", prepareEpisodeCaptionAllQuality(episode)),
			Photo:     photo,
			ParseMode: models.ParseModeMarkdown,
		})

		fmt.Println("public err:", err)
	} else {
		_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Chat.ID,
			Text:   "کانال عمومی یافت نشد",
		})
	}
}
