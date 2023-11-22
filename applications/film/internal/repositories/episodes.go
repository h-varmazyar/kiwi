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
	ErrNoEpisodeFound = errors.NewWithCode("no_episode_found", 2201)
)

var (
	episodeTable = "episodes"
)

type EpisodeRepository struct {
	log *log.Logger
	*db2.DB
}

func NewEpisodeRepository(ctx context.Context, log *log.Logger, db *db2.DB) (*EpisodeRepository, error) {
	repo := &EpisodeRepository{
		log: log,
		DB:  db,
	}

	if err := repo.migration(ctx, db); err != nil {
		return nil, err
	}

	return repo, nil
}

func (r *EpisodeRepository) migration(_ context.Context, dbInstance *db2.DB) error {
	var err error
	migrations := make(map[string]struct{})
	tags := make([]string, 0)
	err = dbInstance.PostgresDB.Table(db2.MigrationTable).Where("table_name = ?", episodeTable).Select("tag").Find(&tags).Error
	if err != nil {
		return err
	}

	for _, tag := range tags {
		migrations[tag] = struct{}{}
	}

	newMigrations := make([]*db2.Migration, 0)
	err = dbInstance.PostgresDB.Transaction(func(tx *gorm.DB) error {
		if _, ok := migrations["v1.0.0"]; !ok {
			err = tx.AutoMigrate(new(entities.Episode))
			if err != nil {
				return err
			}
			newMigrations = append(newMigrations, &db2.Migration{
				TableName:   episodeTable,
				Tag:         "v1.0.0",
				Description: "create episodes table",
			})
		}
		if _, ok := migrations["v1.0.1"]; !ok {
			err = tx.Migrator().AddColumn(&entities.Episode{}, "VisitCount")
			if err != nil {
				return err
			}
			newMigrations = append(newMigrations, &db2.Migration{
				TableName:   episodeTable,
				Tag:         "v1.0.1",
				Description: "add visit_count column",
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

func (r *EpisodeRepository) Create(_ context.Context, episode *entities.Episode) error {
	if err := r.PostgresDB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(episode).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func (r *EpisodeRepository) Return(_ context.Context, id uint) (*entities.Episode, error) {
	episode := new(entities.Episode)
	if err := r.PostgresDB.Model(new(entities.Episode)).Preload("Banner").Preload("Videos").Where("id = ?", id).First(episode).Error; err != nil {
		if e.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNoEpisodeFound
		}
		return nil, err
	}
	return episode, nil
}

func (r *EpisodeRepository) Visit(_ context.Context, id uint) error {
	if err := r.PostgresDB.Model(new(entities.Episode)).Where("id = ?", id).Update("visit_count", gorm.Expr("visit_count + 1")).Error; err != nil {
		if e.Is(err, gorm.ErrRecordNotFound) {
			return ErrNoEpisodeFound
		}
		return err
	}
	return nil
}

func (r *EpisodeRepository) List(_ context.Context, seasonId uint) ([]*entities.Episode, error) {
	episodes := make([]*entities.Episode, 0)
	if err := r.PostgresDB.Model(new(entities.Episode)).Where("season_id = ?", seasonId).Find(&episodes).Error; err != nil {
		return nil, err
	}
	if len(episodes) == 0 {
		return nil, ErrNoEpisodeFound
	}
	return episodes, nil
}
