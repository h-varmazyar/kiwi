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
	ErrNoMediaFound = errors.NewWithCode("no_media_found", 2001)
)

var mediaTableName = "media"

type MediaRepository struct {
	log *log.Logger
	*db2.DB
}

func NewMediaRepository(ctx context.Context, log *log.Logger, db *db2.DB) (*MediaRepository, error) {
	repo := &MediaRepository{
		log: log,
		DB:  db,
	}

	if err := repo.migration(ctx, db); err != nil {
		return nil, err
	}

	return repo, nil
}

func (r *MediaRepository) migration(_ context.Context, dbInstance *db2.DB) error {
	var err error
	migrations := make(map[string]struct{})
	tags := make([]string, 0)
	err = dbInstance.PostgresDB.Table(db2.MigrationTable).Where("table_name = ?", mediaTableName).Select("tag").Find(&tags).Error
	if err != nil {
		return err
	}

	for _, tag := range tags {
		migrations[tag] = struct{}{}
	}

	newMigrations := make([]*db2.Migration, 0)
	err = dbInstance.PostgresDB.Transaction(func(tx *gorm.DB) error {
		if _, ok := migrations["v1.0.0"]; !ok {
			err = tx.AutoMigrate(new(entities.Media))
			if err != nil {
				return err
			}
			newMigrations = append(newMigrations, &db2.Migration{
				TableName:   mediaTableName,
				Tag:         "v1.0.0",
				Description: "create media table",
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

func (r *MediaRepository) Create(_ context.Context, media *entities.Media) error {
	if err := r.PostgresDB.Save(media).Error; err != nil {
		return err
	}
	return nil
}

func (r *MediaRepository) Return(_ context.Context, id uint) (*entities.Media, error) {
	media := new(entities.Media)
	if err := r.PostgresDB.Model(new(entities.Media)).Where("id = ?", id).First(media).Error; err != nil {
		if e.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNoSeriesFound
		}
		return nil, err
	}
	return media, nil
}
