package main

import (
	"context"
	"github.com/joho/godotenv"
	LibraryManagement "library"
	"library/pkg/handler"
	"library/pkg/repository"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Инициализируем базу данных
	db, err := repository.InitDB()
	if err != nil {
		log.Fatal(err)
	}

	//Проводим миграцию
	if err = repository.MigrateDB(); err != nil {
		log.Fatal(err)
	}

	repos := repository.NewRepository(db)
	hund := handler.NewHandler(repos)

	srv := new(LibraryManagement.Server)
	go func() {
		err = srv.Run("8080", hund.InitRoutes())
		if err != nil {
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-quit

	log.Println("Shutting down server...")
	if err = srv.Shutdown(context.Background()); err != nil {
		log.Println("Error on shutdown server:", err)
	}

	log.Println("Close database connection")
	if err = repository.CloseDB(); err != nil {
		log.Println("Error on close DB:", err)
	}
}
