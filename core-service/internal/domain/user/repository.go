package user

type UserRepo interface {
	Create(user *User) error
	FindByID(id int64) *User
}
