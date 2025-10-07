package user

type User struct {
	ID       int64
	Name     string
	IsActive bool
}

func (u *User) Activate() {
	u.IsActive = true
}

func (u *User) Dectivate() {
	u.IsActive = false
}
