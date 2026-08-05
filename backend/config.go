package main

import "os"

type Config struct {
	SecretKey   string
	DatabaseDSN string
	Port        string
	Debug       bool
	LogSecret   string
	AllowedCORS []string
}

func loadConfig() Config {
	dsn := os.Getenv("DATABASE_DSN")
	if dsn == "" {
		dsn = "klimaguessr.db" 
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	return Config{
		SecretKey:   os.Getenv("secretkey"),
		DatabaseDSN: dsn,
		Port:        port,
		Debug:       os.Getenv("DEBUG") == "true" || os.Getenv("DEBUG") == "1",
		LogSecret:   os.Getenv("log_secret"),
		AllowedCORS: []string{
			"https://klimaguessr.cns-studios.com",
			"https://klima-test.cns-studios.com",
			"http://klima-test.cns-studios.com",
			"http://localhost:" + port,
			"http://127.0.0.1:" + port,
		},
	}
}
