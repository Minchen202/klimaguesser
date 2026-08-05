package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/zishang520/socket.io/v2/socket"
)




func dataArg(args []any) map[string]interface{} {
	if len(args) == 0 {
		return nil
	}
	if m, ok := args[0].(map[string]interface{}); ok {
		return m
	}
	return nil
}

func registerSocketHandlers(io *socket.Server) {
	io.On("connection", func(clients ...any) {
		client := clients[0].(*socket.Socket)

		fmt.Printf("Client %s connected\n", client.Id())
		client.Emit("connected", map[string]interface{}{"message": "Successfully connected to server"})

		client.On("disconnect", func(args ...any) {
			fmt.Printf("Client %s disconnected\n", client.Id())
		})

		client.On("start_solo_game", func(args ...any) { handleStartSoloGame(client) })
		client.On("resize_chart", func(args ...any) { handleResizeChart(client, dataArg(args)) })
		client.On("submit_solo_guess", func(args ...any) { handleSubmitSoloGuess(client, dataArg(args)) })
		client.On("start_solo_round", func(args ...any) { handleStartSoloRound(client) })
		client.On("delete_solo_game", func(args ...any) { handleDeleteSoloGame(client) })
		client.On("save_solo_game", func(args ...any) { handleSaveSoloGame(client, dataArg(args)) })
		client.On("get_leaderboard", func(args ...any) { handleGetLeaderboard(io) })

		client.On("register", func(args ...any) { handleRegister(client, dataArg(args)) })
		client.On("login", func(args ...any) { handleLogin(client, dataArg(args)) })
		client.On("authenticate", func(args ...any) { handleAuthenticate(client, dataArg(args)) })

		client.On("create_lobby", func(args ...any) { handleCreateLobby(client) })
		client.On("register_player", func(args ...any) { handleRegisterPlayer(io, client, dataArg(args)) })
		client.On("join_lobby", func(args ...any) { handleJoinLobby(io, client, dataArg(args)) })
		client.On("start_lobby", func(args ...any) { handleStartLobby(io, client, dataArg(args)) })
		client.On("make_guess", func(args ...any) { handleMakeGuess(io, client, dataArg(args)) })
		client.On("end_round", func(args ...any) { handleEndRound(io, client, dataArg(args)) })
		client.On("start_new_round", func(args ...any) { handleStartNewRound(io, client, dataArg(args)) })
		client.On("end_game", func(args ...any) { handleEndGame(io, client, dataArg(args)) })
		client.On("get_lobby_info", func(args ...any) { handleGetLobbyInfo(client, dataArg(args)) })
	})
}



func handleStartSoloGame(s *socket.Socket) {
	game := &SoloGame{CurrentRound: 1, Score: 0, Climate: randomClimate()}

	soloGamesMu.Lock()
	activeSoloGames[string(s.Id())] = game
	soloGamesMu.Unlock()

	logInfo("Solo game started")
	s.Emit("solo_game_start_response", map[string]interface{}{
		"success": true,
		"message": "Solo game started!",
		"climate": game.Climate,
	})
}

func handleResizeChart(s *socket.Socket, data map[string]interface{}) {
	soloGamesMu.Lock()
	game, ok := activeSoloGames[string(s.Id())]
	soloGamesMu.Unlock()

	if !ok {
		s.Emit("solo_guess_response", map[string]interface{}{"success": false, "message": "No active solo game found."})
		logWarn("No active solo game found")
		return
	}

	var big any
	if data != nil {
		big = data["big"]
	}

	s.Emit("resize_chart_response", map[string]interface{}{
		"climate": game.Climate,
		"big":     big,
	})
}

