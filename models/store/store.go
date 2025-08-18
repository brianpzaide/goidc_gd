package store

import (
	"goidc_gd/models"
	"goidc_gd/models/store/sqlite"
)

type UserModelInterface interface {
	CreateUser(sub, email, name string) (models.User, error)
	UserExists(sub string) (models.User, error)
}

type Models interface {
	UserModelInterface
	Close()
}

func New(dsn string) (Models, error) {
	m, err := sqlite.NewSqliteModel(dsn)
	if err != nil {
		return nil, err
	}
	return m, nil
}
