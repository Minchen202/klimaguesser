package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	engineTypes "github.com/zishang520/engine.io/v2/types"
	"github.com/zishang520/socket.io/v2/socket"
)

func main() {
	_ = godotenv.Load()

	setupLogging()
	logInfo("Starting application...")

	cfg := loadConfig()

	if err := initDB(cfg.DatabaseDSN); err != nil {
		logError("Critical error during database initializa<qqqq2tion: %v", err)
	} else {
		logInfo("Database tables initialized successfully.")
	}

	if err := loadClimateData("../climate_data.json"); err != nil {
		logError("Failed to load climate_data.json: %v", err)
	}

	opts := socket.DefaultServerOptions()
	opts.SetServeClient(false)

	origins := make([]any, len(cfg.AllowedCORS))
	for i, o := range cfg.AllowedCORS {
		origins[i] = o
	}
	opts.SetCors(&engineTypes.Cors{
		Origin:      origins,
		Credentials: true,
	})

	io := socket.NewServer(nil, opts)
	registerSocketHandlers(io)
	defer io.Close(nil)

	router := mux.NewRouter()

	router.PathPrefix("/socket.io/").Handler(io.ServeHandler(nil))

	router.HandleFunc("/is_lobby_joinable", handleIsLobbyJoinable).Methods(http.MethodPost)
	router.HandleFunc("/closest_loc", handleClosestLoc).Methods(http.MethodPost)
	router.HandleFunc("/logs", handleLogs(cfg)).Methods(http.MethodGet)
	router.HandleFunc("/server", handleServerInfo(cfg)).Methods(http.MethodGet)

	router.HandleFunc("/", serveTemplate("index.html")).Methods(http.MethodGet)
	router.HandleFunc("/multiplayerhost", serveTemplate("multiplayerhost.html")).Methods(http.MethodGet)
	router.HandleFunc("/multiplayer", serveTemplate("multiplayer.html")).Methods(http.MethodGet)
	router.HandleFunc("/singleplayer", serveTemplate("singleplayer.html")).Methods(http.MethodGet)
	router.HandleFunc("/legal", serveTemplate("legal.html")).Methods(http.MethodGet)
	router.HandleFunc("/singleplayerlegacy", serveTemplate("singleplayerlegacy.html")).Methods(http.MethodGet)
	router.HandleFunc("/climamap", serveTemplate("climamap.html")).Methods(http.MethodGet)
	router.HandleFunc("/test", serveTemplate("test.html")).Methods(http.MethodGet)

	router.HandleFunc("/profile.png", servePicture("profile.png")).Methods(http.MethodGet)
	router.HandleFunc("/settings.png", servePicture("settings.png")).Methods(http.MethodGet)

	logInfo("Listening on port %s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, router); err != nil {
		logError("Server failed: %v", err)
	}
}