func handleSubmitSoloGuess(s *socket.Socket, data map[string]interface{}) {
	soloGamesMu.Lock()
	game, ok := activeSoloGames[string(s.Id())]
	soloGamesMu.Unlock()

	if !ok {
		s.Emit("solo_guess_response", map[string]interface{}{"success": false, "message": "No active solo game found."})
		logWarn("No active solo game found")
		return
	}

	actualLat := game.Climate.lat()
	actualLng := game.Climate.lng()

	var guessLat, guessLng interface{}
	var latOK, lngOK bool
	if data != nil {
		guessLat, latOK = data["guessLat"]
		guessLng, lngOK = data["guessLng"]
	}

	if data == nil || !latOK || !lngOK {
		s.Emit("solo_guess_response", map[string]interface{}{
			"success":         true,
			"message":         fmt.Sprintf("Round %d started!", game.CurrentRound),
			"current_round":   game.CurrentRound,
			"name":            game.Climate.name(),
			"score":           0,
			"actual_location": map[string]interface{}{"lat": actualLat, "lng": actualLng},
			"total_points":    game.Score,
			"points_earned":   0,
		})
		return
	}

	gLat := toFloat(guessLat)
	gLng := toFloat(guessLng)

	latDistance := absFloat(gLat-actualLat) * 111
	lonDistance := absFloat(gLng-actualLng) * 111

	latPoints := 0
	if latDistance <= 2500 {
		latPoints = int(roundFloat(2500 - latDistance))
	}
	lonPoints := 0
	if lonDistance <= 2500 {
		lonPoints = int(roundFloat(2500 - lonDistance))
	}
	points := latPoints + lonPoints

	soloGamesMu.Lock()
	game.Score += points
	soloGamesMu.Unlock()

	logInfo("Solo guess submitted")
	s.Emit("solo_guess_response", map[string]interface{}{
		"success":         true,
		"message":         fmt.Sprintf("Round %d started!", game.CurrentRound),
		"current_round":   game.CurrentRound,
		"name":            game.Climate.name(),
		"score":           game.Score,
		"actual_location": map[string]interface{}{"lat": actualLat, "lng": actualLng},
		"total_points":    game.Score,
		"points_earned":   points,
	})
}

func handleStartSoloRound(s *socket.Socket) {
	soloGamesMu.Lock()
	game, ok := activeSoloGames[string(s.Id())]
	if ok {
		game.CurrentRound++
		game.Climate = randomClimate()
	}
	soloGamesMu.Unlock()

	if !ok {
		s.Emit("solo_guess_response", map[string]interface{}{"success": false, "message": "No active solo game found."})
		logWarn("No active solo game found")
		return
	}

	logInfo("Solo round started")
	s.Emit("solo_round_started", map[string]interface{}{
		"success":       true,
		"message":       fmt.Sprintf("Round %d started!", game.CurrentRound),
		"climate":       game.Climate,
		"current_round": game.CurrentRound,
		"total_score":   game.Score,
	})
}

func handleDeleteSoloGame(s *socket.Socket) {
	soloGamesMu.Lock()
	_, ok := activeSoloGames[string(s.Id())]
	if ok {
		delete(activeSoloGames, string(s.Id()))
	}
	soloGamesMu.Unlock()

	if ok {
		s.Emit("solo_game_deleted", map[string]interface{}{"success": true, "message": "Solo game deleted successfully."})
		logInfo("Solo game not saved")
	} else {
		s.Emit("solo_game_deleted", map[string]interface{}{"success": false, "message": "No active solo game found."})
		logWarn("No active solo game found")
	}
}

func handleSaveSoloGame(s *socket.Socket, data map[string]interface{}) {
	var token string
	if data != nil {
		token, _ = data["token"].(string)
	}

	soloGamesMu.Lock()
	game, ok := activeSoloGames[string(s.Id())]
	soloGamesMu.Unlock()

	if !ok {
		logWarn("No active solo game found")
		s.Emit("save_solo_response", map[string]interface{}{"success": false, "message": "No active solo game found."})
		return
	}

	user, err := findUserByToken(token)
	if err != nil {
		logWarn("Invalid username")
		s.Emit("save_solo_response", map[string]interface{}{"success": false, "message": "Invalid username."})
		return
	}

	entry := &LeaderboardEntry{Username: user.Username, Score: game.Score, Timestamp: time.Now()}
	if err := saveLeaderboardEntry(entry); err != nil {
		s.Emit("save_solo_response", map[string]interface{}{"success": false, "message": "Failed to save game."})
		return
	}

	logInfo("Solo game saved for user: %s with score: %d", user.Username, game.Score)
	s.Emit("save_solo_response", map[string]interface{}{"success": true, "message": "Solo game saved successfully!"})
}

func handleGetLeaderboard(io *socket.Server) {
	entries, err := fetchLeaderboard()
	if err != nil {
		logError("Error fetching leaderboard: %v", err)
		return
	}

	leaderboard := make([]map[string]interface{}, 0, len(entries))
	for _, e := range entries {
		leaderboard = append(leaderboard, map[string]interface{}{
			"username":  e.Username,
			"score":     e.Score,
			"timestamp": e.Timestamp.Format("2006-01-02 15:04:05"),
		})
	}

	io.Emit("leaderboard_update", leaderboard)
}



