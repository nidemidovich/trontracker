package main

import (
	"fmt"
	"log"
	"os"

	"github.com/nidemidovich/trontracker/internal/app/service"
)

func main() {
	if err := service.Run(); err != nil {
		log.Println(fmt.Errorf("error running service: %w", err))
		os.Exit(1)
	}
}
