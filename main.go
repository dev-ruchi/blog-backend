package main

import (
	"blog-backend/api"
	"blog-backend/app"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		panic("Error loading .env file")
	}

	app.SetupDatabase()

	api.SetupRoutes()

	defer app.Db.Close()
}