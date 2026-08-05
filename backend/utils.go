package main

import (
	"math"

	"golang.org/x/crypto/bcrypt"
)

func absFloat(f float64) float64 {
	if f < 0 {
		return -f
	}
	return f
}

func roundFloat(f float64) float64 {
	return math.Round(f)
}

func checkPasswordHash(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
