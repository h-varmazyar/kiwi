package repositories

import (
	"context"
	e "errors"
	db "github.com/h-varmazyar/kiwi/applications/film/pkg/db/PostgreSQL"
	"github.com/h-varmazyar/kiwi/applications/film/pkg/entities"
	"github.com/h-varmazyar/kiwi/pkg/errors"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var (
	ErrNoMovieFound = errors.NewWithCode("no_movie_found", 2401)
)

var moviesTableName = "movies"

type MoviesRepository struct {
	log *log.Logger
	*db.DB
}

func NewMoviesRepository(ctx context.Context, log *log.Logger, db *db.DB) (*MoviesRepository, error) {
	repo := &MoviesRepository{
		log: log,
		DB:  db,
	}

	if err := repo.migration(ctx, db); err != nil {
		return nil, err
	}

	return repo, nil
}

func (r *MoviesRepository) migration(_ context.Context, dbInstance *db.DB) error {
	var err error
	migrations := make(map[string]struct{})
	tags := make([]string, 0)
	err = dbInstance.PostgresDB.Table(db.MigrationTable).Where("table_name = ?", moviesTableName).Select("tag").Find(&tags).Error
	if err != nil {
		return err
	}

	for _, tag := range tags {
		migrations[tag] = struct{}{}
	}

	newMigrations := make([]*db.Migration, 0)
	err = dbInstance.PostgresDB.Transaction(func(tx *gorm.DB) error {
		if _, ok := migrations["v1.0.0"]; !ok {
			err = tx.AutoMigrate(new(entities.Movie))
			if err != nil {
				return err
			}
			newMigrations = append(newMigrations, &db.Migration{
				TableName:   moviesTableName,
				Tag:         "v1.0.0",
				Description: "create movies table",
			})
		}
		err = tx.Model(new(db.Migration)).CreateInBatches(&newMigrations, 100).Error
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

func (r *MoviesRepository) Create(_ context.Context, movie *entities.Movie) error {
	r.PostgresDB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(movie).Error; err != nil {
			return err
		}
		movie.Banner.OwnerID = movie.ID
		movie.Banner.OwnerType = "movie"

		if err := tx.Save(movie.Banner).Error; err != nil {
			return err
		}
		return nil
	})
	if err := r.PostgresDB.Save(movie).Error; err != nil {
		return err
	}
	return nil
}

func (r *MoviesRepository) Search(_ context.Context, searchQuery string) ([]*entities.Movie, error) {
	movies := make([]*entities.Movie, 0)
	if err := r.PostgresDB.Model(new(entities.Movie)).
		Where("title like ?", "%"+searchQuery+"%").
		Or("fa_name like ?", "%"+searchQuery+"%").
		Or("en_name like ?", "%"+searchQuery+"%").
		Find(&movies).Limit(5).Error; err != nil {
		return nil, err
	}
	return movies, nil
}

func (r *MoviesRepository) Return(_ context.Context, id uint) (*entities.Movie, error) {
	movie := new(entities.Movie)
	if err := r.PostgresDB.Model(new(entities.Movie)).
		Preload("Banner").
		Preload("Videos", "type = ?", entities.MediaTypeVideo).
		Where("id = ?", id).First(&movie).Error; err != nil {
		return nil, err
	}
	return movie, nil
}

func (r *MoviesRepository) Visit(_ context.Context, id uint) error {
	if err := r.PostgresDB.Model(new(entities.Movie)).Where("id = ?", id).Update("visit_count", gorm.Expr("visit_count + 1")).Error; err != nil {
		if e.Is(err, gorm.ErrRecordNotFound) {
			return ErrNoMovieFound
		}
		return err
	}
	return nil
}
