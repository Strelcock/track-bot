package postgres

import (
	"core-service/internal/domain/user"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Database struct {
	*sqlx.DB
}

const driver = "postgres"

func New(dsn string) *Database {
	db, err := sqlx.Connect(driver, dsn)
	if err != nil {
		panic(err)
	}

	return &Database{db}
}

func (db *Database) Create(u *user.User) error {

	model := userToModel(u)
	_, err := db.DB.NamedQuery("INSERT INTO users (id, name, is_active) VALUES (:id, :name, :is_active)", model)

	return err
}

func (db *Database) FindByID(id int64) (*user.User, error) {
	model := &userModel{}

	err := db.DB.Get(model, "SELECT id, name, is_active FROM users WHERE id = $1", id)
	if err != nil {
		return nil, err
	}

	return model.toUser(), nil
}

func (db *Database) Activate(id int64) error {
	_, err := db.DB.Exec("UPDATE TABLE users SET is_active=TRUE WHERE id = $1", id)
	return err
}

func (db *Database) Deactivate(id int64) error {
	_, err := db.DB.Exec("UPDATE TABLE users SET is_active=FALSE WHERE id = $1", id)
	return err
}
