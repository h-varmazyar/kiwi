package repositories

import (
	"context"
	e "errors"
	db2 "github.com/h-varmazyar/kiwi/applications/film/pkg/db/PostgreSQL"
	entities2 "github.com/h-varmazyar/kiwi/applications/film/pkg/entities"
	"github.com/h-varmazyar/kiwi/pkg/errors"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var (
	ErrNoSeriesFound = errors.NewWithCode("no_series_found", 2001)
)

var seriesTableName = "series"

type SeriesRepository struct {
	log *log.Logger
	*db2.DB
}

func NewSeriesRepository(ctx context.Context, log *log.Logger, db *db2.DB) (*SeriesRepository, error) {
	repo := &SeriesRepository{
		log: log,
		DB:  db,
	}

	if err := repo.migration(ctx, db); err != nil {
		return nil, err
	}

	return repo, nil
}

func (r *SeriesRepository) migration(_ context.Context, dbInstance *db2.DB) error {
	var err error
	migrations := make(map[string]struct{})
	tags := make([]string, 0)
	err = dbInstance.PostgresDB.Table(db2.MigrationTable).Where("table_name = ?", seriesTableName).Select("tag").Find(&tags).Error
	if err != nil {
		return err
	}

	for _, tag := range tags {
		migrations[tag] = struct{}{}
	}

	newMigrations := make([]*db2.Migration, 0)
	err = dbInstance.PostgresDB.Transaction(func(tx *gorm.DB) error {
		if _, ok := migrations["v1.0.0"]; !ok {
			err = tx.AutoMigrate(new(entities2.Movie))
			if err != nil {
				return err
			}
			err = tx.AutoMigrate(new(entities2.Series))
			if err != nil {
				return err
			}
			newMigrations = append(newMigrations, &db2.Migration{
				TableName:   seriesTableName,
				Tag:         "v1.0.0",
				Description: "create series table",
			})
		}
		if _, ok := migrations["v1.0.1"]; !ok {
			err = tx.Migrator().AddColumn(&entities2.Series{}, "PrivateChannelId")
			newMigrations = append(newMigrations, &db2.Migration{
				TableName:   seriesTableName,
				Tag:         "v1.0.1",
				Description: "add private channel id",
			})
			if err != nil {
				return err
			}
		}
		err = tx.Model(new(db2.Migration)).CreateInBatches(&newMigrations, 100).Error
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (r *SeriesRepository) Create(_ context.Context, series *entities2.Series) error {
	if err := r.PostgresDB.Transaction(func(tx *gorm.DB) error {
		//if err := tx.Model(new(entities.Media)).Save(series.Banner).Error; err != nil {
		//	return err
		//}
		//series.BannerId = series.Banner.ID
		if err := tx.Save(series).Error; err != nil {
			return err
		}
		//
		//for _, genre := range series.Genres {
		//
		//}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func (r *SeriesRepository) Search(_ context.Context, searchQuery string) ([]*entities2.Series, error) {
	series := make([]*entities2.Series, 0)
	if err := r.PostgresDB.Model(new(entities2.Series)).Find(&series).Limit(5).Error; err != nil {
		return nil, err
	}
	if len(series) == 0 {
		return nil, ErrNoSeriesFound
	}

	return series, nil
}

func (r *SeriesRepository) Return(_ context.Context, id uint) (*entities2.Series, error) {
	series := new(entities2.Series)
	if err := r.PostgresDB.Model(new(entities2.Series)).Preload("Banner").Preload("Genres").Where("id = ?", id).First(series).Error; err != nil {
		if e.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNoSeriesFound
		}
		return nil, err
	}
	return series, nil
}
