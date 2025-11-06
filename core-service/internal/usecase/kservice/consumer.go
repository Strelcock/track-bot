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

	tracks, err := c.TrackRepo.FindByNumber(res.Barcode)
	if err != nil {
		return nil, err
	}

	var to []int64

	for _, t := range tracks {
		to = append(to, t.Get().User)
	}

	return &StatusTo{
			Status: res,
			To:     to},
		nil
}
