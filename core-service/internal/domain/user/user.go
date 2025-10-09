package user

type User struct {
	id       int64
	name     string
	isActive bool
}

type UserSnap struct {
	ID       int64
	Name     string
	IsActive bool
}

func New(id int64, name string) *User {
	return &User{id, name, true}
}

func (u *User) Activate() {
	u.isActive = true
}

func (u *User) Dectivate() {
	u.isActive = false
}

func (u *User) Get() *UserSnap {
	return &UserSnap{u.id, u.name, u.isActive}
}
