package main

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"
)

func writeJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func readJSONBody(r *http.Request) map[string]interface{} {
	if r.Body == nil {
		return nil
	}
	var data map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return nil
	}
	return data
}

func serveTemplate(name string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../frontend/"+name)
	}
}

func servePicture(name string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../pictures/"+name)
	}
}

func handleIsLobbyJoinable(w http.ResponseWriter, r *http.Request) {
	data := readJSONBody(r)
	lobbyCode, _ := data["lobby_code"].(string)

	lobbiesMu.Lock()
	lobby, ok := activeLobbies[lobbyCode]
	joinable := ok && lobby.Status == "waiting"
	lobbiesMu.Unlock()

	writeJSON(w, http.StatusOK, map[string]interface{}{"joinable": joinable})
}

func handleClosestLoc(w http.ResponseWriter, r *http.Request) {
	data := readJSONBody(r)
	if data == nil {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{"error": "invalid request body"})
		return
	}
	lat := toFloat(data["lat"])
	lng := toFloat(data["lng"])

	closest := closestClimate(lat, lng)
	writeJSON(w, http.StatusOK, map[string]interface{}{"closest_climate": closest})
}

func requireLogSecret(cfg Config, w http.ResponseWriter, r *http.Request) bool {
	data := readJSONBody(r)
	secret, _ := data["secret_key"].(string)
	if data == nil || secret != cfg.LogSecret {
		logWarn("Unauthorized access attempt to logs from %s", r.RemoteAddr)
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{"success": false, "message": "Invalid request."})
		return false
	}
	return true
}

func handleLogs(cfg Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !requireLogSecret(cfg, w, r) {
			return
		}
		http.ServeFile(w, r, "app.log")
	}
}

func handleServerInfo(cfg Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !requireLogSecret(cfg, w, r) {
			return
		}

		lastRestart := ""
		if data, err := os.ReadFile("app.log"); err == nil {
			firstLine := strings.SplitN(string(data), "\n", 2)[0]
			lastRestart = strings.SplitN(firstLine, "[", 2)[0]
		}

		lobbiesMu.Lock()
		lobbiesCopy := make(map[string]*Lobby, len(activeLobbies))
		for k, v := range activeLobbies {
			lobbiesCopy[k] = v
		}
		lobbiesMu.Unlock()

		soloGamesMu.Lock()
		soloCopy := make(map[string]*SoloGame, len(activeSoloGames))
		for k, v := range activeSoloGames {
			soloCopy[k] = v
		}
		soloGamesMu.Unlock()

		writeJSON(w, http.StatusOK, map[string]interface{}{
			"success":           true,
			"cpu_usage":         cpuPercent(),
			"ram_usage":         ramPercent(),
			"last_restart":      lastRestart,
			"active_lobbies":    lobbiesCopy,
			"active_solo_games": soloCopy,
		})
	}
}
