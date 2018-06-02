package models

import "github.com/jinzhu/gorm"

const connectionString = "postgresql://postgres:development@db:5432/stereodose?sslmode=disable"

type StereoDoseDB struct {
	db    *gorm.DB
	Users UserService
}

// NewStereodoseDB takes a reference to gorm and returns
// an abstraction for use throughout the app
func NewStereodoseDB(db *gorm.DB) *StereoDoseDB {
	db.AutoMigrate(User{})
	database := &StereoDoseDB{}
	database.db = db
	database.Users = &StereodoseUserService{db}
	return database
}

// func GetDB() *gorm.DB {
// 	if db != nil {
// 		return db
// 	}
// 	db, err := gorm.Open("postgres", connectionString)
// 	if err != nil {
// 		panic(err.Error())
// 	}
// 	return db
// }
