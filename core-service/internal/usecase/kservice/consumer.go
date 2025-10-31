package kservice

import (
	"context"
	"core-service/internal/domain/track"
	"core-service/pkg/json"
)

type ConsumerService struct {
	track.TrackConsumer
	track.TrackRepo
}

func NewCons(cons track.TrackConsumer, repo track.TrackRepo) *ConsumerService {
	return &ConsumerService{cons, repo}
}

func (c *ConsumerService) Read(ctx context.Context) (*StatusTo, error) {
	msg, err := c.TrackConsumer.Read(ctx)
	if err != nil {
		return nil, err
	}

	res, err := json.Unmarshal[Status](msg)
	if err != nil {
		return nil, err
	}

	track, err := c.TrackRepo.FindByNumber(res.ID)
	if err != nil {
		return nil, err
	}

	return &StatusTo{
			Status: res,
			To:     track.Get().User},
		nil
}
