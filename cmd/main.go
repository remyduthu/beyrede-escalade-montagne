package main

import (
	"os"

	"github.com/joho/godotenv"
	beyredeescalademontagne "github.com/remyduthu/beyrede-escalade-montagne"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	server, err := beyredeescalademontagne.New(beyredeescalademontagne.Config{
		DBDir: os.Getenv("DATABASE_DIRECTORY"),
	})
	if err != nil {
		panic(err)
	}

	if err := server.UpdateUser(os.Getenv("ADMIN_USERNAME"), os.Getenv("ADMIN_PASSWORD")); err != nil {
		panic(err)
	}

	// TODO(remyduthu): Start HTTP server.
}
