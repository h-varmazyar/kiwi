package handlers

import (
	"github.com/h-varmazyar/kiwi/pkg/errors"
)

var (
	ErrNoSeriesFound      = errors.NewWithCode("no_series_found", 1001)
	ErrInvalidSeriesData  = errors.NewWithCode("invalid_series_data", 1002)
	ErrInvalidEpisodeData = errors.NewWithCode("invalid_episode_data", 1003)
	ErrNoBannerAdded      = errors.NewWithCode("no_banner_added", 1004)
	ErrNoEpisodeCached    = errors.NewWithCode("no_episode_cached", 1005)
)
