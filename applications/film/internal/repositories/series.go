package repositories

import (
	"context"
	e "errors"
	db2 "github.com/h-varmazyar/kiwi/applications/film/pkg/db/PostgreSQL"
	"github.com/h-varmazyar/kiwi/applications/film/pkg/entities"
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
			err = tx.AutoMigrate(new(entities.Movie))
			if err != nil {
				return err
			}
			err = tx.AutoMigrate(new(entities.Series))
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
			if !tx.Migrator().HasColumn(&entities.Series{}, "PrivateChannelId") {
				err = tx.Migrator().AddColumn(&entities.Series{}, "PrivateChannelId")
				newMigrations = append(newMigrations, &db2.Migration{
					TableName:   seriesTableName,
					Tag:         "v1.0.1",
					Description: "add private channel id",
				})
				if err != nil {
					return err
				}
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

func (r *SeriesRepository) Create(_ context.Context, series *entities.Series) error {
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

func (r *SeriesRepository) Search(_ context.Context, searchQuery string) ([]*entities.Series, error) {
	series := make([]*entities.Series, 0)
	if err := r.PostgresDB.Model(new(entities.Movie)).
		Where("title like %?%", searchQuery).
		Or("fa_name like %?%", searchQuery).
		Or("en_name like %?%", searchQuery).
		Find(&series).Limit(5).Error; err != nil {
		return nil, err
	}
	return series, nil
}

func (r *SeriesRepository) Return(_ context.Context, id uint) (*entities.Series, error) {
	series := new(entities.Series)
	if err := r.PostgresDB.Model(new(entities.Series)).Preload("Banner").Preload("Genres").Where("id = ?", id).First(series).Error; err != nil {
		if e.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNoSeriesFound
		}
		return nil, err
	}
	return series, nil
}
