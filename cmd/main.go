package main

import (
	"log"

	"main/cmd/api"
	"main/database"
	"main/internal/evn"
)

// @title Nkata API
// @version 1.0
// @description Nkata monolight server
// @host localhost:5557
// @BasePath 
func main() {

	if err := evn.InitEvn(); err != nil {
		log.Print("Failed to load .env file useing hard coded defaults")
	}

	databaseConfig := database.DatabaseConfig{
		Addr:         evn.GetString("postgres://postgres:12345@localhost/nkata?sslmode=disable", "DATABASE_ADDR"),
		MaxOpenConn:  evn.GetInt(20, "MAX_DATABASE_OPEN_CONN"),
		MaxIdealConn: evn.GetInt(20, "MAX_DATABASE_IDEAL_CONN"),
		MaxIdealTime: "15m",
	}

	db, err := database.ConnectDatabase(databaseConfig)

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()
	log.Print("Database conection established")
	api.IntiApi(db)

}
