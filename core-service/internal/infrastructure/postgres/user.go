package postgres

import (
	"core-service/internal/domain/user"
)

type userModel struct {
	Id       int64
	Name     string `db:"name"`
	IsActive bool   `db:"is_active"`
	IsAdmin  bool   `db:"is_admin"`
	// createdAt time.Time `db:"created_at"`
	// updatedAt time.Time `db:"updated_at"`
	// deletedAt time.Time `db:"deleted_at"`
}

func userToModel(u *user.User) *userModel {
	return &userModel{
		Id:       u.Get().ID,
		Name:     u.Get().Name,
		IsActive: u.Get().IsActive,
	}
}

func (m *userModel) toUser() *user.User {
	return user.New(m.Id, m.Name, m.IsAdmin)
}
