package postgres

import (
	"core-service/internal/domain/user"
	"time"
)

type userModel struct {
	id        int64
	name      string
	isActive  bool      `db:"is_active"`
	createdAt time.Time `db:"created_at"`
	updatedAt time.Time `db:"updated_at"`
	deletedAt time.Time `db:"deleted_at"`
}

func userToModel(u *user.User) *userModel {
	return &userModel{
		id:        u.Get().ID,
		name:      u.Get().Name,
		isActive:  u.Get().IsActive,
		createdAt: time.Now(),
		updatedAt: time.Now(),
	}
}

func (m *userModel) toUser() *user.User {
	return user.New(m.id, m.name)
}
