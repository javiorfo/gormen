package model

type User struct {
	ID       uint
	Username string
	Password string
	Enable   bool
	Person   Person
}
