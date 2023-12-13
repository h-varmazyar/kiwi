package handlers

import (
	"github.com/h-varmazyar/kiwi/pkg/errors"
)

var (
	ErrNoBannerAdded = errors.NewWithCode("no_banner_added", 1000)
	ErrInvalidYear   = errors.NewWithCode("invalid_year", 1001)

	ErrNoSeriesFound     = errors.NewWithCode("no_series_found", 1100)
	ErrInvalidSeriesData = errors.NewWithCode("invalid_series_data", 1101)

	ErrInvalidEpisodeData = errors.NewWithCode("invalid_episode_data", 1200)
	ErrNoEpisodeCached    = errors.NewWithCode("no_episode_cached", 1201)

	ErrInvalidMovieData = errors.NewWithCode("invalid_movie_data", 1300)
	ErrNoMovieCached    = errors.NewWithCode("no_movie_cached", 1301)
)
