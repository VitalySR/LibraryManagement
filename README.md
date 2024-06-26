# Library Management System

## Компиляция и запуск

1. Скопировать репозиторий: `git clone https://github.com/VitalySR/LibraryManagement.git`
2. Перейти в директорию проекта: `cd LibraryManagement`
3. Для компиляции Docker образа: `docker-compose build`
4. Запуск приложения: `docker-compose up -d`. Флаг `-d` можно опустить, если требуется оставить контейнеры после остановки
5. Остановка приложения: `docker-compose down`

## API Endpoints

### Для книг

* `POST/books`: Добавить новую книгу
* `GET /books`: Получить все книги
* `GET /books/{id}`: Получить книгу по ее идентификатору
* `PUT /books/{id}`: обновить книгу по ее идентификатору
* `DELETE /books/{id}`: Удалить книгу по ее идентификатору

### Для авторов

* `POST/authors`: Добавить нового автора;
* `GET /authors`: Получить всех авторов;
* `GET /authors/{id}`: получить автора по его идентификатору;
* `PUT /authors/{id}`: обновить автора по его идентификатору;
* `DELETE /authors/{id}`: удалить автора по его идентификатору.

### Транзакционное обновление

* `PUT /books/{book_id}/authors/{author_id}`: одновременно обновить сведения
о книге и авторе.

## JSON формат

### Для авторов
```json
  {
  "id": <уникальный идентификатор автора>,
  "first_name": <имя автора>,
  "last_name": <фамилия автора>,
  "biography": <краткое описание деятельности автора>,
  "birth_date": <дата в формате yyyy-mm-dd>
  }
 ```
"first_name" и "last_name" являются обязательными полями
### Для книг
```json
{
  "id": <уникальный идентификатор книги>,
  "title": <название книги>,
  "author": {
    "id": <уникальный идентификатор автора>,
    "first_name": <имя автора>,
    "last_name": <фамилия автора>,
    "biography": <краткое описание деятельности автора>,
    "birth_date": <дата в формате yyyy-mm-dd>
  },
  "year": <Год издания книги>,
  "isbn": <Уникальный ISBN книги>
}
```
"title" является обязательным полем

### Примеры запросов
<details>
  <summary>Для книг</summary>

Создание книги без автора:
```json
{
    "title": "Как закалялась сталь",
    "year": 1934,
    "isbn": "978-5-17-121544-6"
}
```

Создание книги с автором (при наличии его в базе):
```json
{
    "title": "Книга 2",
    "year": 2013,
    "isbn": "4545-5565",
    "author": {
            "id": 1
    }
}
```
Обновление данных по книге (необходимо передавать все поля, если не требуется их затереть)
```json
{
    "Id": 3,
    "title": "Книга 22",
    "year": 2013,
    "author": {
            "id": 1
    }
}
```
</details>

<details>
    <summary>Для авторов</summary>

Создание автора
```json
{
    "first_name": "Иван",
    "last_name": "Иванов",
    "biography": "Может быть у него и биография есть",
    "birth_date": "1985-06-30"
}
```

Обновление автора
```json
{
  "id": 2,
  "first_name": "Петр",
  "last_name": "Петров",
  "biography": "Самобытный писатель. Начинал печататься со статей в газетах",
  "birth_date": "1974-12-01"
}
```
</details>

<details>
    <summary>Одновременное обновление автора и книги</summary>

```json
{
    "id": 1,
    "title": "Книга",
    "author": {
        "id": 1,
        "first_name": "Александр",
        "last_name": "Александров",
        "biography": "Писатель. Поэт",
        "birth_date": "1974-12-01"
    },
    "year": 2013,
    "isbn": "4545-5565"
}
```
</details>
