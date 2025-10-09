package user

type UserRepo interface {
	Create(user *User) error
	FindByID(id int64) (*User, error)
	Activate(id int64) error
	Deactivate(id int64) error
}
