package repositories

import (
	"context"
	db2 "github.com/h-varmazyar/kiwi/applications/film/pkg/db/PostgreSQL"
	entities2 "github.com/h-varmazyar/kiwi/applications/film/pkg/entities"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var genreTableName = "genres"

type GenreRepository struct {
	log *log.Logger
	*db2.DB
}

func NewGenreRepository(ctx context.Context, log *log.Logger, db *db2.DB) (*GenreRepository, error) {
	repo := &GenreRepository{
		log: log,
		DB:  db,
	}

	return repo, nil
}

func (r *GenreRepository) migration(_ context.Context, dbInstance *db2.DB) error {
	var err error
	migrations := make(map[string]struct{})
	tags := make([]string, 0)
	err = dbInstance.PostgresDB.Table(db2.MigrationTable).Where("table_name = ?", genreTableName).Select("tag").Find(&tags).Error
	if err != nil {
		return err
	}

	for _, tag := range tags {
		migrations[tag] = struct{}{}
	}

	newMigrations := make([]*db2.Migration, 0)
	err = dbInstance.PostgresDB.Transaction(func(tx *gorm.DB) error {
		if _, ok := migrations["v1.0.0"]; !ok {
			err = tx.AutoMigrate(new(entities2.Series))
			if err != nil {
				return err
			}
			newMigrations = append(newMigrations, &db2.Migration{
				TableName:   genreTableName,
				Tag:         "v1.0.0",
				Description: "create genres table",
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

func (r *GenreRepository) List(ctx context.Context) ([]*entities2.Genre, error) {
	genres := make([]*entities2.Genre, 0)
	if err := r.PostgresDB.Model(new(entities2.Genre)).Find(&genres).Error; err != nil {
		return nil, err
	}
	return genres, nil
}

func (r *GenreRepository) Return(ctx context.Context, id uint) (*entities2.Genre, error) {
	genre := new(entities2.Genre)
	if err := r.PostgresDB.Model(new(entities2.Genre)).Where("id = ?", id).Find(genre).Error; err != nil {
		return nil, err
	}
	return genre, nil
}
