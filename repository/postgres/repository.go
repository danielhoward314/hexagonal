package postgres

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/danielhoward314/hexagonal/shortener"
	"github.com/pkg/errors"
)

type postgresRepository struct {
	client *gorm.DB
}

func NewPostgresRepository(DbHost, DbUser, DbPassword, DbName, DbPort string) (shortener.RedirectRepository, error) {
	url := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", DbHost, DbUser, DbPassword, DbName, DbPort)
	client, err := gorm.Open(postgres.Open(url), &gorm.Config{})
	if err != nil {
		return nil, errors.Wrap(err, "repository.NewPostgresRepo")
	}
	client.Debug().AutoMigrate(&shortener.Redirect{})
	repo := &postgresRepository{client: client}
	return repo, nil
}

func (r *postgresRepository) Find(code string) (*shortener.Redirect, error) {
	redirect := &shortener.Redirect{}
	err := r.client.Model(shortener.Redirect{}).Where("code = ?", code).Take(redirect).Error
	if err != nil {
		return nil, errors.Wrap(err, "repository.NewPostgresRepo")
	}
	return redirect, nil
}

func (r *postgresRepository) Store(redirect *shortener.Redirect) error {
	err := r.client.Create(redirect).Error
	if err != nil {
		return err
	}
	return nil
}
