package main

import (
	"log"

	"github.com/isyayanoaa/go-blog/internal/router"
)

func main() {
	r := router.Setup()
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
