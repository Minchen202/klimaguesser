package main

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func initDB(dsn string) error {
	var err error
	db, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	return db.AutoMigrate(&User{}, &LeaderboardEntry{})
}

var errUserExists = errors.New("username already exists")
var errUserNotFound = errors.New("user not found")

func createUser(username, password string) error {
	var existing User
	if err := db.Where("username = ?", username).First(&existing).Error; err == nil {
		return errUserExists
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := User{Username: username, Password: string(hashed)}
	return db.Create(&user).Error
}

func findUserByUsername(username string) (*User, error) {
	var user User
	if err := db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, errUserNotFound
	}
	return &user, nil
}

func findUserByToken(token string) (*User, error) {
	var user User
	if err := db.Where("token = ?", token).First(&user).Error; err != nil {
		return nil, errUserNotFound
	}
	return &user, nil
}

func setUserToken(username, token string) error {
	return db.Model(&User{}).Where("username = ?", username).Update("token", token).Error
}

func saveLeaderboardEntry(entry *LeaderboardEntry) error {
	return db.Create(entry).Error
}

func fetchLeaderboard() ([]LeaderboardEntry, error) {
	var entries []LeaderboardEntry
	err := db.Order("score desc").Find(&entries).Error
	return entries, err
}
