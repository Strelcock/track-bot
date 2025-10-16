package track

type TrackRepo interface {
	Create(track *Track) error
	FindByNumber(number string) (*Track, error)
}
