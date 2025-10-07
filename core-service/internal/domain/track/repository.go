package track

type TrackRepo interface {
	Create(track *Track) error
	FindByID(id int64) *Track
}
