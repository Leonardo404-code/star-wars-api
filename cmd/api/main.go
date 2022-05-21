package main

import (
	"ame-challenge/config"
	"ame-challenge/internal/database"
	"ame-challenge/internal/routes"
	"fmt"
	"log"
)

func init() {
	config.LoadDotEnv()

	database.DBConn = database.InitDb()
}

func main() {
	r := routes.SetupRoutes()

	port := "3000"

	err := r.Run(fmt.Sprintf(":%s", port))

	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
