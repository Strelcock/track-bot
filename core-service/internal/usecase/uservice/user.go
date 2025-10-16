package uservice

import (
	"core-service/config"
	"core-service/internal/domain/user"
	"database/sql"
	"errors"
	"strconv"
)

var adminId = config.Load().AdminId

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

	if strconv.FormatInt(user.Get().ID, 10) == adminId {
		err = s.repo.Admin(user.Get().ID)
		if err != nil {
			return err
		}
	}

	err = s.repo.Activate(user.Get().ID)
	return err
}

func (s *UserService) IsAdmin(id int64) (bool, error) {
	foundUser, err := s.repo.FindByID(id)
	if err != nil {
		return false, err
	}
	return foundUser.Get().IsAdmin, nil
}
