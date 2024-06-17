package handler

import (
	"encoding/json"
	"fmt"
	"library/pkg/repository"
	"log"
	"net/http"
)

type Hundler struct {
	repository *repository.Repository
}

func NewHandler(repository *repository.Repository) *Hundler {
	return &Hundler{repository: repository}
}

func (h *Hundler) InitRoutes() {
	http.HandleFunc("/books", h.bookWorker)
}

func (h *Hundler) RunServer() error {
	return http.ListenAndServe(":8080", nil)
}

func (h *Hundler) bookWorker(w http.ResponseWriter, r *http.Request) {
	log.Printf("Call link %s, method: %s", r.RequestURI, r.Method)

	switch r.Method {
	case http.MethodGet:
		log.Println("Processing method GET")
		books, err := h.repository.BookWorker.GetAll()
		if err != nil {
			log.Printf("Error getting books: %v", err)
			http.Error(w, "Error getting books", http.StatusInternalServerError)
			return
		}

		//err = json.NewEncoder(w).Encode(books)
		booksJSON, err := json.Marshal(&books)
		if err != nil {
			log.Printf("Error encoding books: %v", err)
			http.Error(w, "Error encoding books", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(booksJSON)
		if err != nil {
			log.Printf("Error writing books: %v", err)
			http.Error(w, "Error writing books", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	case http.MethodPost:
		book := &repository.Book{}

		err := json.NewDecoder(r.Body).Decode(book)
		if err != nil {
			log.Println("Error decode input:", r.Body, " - ", err)
			http.Error(w, "Wrong JSON request", http.StatusBadRequest)
			return
		}
		log.Printf("Getting book: %+v\n", book)

		if book.Title == "" {
			http.Error(w, "Title is mandatory field", http.StatusNotAcceptable)
			return
		}

		id, err := h.repository.BookWorker.Create(*book)
		if err != nil {
			log.Println("Error to create book in db:", err)
			http.Error(w, "Error to create book", http.StatusInternalServerError)
			return
		}
		log.Printf("Created book with id %d", id)

		w.WriteHeader(http.StatusCreated)
		_, err = fmt.Fprint(w, "Book created successfully with id = ", id)
		if err != nil {
			log.Println("Error write answer:", err)
		}
	default:
		log.Println("Unsupported method")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
