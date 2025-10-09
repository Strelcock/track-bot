package postgres

import (
	"core-service/internal/domain/user"

	"github.com/jmoiron/sqlx"
)

type Database struct {
	*sqlx.DB
}

const driver = "postgres"

func New(dsn string) (*Database, error) {
	db, err := sqlx.Connect(driver, dsn)
	if err != nil {
		panic(err)
	}

	return &Database{db}, nil
}

func (db *Database) Create(u *user.User) error {

	model := userToModel(u)
	db.DB.NamedQuery("INSERT INTO users (name, is_active) VALUES (:name, :is_active)", model)
	return nil
}
