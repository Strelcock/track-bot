package tservice

import (
	"core-service/internal/domain/track"
	"fmt"
)

type TrackService struct {
	repo track.TrackRepo
}

func New(repo track.TrackRepo) *TrackService {
	return &TrackService{repo: repo}
}

func (t *TrackService) Create(track *track.Track, commit chan bool) error {
	err := t.repo.Create(track, commit)
	if err != nil {
		return fmt.Errorf("trackService: %w", err)
	}

	return nil
}

func (t *TrackService) GetInfo(user int64) ([]string, error) {
	tracks, err := t.repo.GetInfo(user)
	if err != nil {
		return nil, err
	}
	res := []string{}
	for _, num := range tracks {
		res = append(res, num.Get().Number)
	}
	return res, nil
}
