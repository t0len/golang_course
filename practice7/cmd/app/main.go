package main

import (
	"practice-7/config"
	"practice-7/internal/app"
)

func main() {
	cfg := config.NewConfig()
	app.Run(cfg)
}
