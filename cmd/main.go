package main

import (
	"log/slog"
	"os"
	"subs_service/internal/app"
)

// @title Subs service on Go
// @version 1.0
// @description API на Go для отслеживания подписок и суммы стоимости с фильтрами
// @host localhost:8080
// @BasePath /
// @schemes http

func main() {
	App, err := app.New()
	if err != nil {
		slog.Error("Failed to build app: ", "Erorr", err)
		os.Exit(1)
	}

	err = App.Run()
	if err != nil {
		slog.Error("Failed to run app: ", "Erorr", err)
		os.Exit(1)
	}
}
