package postgres

import (
	"core-service/internal/domain/track"
	"core-service/internal/domain/user"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Database struct {
	*UserDb
	*TrackDb
}

type UserDb struct{ *sqlx.DB }
type TrackDb struct{ *sqlx.DB }

const driver = "postgres"

func New(dsn string) *Database {

	for i := range 5 {
		time.Sleep(time.Second * time.Duration(i))
		db, err := sqlx.Connect(driver, dsn)
		if err == nil {
			log.Println("Service connected to db successfully")
			return &Database{&UserDb{db}, &TrackDb{db}}
		}
		log.Println("cant connect to db: ", err)
	}

	return nil
}

// USERS#########################################
func (db *UserDb) Create(u *user.User) error {

	model := userToModel(u)
	_, err := db.DB.NamedQuery("INSERT INTO users (id, name, is_active) VALUES (:id, :name, :is_active)", model)

	return err
}

func (db *UserDb) FindByID(id int64) (*user.User, error) {
	model := &userModel{}

	err := db.DB.Get(model, "SELECT id, name, is_active, is_admin FROM users WHERE id = $1", id)
	if err != nil {
		return nil, err
	}

	return model.toUser(), nil
}

func (db *UserDb) Activate(id int64) error {
	_, err := db.DB.Exec("UPDATE users SET is_active=TRUE WHERE id = $1", id)
	return err
}

func (db *UserDb) Deactivate(id int64) error {
	_, err := db.DB.Exec("UPDATE users SET is_active=FALSE WHERE id = $1", id)
	return err
}

func (db *UserDb) Admin(id int64) error {
	_, err := db.DB.Exec("UPDATE users SET is_admin=TRUE WHERE id = $1", id)
	return err
}

//TRACKS#########################################

func (db *TrackDb) Create(track *track.Track) error {
	model := trackToModel(track)
	_, err := db.DB.NamedQuery("INSERT INTO tracks (number, user_id) VALUES (:number, :user_id)", model)

	return err
}

func (db *TrackDb) FindByNumber(number string) ([]track.Track, error) {
	model := []trackModel{}

	err := db.DB.Select(&model, "SELECT number, user_id FROM tracks WHERE TRIM(number) = $1", number)
	if err != nil {
		return nil, err
	}
	res := []track.Track{}
	for _, m := range model {
		res = append(res, *m.toTrack())
	}
	return res, nil
}
