package kservice

import (
	"context"
	"core-service/internal/domain/track"
)

type Consumer struct {
	track.TrackConsumer
}

func New(cons track.TrackConsumer) *Consumer {
	return &Consumer{cons}
}

func Read(ctx context.Context) {

}
