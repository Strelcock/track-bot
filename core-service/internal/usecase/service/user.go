package service

import (
	"core-service/internal/domain/user"
	"database/sql"
	"errors"
)

type UserService struct {
	repo user.UserRepo
}

func New(repo user.UserRepo) *UserService {
	return &UserService{repo}
}

func (s *UserService) Create(user *user.User) error {
	_, err := s.repo.FindByID(user.Get().ID)

	if errors.Is(err, sql.ErrNoRows) {
		err = s.repo.Create(user)
		if err != nil {
			return err
		}
		return nil
	}

	if err != nil {
		return err
	}

	err = s.repo.Activate(user.Get().ID)
	return err
}
