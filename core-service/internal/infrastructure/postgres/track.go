package postgres

import (
	"core-service/internal/domain/track"
)

type trackModel struct {
	Id     int64  `db:"id"`
	Number string `db:"number"`
	User   int64  `db:"user_id"`
	// createdAt time.Time `db:"created_at"`
	// updatedAt time.Time `db:"updated_at"`
	// deletedAt time.Time `db:"deleted_at"`
}

func trackToModel(t *track.Track) *trackModel {
	return &trackModel{
		Number: t.Get().Number,
		User:   t.Get().User,
	}
}

func (m *trackModel) toTrack() *track.Track {
	return track.New(m.Number, m.User)
}
