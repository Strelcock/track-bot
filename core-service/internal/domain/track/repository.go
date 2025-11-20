package track

import "context"

type TrackRepo interface {
	Create(track *Track, commit chan bool) error
	FindByNumber(number string) ([]Track, error)
	GetInfo(user int64) ([]Track, error)
}

type TrackConsumer interface {
	Read(ctx context.Context) ([]byte, error)
}

type TrackProducer interface {
	Write(ctx context.Context, msg []byte) error
}
