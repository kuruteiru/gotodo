package models

type User struct {
    ID       uint64
    Username string
    Email    string
    Password string
}

func NewUser(id uint64, username, email, password string) User {
    return User{
        ID:       id
        Username: username
        Email:    email
        Password: password
    }
}