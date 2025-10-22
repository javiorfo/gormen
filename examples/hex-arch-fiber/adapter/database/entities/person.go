package entities

import "hex-arch-fiber/application/model"

type PersonDB struct {
	ID    uint   `gorm:"primaryKey;autoIncrement"`
	Name  string `gorm:"not null"`
	Email string `gorm:"not null"`
}

func (pdb PersonDB) TableName() string {
	return "people"
}

func (pdb *PersonDB) From(p model.Person) {
	pdb.ID = p.ID
	pdb.Name = p.Name
	pdb.Email = p.Email
}

func (pdb PersonDB) Into() model.Person {
	return model.Person{
		ID:    pdb.ID,
		Name:  pdb.Name,
		Email: pdb.Email,
	}
}
