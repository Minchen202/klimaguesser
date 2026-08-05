package main

import (
	"crypto/rand"
	"math/big"
	"sync"
)





type Climate map[string]interface{}

func (c Climate) lat() float64 { return toFloat(c["lat"]) }
func (c Climate) lng() float64 { return toFloat(c["lng"]) }
func (c Climate) name() string {
	if v, ok := c["name"].(string); ok {
		return v
	}
	return ""
}

func toFloat(v interface{}) float64 {
	switch n := v.(type) {
	case float64:
		return n
	case int:
		return float64(n)
	}
	return 0
}


type Guess struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}


type Player struct {
	Nickname       string `json:"nickname"`
	SessionID      string `json:"session_id"`
	AlreadyGuessed bool   `json:"alreadyguessed"`
	Guess          *Guess `json:"guess"`
	Score          int    `json:"score"`
}


type Lobby struct {
	Status        string             `json:"status"` 
	Players       map[string]*Player `json:"players"`
	ActiveClimate Climate            `json:"active_climate"`
	Round         int                `json:"round"`
}


type SoloGame struct {
	CurrentRound int
	Score        int
	Climate      Climate
}

var (
	lobbiesMu     sync.Mutex
	activeLobbies = map[string]*Lobby{}

	soloGamesMu     sync.Mutex
	activeSoloGames = map[string]*SoloGame{}

	climateData []Climate
)

const lobbyCodeChars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"


func generateUniqueLobbyCode(length int) string {
	for {
		code := randomCode(length)
		lobbiesMu.Lock()
		_, exists := activeLobbies[code]
		lobbiesMu.Unlock()
		if !exists {
			return code
		}
	}
}

func randomCode(length int) string {
	b := make([]byte, length)
	for i := range b {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(lobbyCodeChars))))
		b[i] = lobbyCodeChars[n.Int64()]
	}
	return string(b)
}

func randomClimate() Climate {
	n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(climateData))))
	return climateData[n.Int64()]
}
