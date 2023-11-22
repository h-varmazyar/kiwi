package repositories

import (
	"context"
	db "github.com/h-varmazyar/kiwi/applications/proxy/pkg/db/PostgreSQL"
	"github.com/h-varmazyar/kiwi/applications/proxy/pkg/entities"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type PostRepository struct {
	log *log.Logger
	*db.DB
}

var postTableName = "posts"

func NewPostRepository(ctx context.Context, log *log.Logger, db *db.DB) (*PostRepository, error) {
	repo := &PostRepository{
		log: log,
		DB:  db,
	}

	if err := repo.migration(ctx, db); err != nil {
		return nil, err
	}

	return repo, nil
}

func (r *PostRepository) migration(_ context.Context, dbInstance *db.DB) error {
	var err error
	migrations := make(map[string]struct{})
	tags := make([]string, 0)
	err = dbInstance.PostgresDB.Table(db.MigrationTable).Where("table_name = ?", postTableName).Select("tag").Find(&tags).Error
	if err != nil {
		return err
	}

	for _, tag := range tags {
		migrations[tag] = struct{}{}
	}

	newMigrations := make([]*db.Migration, 0)
	err = dbInstance.PostgresDB.Transaction(func(tx *gorm.DB) error {
		if _, ok := migrations["v1.0.0"]; !ok {
			err = tx.AutoMigrate(new(entities.Post))
			if err != nil {
				return err
			}
			newMigrations = append(newMigrations, &db.Migration{
				TableName:   postTableName,
				Tag:         "v1.0.0",
				Description: "create posts table",
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

func (r *PostRepository) Create(_ context.Context, post *entities.Post) error {
	return r.PostgresDB.Save(post).Error
}

func (r *PostRepository) NewUnused(_ context.Context) (*entities.Post, error) {
	st := new(entities.Post)
	err := r.PostgresDB.Model(new(entities.Post)).Order("created_at ASC").Where("used = false").First(st).Error
	if err != nil {
		return nil, err
	}

	err = r.PostgresDB.Model(new(entities.Post)).Where("id = ?", st.ID).Update("used", true).Error
	if err != nil {
		return nil, err
	}
	return st, nil
}
