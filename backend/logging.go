package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var logger *log.Logger




func setupLogging() {
	_ = os.MkdirAll("logs", 0o755)

	archiveOldLog()

	var writer io.Writer = os.Stdout
	if logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644); err == nil {
		writer = io.MultiWriter(os.Stdout, logFile)
	}

	logger = log.New(writer, "", 0)
}

func archiveOldLog() {
	data, err := os.ReadFile("app.log")
	if err != nil {
		writeFreshLogHeader()
		return
	}

	lines := strings.SplitN(string(data), "\n", 2)
	firstLine := strings.TrimPrefix(lines[0], "#")
	safeDate := strings.ReplaceAll(firstLine, ":", "-")

	if safeDate != "" {
		archivePath := filepath.Join("logs", safeDate+".log")
		_ = os.WriteFile(archivePath, data, 0o644)
	}

	writeFreshLogHeader()
}

func writeFreshLogHeader() {
	header := "#" + time.Now().Format("2006-01-02 15:04:05") + "\n"
	_ = os.WriteFile("app.log", []byte(header), 0o644)
}

func logLine(level, format string, args ...interface{}) string {
	msg := fmt.Sprintf(format, args...)
	return fmt.Sprintf("%s [%s] %s", time.Now().Format("2006-01-02 15:04:05"), level, msg)
}

func logInfo(format string, args ...interface{}) {
	if logger == nil {
		return
	}
	logger.Println(logLine("INFO", format, args...))
}

func logWarn(format string, args ...interface{}) {
	if logger == nil {
		return
	}
	logger.Println(logLine("WARNING", format, args...))
}

func logError(format string, args ...interface{}) {
	if logger == nil {
		return
	}
	logger.Println(logLine("ERROR", format, args...))
}
