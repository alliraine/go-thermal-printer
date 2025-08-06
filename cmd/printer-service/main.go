package main

import (
	"log"

	"github.com/jonasclaes/go-thermal-printer/pkg/bootstrap"
)

func main() {
	err := bootstrap.Bootstrap()
	if err != nil {
		log.Fatalf("failed to bootstrap application: %v", err)
	}
}
