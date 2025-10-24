package entities

import "hex-arch-fiber/application/model"

type UserDB struct {
	ID       uint     `gorm:"primaryKey;autoIncrement"`
	Username string   `gorm:"not null"`
	Password string   `gorm:"not null"`
	Enable   bool     `gorm:"not null"`
	PersonID uint     `gorm:"not null"`
	Person   PersonDB `gorm:"column:person_id;not null"`
}

func (udb UserDB) TableName() string {
	return "users"
}

func (udb *UserDB) From(u model.User) {
	udb.ID = u.ID
	udb.Username = u.Username
	udb.Password = u.Password
	udb.Enable = u.Enable
	udb.Person.From(u.Person)
}

func (udb UserDB) Into() model.User {
	return model.User{
		ID:       udb.ID,
		Username: udb.Username,
		Password: udb.Password,
		Enable:   udb.Enable,
		Person:   udb.Person.Into(),
	}
}

type UserFilter struct {
	Username    string `filter:"username = ?"`
	PersonEmail string `filter:"people.email in (?);join:inner join people on people.id = users.person_id"`
}
