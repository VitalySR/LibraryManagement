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

/*
Для книг:
• POST/books — Добавить новую книгу;
• GET /books — Получить все книги;
• GET /books/{id} — Получить книгу по ее идентификатору;
• PUT /books/{id} — обновить книгу по ее идентификатору;
• DELETE /books/{id} — Удалить книгу по ее идентификатору.
Для авторов:
• POST/authors — Добавить нового автора;
• GET /authors — Получить всех авторов;
• GET /authors/{id} — получить автора по его идентификатору;
• PUT /authors/{id} — обновить автора по его идентификатору;
• DELETE /authors/{id} — удалить автора по его идентификатору.
Транзакционное обновление:
• PUT /books/{book_id}/authors/{author_id} — одновременно обновить сведения
о книге и авторе.
*/
