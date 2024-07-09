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
	"strconv"
	"syscall"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	//Получаем значения переменных для работы сервера
	httpPort := os.Getenv("PORT")
	httpReadTimeout, err := strconv.Atoi(os.Getenv("HTTP_READ_TIMEOUT"))
	if err != nil || httpReadTimeout <= 0 {
		log.Fatal("HTTP_READ_TIMEOUT is wrong. See .env file")
	}
	httpWriteTimeout, err := strconv.Atoi(os.Getenv("HTTP_WRITE_TIMEOUT"))
	if err != nil || httpWriteTimeout <= 0 {
		log.Fatal("HTTP_WRITE_TIMEOUT is wrong. See .env file")
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

	srv := LibraryManagement.NewServer(httpPort, hund.InitRoutes(), httpReadTimeout, httpWriteTimeout)
	go func() {
		err = srv.Run()
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
