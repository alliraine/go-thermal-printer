package main

import (
	"log"

	"github.com/jonasclaes/go-thermal-printer/pkg/bootstrap"
)

// @title			Thermal Printer API
// @version		1.0.0
// @description	Thermal Printer API written in Go, which uses ESC/POS commands.
// @host			localhost:8080
// @basepath		/
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name X-API-Key
func main() {
	err := bootstrap.Bootstrap()
	if err != nil {
		log.Fatalf("failed to bootstrap application: %v", err)
	}
}
