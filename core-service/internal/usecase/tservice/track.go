package tservice

import (
	"core-service/internal/domain/track"
	"errors"
	"strings"
)

type TrackService struct {
	repo track.TrackRepo
}

func New(repo track.TrackRepo) *TrackService {
	return &TrackService{repo: repo}
}

func (t *TrackService) Create(tracks []track.Track) error {

	var errs = []string{}
	for _, tr := range tracks {
		err := t.repo.Create(&tr)
		if err != nil {
			errs = append(errs, err.Error())
		}
	}

	if len(errs) != 0 {
		errStr := strings.Join(errs, ";")
		err := errors.New(errStr)
		return err
	}
	return nil
}