func handleRegister(s *socket.Socket, data map[string]interface{}) {
	var username, password string
	if data != nil {
		username, _ = data["username"].(string)
		password, _ = data["password"].(string)
	}

	if username == "" || password == "" {
		s.Emit("registration_response", map[string]interface{}{"success": false, "message": "Username and password are required."})
		logWarn("Username and password missing")
		return
	}

	if err := createUser(username, password); err != nil {
		if err == errUserExists {
			s.Emit("registration_response", map[string]interface{}{"success": false, "message": "Username already exists."})
			logWarn("Username already exists")
			return
		}
		s.Emit("registration_response", map[string]interface{}{"success": false, "message": "Registration failed."})
		return
	}

	logInfo("New user registered: %s", username)
	s.Emit("registration_response", map[string]interface{}{"success": true, "message": "Registration successful!"})
}

func handleLogin(s *socket.Socket, data map[string]interface{}) {
	var username, password string
	if data != nil {
		username, _ = data["username"].(string)
		password, _ = data["password"].(string)
	}

	if username == "" || password == "" {
		s.Emit("login_response", map[string]interface{}{"success": false, "message": "Username and password are required."})
		logWarn("Username and password missing")
		return
	}

	user, err := findUserByUsername(username)
	if err != nil {
		s.Emit("login_response", map[string]interface{}{"success": false, "message": "Username does not exist."})
		logWarn("Username does not exist")
		return
	}

	if !checkPasswordHash(password, user.Password) {
		s.Emit("login_response", map[string]interface{}{"success": false, "message": "Invalid username or password."})
		logWarn("Invalid username or password")
		return
	}

	token := uuid.NewString()
	if err := setUserToken(username, token); err != nil {
		s.Emit("login_response", map[string]interface{}{"success": false, "message": "Login failed."})
		return
	}

	logInfo("User logged in: %s with token: %s", username, token)
	s.Emit("login_response", map[string]interface{}{
		"success": true, "message": "Login successful!", "token": token, "username": username,
	})
}

func handleAuthenticate(s *socket.Socket, data map[string]interface{}) {
	var token string
	if data != nil {
		token, _ = data["token"].(string)
	}

	user, err := findUserByToken(token)
	if err != nil {
		s.Emit("auth_response", map[string]interface{}{"success": false, "message": "Authentication failed. Invalid or expired token."})
		logWarn("Authentication failed. Invalid or expired token.")
		return
	}

	logInfo("User authenticated: %s", user.Username)
	s.Emit("auth_response", map[string]interface{}{
		"success":  true,
		"message":  fmt.Sprintf("Welcome, %s! Your token is valid.", user.Username),
		"username": user.Username,
	})
}



func broadcastLobbyUpdate(io *socket.Server, lobbyCode string) {
	lobbiesMu.Lock()
	lobby, ok := activeLobbies[lobbyCode]
	lobbiesMu.Unlock()
	if !ok {
		return
	}
	io.To(socket.Room(lobbyCode)).Emit("lobby_update", map[string]interface{}{
		"lobby_code": lobbyCode,
		"details":    lobby,
	})
}

func handleCreateLobby(s *socket.Socket) {
	code := generateUniqueLobbyCode(6)

	lobbiesMu.Lock()
	activeLobbies[code] = &Lobby{
		Status:  "waiting",
		Players: map[string]*Player{},
	}
	lobbiesMu.Unlock()

	s.Join(socket.Room(code))
	logInfo("Lobby created: %s.", code)

	s.Emit("lobby_created", map[string]interface{}{
		"success":    true,
		"lobby_code": code,
		"message":    "Lobby created successfully. Share this code with your friends!",
	})
}

