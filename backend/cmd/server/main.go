package main

import (
	"fmt"
	"log"

	"github.com/isyayanoaa/go-blog/internal/config"
	"github.com/isyayanoaa/go-blog/internal/router"
	"github.com/isyayanoaa/go-blog/pkg/database"
	"github.com/isyayanoaa/go-blog/pkg/rdb"
)

func main() {
	if err := config.Load("config.yaml"); err != nil {
		log.Fatal("load config:", err)
	}

	if err := database.Init(); err != nil {
		log.Fatal("db init:", err)
	}
	log.Println("PostgreSQL connected")

	rdb.Init()
	log.Println("Redis connected")

	r := router.Setup()
	addr := fmt.Sprintf(":%d", config.C.Server.Port)
	log.Printf("Server running on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatal(err)
	}
}
