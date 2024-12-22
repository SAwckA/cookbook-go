package main

import (
	"cookbook/database"
	_ "cookbook/docs"
	httpserver "cookbook/http_server"
	"os"

	"github.com/joho/godotenv"
)

// @SecurityDefinitions.apikey	CookieSID
// @in							cookie
// @name						sid
func main() {

	_ = godotenv.Load()

	db := database.New()

	server := httpserver.New(
		db,
	)

	server.Listen(os.Getenv("LISTEN_ADDR"))
}
