package main

import (
	dotenv "github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"os"
)


func run() {
	user := os.Getenv("DATABASE_USER")
	password := os.Getenv("DATABASE_PASSWORD")
	dbname := os.Getenv("DATABASE_NAME")
	logPath := os.Getenv("LOG_PATH")
	addr := os.Getenv("RUN_ADDR")
	app := App{}
	app.Initialize(user, password, dbname, logPath)
	app.Run(addr)
}

func init() {
	if err := dotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	run()
}
