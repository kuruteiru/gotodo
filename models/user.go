package models

type User struct {
	ID       uint64
	Username string
	Email    string
	Password string
}

func NewUser(username, email, password string) (User, error) {
	//todo: hash passwords
	password = "hashed_password"

	return User{
		Username: username,
		Email:    email,
		Password: password,
	}, nil
}
