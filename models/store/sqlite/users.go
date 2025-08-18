package sqlite

import (
	"database/sql"
	"goidc_gd/models"
)

func (m *SqliteModel) CreateUser(id, email, name string) (models.User, error) {
	db, err := getDBConnection(m.dsn)
	if err != nil {
		return models.User{}, err
	}
	defer db.Close()
	_, err = db.Exec(create_user, id, email, name)
	if err != nil {
		return models.User{}, err
	}

	return models.User{
		ID:    id,
		Email: email,
		Name:  name,
	}, nil
}

func (m *SqliteModel) UserExists(sub string) (models.User, error) {

	db, err := getDBConnection(m.dsn)
	if err != nil {
		return models.User{}, err
	}
	defer db.Close()

	var (
		idStr, emailStr, nameStr sql.NullString
	)

	user := models.User{}

	err = db.QueryRow(user_exists, sub).Scan(&idStr, &emailStr, &nameStr)
	if err != nil {
		return models.User{}, err
	}

	if idStr.Valid && emailStr.Valid && nameStr.Valid {
		user.ID = idStr.String
		user.Email = emailStr.String
		user.Name = nameStr.String
		return user, nil
	}
	return models.User{}, err
}
