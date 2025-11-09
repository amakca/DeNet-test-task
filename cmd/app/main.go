package main

import (
	"denet-test-task/internal/app"
)

const configPath = "config/config.yaml"

func main() {
	app.Run(configPath)
}
