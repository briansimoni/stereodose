package models

import (
	"github.com/gorilla/sessions"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type UserService interface {
	ByID(ID uint) (*User, error)
	FirstOrCreate(user *User) (*User, error)
}

type StereodoseUserService struct {
	store *sessions.CookieStore
	db    *gorm.DB
}

type User struct {
	gorm.Model
	Birthdate   string
	DisplayName string
	Email       string
	// TODO: may want to change this to not unique to handle soft delete cases
	SpotifyID    string `gorm:"unique;not null"`
	RefreshToken string
	//Images      []string
}

// Me first checks to see if the user already exists
// if it doesn't it creates one, otherwise it returns a pointer to user
// TODO: probably just get the user by id (not create)
func (u *StereodoseUserService) ByID(ID uint) (*User, error) {
	user := &User{}

	err := u.db.Find(user, "id = ?", ID).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *StereodoseUserService) FirstOrCreate(user *User) (*User, error) {
	err := u.db.FirstOrCreate(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *StereodoseUserService) UpdateUser(user *User) error {
	err := u.db.Update(user).Error
	if err != nil {
		return err
	}
	return nil
}

func (u *StereodoseUserService) DeleteUser(user *User) error {
	err := u.db.Delete(user).Error
	if err != nil {
		return err
	}
	return nil
}