func handleRegisterPlayer(io *socket.Server, s *socket.Socket, data map[string]interface{}) {
	var nickname, lobbyCode string
	if data != nil {
		nickname, _ = data["nickname"].(string)
		lobbyCode, _ = data["lobby_code"].(string)
	}

	logInfo("Registering player '%s' in lobby '%s'.", nickname, lobbyCode)

	if lobbyCode == "" || nickname == "" {
		s.Emit("register_result", map[string]interface{}{
			"error_code": "missing_fields", "success": false, "message": "Lobby code and nickname are required.",
		})
		logWarn("Lobby code and nickname missing")
		return
	}

	if len(nickname) > 10 {
		s.Emit("register_result", map[string]interface{}{"success": false, "message": "Nickname is too long."})
		logWarn("Nickname is too long")
		return
	}

	lobbiesMu.Lock()
	lobby, ok := activeLobbies[lobbyCode]
	if ok {
		if p, exists := lobby.Players[nickname]; exists {
			p.SessionID = string(s.Id())
		}
	}
	lobbiesMu.Unlock()

	if !ok {
		s.Emit("register_result", map[string]interface{}{
			"error_code": "lobby_not_found", "success": false, "message": "Lobby not found.",
		})
		logWarn("Lobby not found")
		return
	}

	s.Join(socket.Room(lobbyCode))
	logInfo("Player '%s' registered in lobby '%s'.", nickname, lobbyCode)
	s.Emit("register_result", map[string]interface{}{"success": true, "message": "Successfully registered."})
	broadcastLobbyUpdate(io, lobbyCode)
}

func handleJoinLobby(io *socket.Server, s *socket.Socket, data map[string]interface{}) {
	var lobbyCodeRaw, nickname string
	if data != nil {
		lobbyCodeRaw, _ = data["lobby_code"].(string)
		nickname, _ = data["nickname"].(string)
	}
	lobbyCode := strings.ToUpper(lobbyCodeRaw)

	if lobbyCode == "" || nickname == "" {
		logWarn("Lobby code and nickname missing")
		s.Emit("join_result", map[string]interface{}{"success": false, "message": "Lobby code and nickname are required."})
		return
	}

	lobbiesMu.Lock()
	lobby, ok := activeLobbies[lobbyCode]
	if !ok {
		lobbiesMu.Unlock()
		logWarn("Lobby '%s' not found.", lobbyCode)
		s.Emit("join_result", map[string]interface{}{"success": false, "message": fmt.Sprintf("Lobby '%s' not found.", lobbyCode)})
		return
	}

	if len(lobby.Players) >= 8 {
		lobbiesMu.Unlock()
		logWarn("Lobby '%s' is full.", lobbyCode)
		s.Emit("join_result", map[string]interface{}{"success": false, "message": "Lobby is full."})
		return
	}

	if len(nickname) >= 10 {
		lobbiesMu.Unlock()
		logWarn("Nickname is too long")
		s.Emit("join_result", map[string]interface{}{"success": false, "message": "Nickname is too long."})
		return
	}

	if _, exists := lobby.Players[nickname]; exists {
		lobbiesMu.Unlock()
		logWarn("Nickname already taken")
		s.Emit("join_result", map[string]interface{}{"success": false, "message": "Nickname already taken."})
		return
	}

	lobby.Players[nickname] = &Player{
		Nickname:  nickname,
		SessionID: string(s.Id()),
		Score:     0,
	}
	lobbiesMu.Unlock()

	s.Join(socket.Room(lobbyCode))
	logInfo("Player '%s' joined lobby '%s'", nickname, lobbyCode)
	s.Emit("join_result", map[string]interface{}{
		"success": true, "nickname": nickname,
		"message":    fmt.Sprintf("Successfully joined lobby '%s' as %s", lobbyCode, nickname),
		"lobby_code": lobbyCode,
	})

	broadcastLobbyUpdate(io, lobbyCode)
}

func handleStartLobby(io *socket.Server, s *socket.Socket, data map[string]interface{}) {
	var lobbyCodeRaw string
	if data != nil {
		lobbyCodeRaw, _ = data["lobby_code"].(string)
	}
	lobbyCode := strings.ToUpper(lobbyCodeRaw)

	lobbiesMu.Lock()
	lobby, ok := activeLobbies[lobbyCode]
	if ok {
		lobby.Status = "playing"
		lobby.ActiveClimate = randomClimate()
		lobby.Round = 1
	}
	lobbiesMu.Unlock()

	if !ok {
		logWarn("Lobby '%s' not found.", lobbyCode)
		s.Emit("start_result", map[string]interface{}{"success": false, "message": fmt.Sprintf("Lobby '%s' not found.", lobbyCode)})
		return
	}

	logInfo("Lobby '%s' started.", lobbyCode)
	io.To(socket.Room(lobbyCode)).Emit("game_started", map[string]interface{}{
		"success": true,
		"message": fmt.Sprintf("Game started in lobby '%s'", lobbyCode),
		"clima":   lobby.ActiveClimate,
		"round":   lobby.Round,
	})
}

