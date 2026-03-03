package user

type Repository interface {
	Create(user *User) error
	FindByEmail(email string) (*User, error)
	FindByID(id uint) (*User, error)
	FindAll() ([]*User, error)
	Update(user *User) error

	saveVerificationToken()
}