package main

import "time"


type User struct {
	ID       uint   `gorm:"primaryKey"`
	Username string `gorm:"uniqueIndex;not null"`
	Password string `gorm:"not null"` 
	Token    string `gorm:"index"`
}


type LeaderboardEntry struct {
	ID        uint      `gorm:"primaryKey"`
	Username  string    `gorm:"not null"`
	Score     int       `gorm:"not null"`
	Timestamp time.Time `gorm:"not null"`
}
