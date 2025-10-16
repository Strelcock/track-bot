package tservice

import "core-service/internal/domain/track"

type TrackService struct {
	repo track.TrackRepo
}

func New(repo track.TrackRepo) *TrackService {
	return &TrackService{repo: repo}
}

func (t *TrackService) Create(track *track.Track) error {
	return nil
}
