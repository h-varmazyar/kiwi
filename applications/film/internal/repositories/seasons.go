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
	ErrNoSeasonFound = errors.NewWithCode("no_season_found", 2101)
)

var seasonTableName = "seasons"

type SeasonRepository struct {
	log *log.Logger
	*db2.DB
}

func NewSeasonRepository(ctx context.Context, log *log.Logger, db *db2.DB) (*SeasonRepository, error) {
	repo := &SeasonRepository{
		log: log,
		DB:  db,
	}

	if err := repo.migration(ctx, db); err != nil {
		return nil, err
	}

	return repo, nil
}

func (r *SeasonRepository) migration(_ context.Context, dbInstance *db2.DB) error {
	var err error
	migrations := make(map[string]struct{})
	tags := make([]string, 0)
	err = dbInstance.PostgresDB.Table(db2.MigrationTable).Where("table_name = ?", seasonTableName).Select("tag").Find(&tags).Error
	if err != nil {
		return err
	}

	for _, tag := range tags {
		migrations[tag] = struct{}{}
	}

	newMigrations := make([]*db2.Migration, 0)
	err = dbInstance.PostgresDB.Transaction(func(tx *gorm.DB) error {
		if _, ok := migrations["v1.0.0"]; !ok {
			err = tx.AutoMigrate(new(entities.Season))
			if err != nil {
				return err
			}
			newMigrations = append(newMigrations, &db2.Migration{
				TableName:   seasonTableName,
				Tag:         "v1.0.0",
				Description: "create seasons table",
			})
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

func (r *SeasonRepository) Create(_ context.Context, season *entities.Season) error {
	if err := r.PostgresDB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(season).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func (r *SeasonRepository) Return(_ context.Context, id uint) (*entities.Season, error) {
	season := new(entities.Season)
	if err := r.PostgresDB.Model(new(entities.Season)).Preload("Banner").Where("id = ?", id).First(season).Error; err != nil {
		if e.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNoSeasonFound
		}
		return nil, err
	}
	return season, nil
}

func (r *SeasonRepository) SeriesSeasons(_ context.Context, seriesId uint) ([]*entities.Season, error) {
	seasons := make([]*entities.Season, 0)
	if err := r.PostgresDB.Model(new(entities.Season)).Where("series_id = ?", seriesId).Find(&seasons).Error; err != nil {
		return nil, err
	}
	if len(seasons) == 0 {
		return nil, ErrNoSeasonFound
	}
	return seasons, nil
}