func handleMakeGuess(io *socket.Server, s *socket.Socket, data map[string]interface{}) {
	var lobbyCodeRaw, nickname string
	var guessRaw map[string]interface{}
	if data != nil {
		lobbyCodeRaw, _ = data["lobby_code"].(string)
		nickname, _ = data["nickname"].(string)
		guessRaw, _ = data["guess"].(map[string]interface{})
	}
	lobbyCode := strings.ToUpper(lobbyCodeRaw)

	lobbiesMu.Lock()
	lobby, ok := activeLobbies[lobbyCode]
	if !ok || lobby.Status != "playing" {
		lobbiesMu.Unlock()
		s.Emit("guess_result", map[string]interface{}{"success": false, "message": fmt.Sprintf("Lobby '%s' not found or not started yet.", lobbyCode)})
		logWarn("Lobby '%s' not found or not started yet.", lobbyCode)
		return
	}

	player, ok := lobby.Players[nickname]
	if !ok {
		lobbiesMu.Unlock()
		logWarn("Player '%s' not found in lobby '%s'.", nickname, lobbyCode)
		s.Emit("guess_result", map[string]interface{}{"success": false, "message": fmt.Sprintf("Player '%s' not found in lobby '%s'.", nickname, lobbyCode)})
		return
	}

	if player.AlreadyGuessed {
		lobbiesMu.Unlock()
		s.Emit("guess_result", map[string]interface{}{"success": false, "message": fmt.Sprintf("Player '%s' has already guessed.", nickname)})
		logWarn("Player '%s' has already guessed in lobby '%s'.", nickname, lobbyCode)
		return
	}

	player.AlreadyGuessed = true
	player.Guess = &Guess{Lat: toFloat(guessRaw["lat"]), Lng: toFloat(guessRaw["lng"])}

	var alreadyGuessed, notGuessed []string
	for _, p := range lobby.Players {
		if p.AlreadyGuessed {
			alreadyGuessed = append(alreadyGuessed, p.Nickname)
		} else {
			notGuessed = append(notGuessed, p.Nickname)
		}
	}
	lobbiesMu.Unlock()

	s.Emit("guess_result", map[string]interface{}{"success": true, "message": "Guess submitted successfully.", "guess": guessRaw})
	logInfo("Player '%s' made a guess.", nickname)

	io.To(socket.Room(lobbyCode)).Emit("player_guessed", map[string]interface{}{
		"already_guessed": alreadyGuessed,
		"not_guessed":     notGuessed,
	})
}

func handleEndRound(io *socket.Server, s *socket.Socket, data map[string]interface{}) {
	var lobbyCodeRaw string
	if data != nil {
		lobbyCodeRaw, _ = data["lobby_code"].(string)
	}
	lobbyCode := strings.ToUpper(lobbyCodeRaw)

	lobbiesMu.Lock()
	lobby, ok := activeLobbies[lobbyCode]
	if !ok || lobby.Status != "playing" {
		lobbiesMu.Unlock()
		logWarn("Lobby '%s' not found or not in playing state.", lobbyCode)
		s.Emit("end_result", map[string]interface{}{"success": false, "message": fmt.Sprintf("Lobby '%s' not found or not in playing state.", lobbyCode)})
		return
	}

	lobby.Status = "result"

	for _, player := range lobby.Players {
		if !player.AlreadyGuessed || player.Guess == nil {
			continue
		}

		distance := calculateDistance(
			lobby.ActiveClimate.lat(), lobby.ActiveClimate.lng(),
			player.Guess.Lat, player.Guess.Lng,
		)

		var points int
		if distance > 5000 {
			points = 0
		} else {
			points = int(roundFloat(5000 - distance))
		}
		if distance < 300 {
			points = 5000
		} else if distance < 1500 {
			points = int(roundFloat(4000 + (500-distance)*2))
		} else if distance < 3000 {
			points = int(roundFloat(3000 + (1000 - distance)))
		}

		player.Score += points
	}

	lobbySnapshot := *lobby
	actualLat := lobby.ActiveClimate.lat()
	actualLng := lobby.ActiveClimate.lng()
	lobbiesMu.Unlock()

	logInfo("Round ended in lobby '%s'.", lobbyCode)
	io.To(socket.Room(lobbyCode)).Emit("round_ended", map[string]interface{}{
		"success":             true,
		"message":             fmt.Sprintf("Round ended in lobby '%s'", lobbyCode),
		"details":             lobbySnapshot,
		"actual_location_lat": actualLat,
		"actual_location_lng": actualLng,
	})
}

