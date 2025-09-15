package main

import (
	"order-manager/internal/api"
)

// @title order-manager
// @version 1.0
// @description This is a order microservice
// @host localhost:8081
// @BasePath /
func main() {
	api.NewApp().RunApp()
}
