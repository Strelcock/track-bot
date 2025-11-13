package track

type Track struct {
	number string
	user   int64
}

type TrackSnap struct {
	Number string
	User   int64
}

func New(number string, user int64) *Track {
	return &Track{number: number, user: user}
}

func (t *Track) Get() *TrackSnap {
	return &TrackSnap{t.number, t.user}
}
