package main

import (
	"library/pkg/handler"
	"library/pkg/repository"
	"log"
)

func main() {
	// Инициализируем базу данных
	dataSourceName := "postgres://postgres:qwerty@localhost:5432/postgres?sslmode=disable"
	db, err := repository.InitDB(dataSourceName)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err := repository.CloseDB()
		if err != nil {
			log.Println("Error on close DB:", err)
		}
	}()

	//Проводим миграцию
	err = repository.MigrateDB()
	if err != nil {
		log.Fatal(err)
	}

	repos := repository.NewRepository(db)
	hund := handler.NewHandler(repos)
	hund.InitRoutes()

	err = hund.RunServer()
	if err != nil {
		log.Fatal(err)
	}
}
