package user

type User struct {
	id       int64
	name     string
	isActive bool
	isAdmin  bool
}

type UserSnap struct {
	ID       int64  `db:"id"`
	Name     string `db:"name"`
	IsActive bool   `db:"is_active"`
	IsAdmin  bool   `db:"is_admin"`
}

func New(id int64, name string, is_admin bool) *User {
	return &User{id, name, true, is_admin}
}

func (u *User) Activate() {
	u.isActive = true
}

func (u *User) Dectivate() {
	u.isActive = false
}

func (u *User) Get() *UserSnap {
	return &UserSnap{u.id, u.name, u.isActive, u.isAdmin}
}