func handleStartNewRound(io *socket.Server, s *socket.Socket, data map[string]interface{}) {
	var lobbyCodeRaw string
	if data != nil {
		lobbyCodeRaw, _ = data["lobby_code"].(string)
	}
	lobbyCode := strings.ToUpper(lobbyCodeRaw)

	lobbiesMu.Lock()
	lobby, ok := activeLobbies[lobbyCode]
	if !ok || lobby.Status != "result" {
		lobbiesMu.Unlock()
		logWarn("Lobby '%s' not found or not in result state.", lobbyCode)
		s.Emit("new_round_result", map[string]interface{}{"success": false, "message": fmt.Sprintf("Lobby '%s' not found or not in result state.", lobbyCode)})
		return
	}

	lobby.Status = "playing"
	lobby.ActiveClimate = randomClimate()
	lobby.Round++

	allPlayers := make([]string, 0, len(lobby.Players))
	for nickname, p := range lobby.Players {
		p.AlreadyGuessed = false
		p.Guess = nil
		allPlayers = append(allPlayers, nickname)
	}
	climate := lobby.ActiveClimate
	round := lobby.Round
	lobbiesMu.Unlock()

	logInfo("New round started in lobby '%s'.", lobbyCode)
	io.To(socket.Room(lobbyCode)).Emit("new_round_started", map[string]interface{}{
		"success":        true,
		"all_players":    allPlayers,
		"message":        "New round started!",
		"active_climate": climate,
		"round":          round,
	})
}

func handleEndGame(io *socket.Server, s *socket.Socket, data map[string]interface{}) {
	var lobbyCodeRaw string
	if data != nil {
		lobbyCodeRaw, _ = data["lobby_code"].(string)
	}
	lobbyCode := strings.ToUpper(lobbyCodeRaw)

	lobbiesMu.Lock()
	lobby, ok := activeLobbies[lobbyCode]
	if !ok || lobby.Status != "result" {
		lobbiesMu.Unlock()
		logWarn("Lobby '%s' not found or not in result state.", lobbyCode)
		s.Emit("end_game_result", map[string]interface{}{"success": false, "message": fmt.Sprintf("Lobby '%s' not found or not in result state.", lobbyCode)})
		return
	}
	delete(activeLobbies, lobbyCode)
	lobbiesMu.Unlock()

	logInfo("Game ended and lobby '%s' deleted.", lobbyCode)

	io.To(socket.Room(lobbyCode)).Emit("end_game", map[string]interface{}{
		"lobby_code": lobbyCode,
		"message":    fmt.Sprintf("Game ended and lobby '%s' deleted.", lobbyCode),
	})

	s.Emit("end_game_result", map[string]interface{}{
		"success": true,
		"message": fmt.Sprintf("Game ended and lobby '%s' deleted.", lobbyCode),
	})
}

func handleGetLobbyInfo(s *socket.Socket, data map[string]interface{}) {
	if data == nil {
		logWarn("No data provided for getting lobby info.")
		s.Emit("lobby_info", map[string]interface{}{"success": false, "message": "No data provided."})
		return
	}
	rawCode, ok := data["lobby_code"]
	if !ok {
		logWarn("Lobby code is required for getting lobby info.")
		s.Emit("lobby_info", map[string]interface{}{"success": false, "message": "Lobby code is required."})
		return
	}

	lobbyCode := strings.ToUpper(fmt.Sprint(rawCode))

	lobbiesMu.Lock()
	lobby, found := activeLobbies[lobbyCode]
	lobbiesMu.Unlock()

	if found {
		players := make([]string, 0, len(lobby.Players))
		for nickname := range lobby.Players {
			players = append(players, nickname)
		}
		s.Emit("lobby_info", map[string]interface{}{
			"success": true, "lobby_code": lobbyCode, "details": lobby, "all_players": players,
		})
		logInfo("Lobby info sent for lobby '%s'.", lobbyCode)
	} else {
		s.Emit("lobby_info", map[string]interface{}{"success": false, "message": fmt.Sprintf("Lobby '%s' not found.", lobbyCode)})
		logWarn("Lobby '%s' not found.", lobbyCode)
	}
}