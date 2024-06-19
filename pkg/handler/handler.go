package handler

import (
	"encoding/json"
	"fmt"
	"library/pkg/repository"
	"log"
	"net/http"
	"strconv"
)

type Hundler struct {
	repository *repository.Repository
}

func NewHandler(repository *repository.Repository) *Hundler {
	return &Hundler{repository: repository}
}

func (h *Hundler) InitRoutes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/books", h.books)
	mux.HandleFunc("/books/{id}", h.bookId)

	return mux
}

func (h *Hundler) books(w http.ResponseWriter, r *http.Request) {
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

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(books)
		if err != nil {
			log.Printf("Error encoding books: %v", err)
			http.Error(w, "Error encoding books", http.StatusInternalServerError)
			return
		}
	case http.MethodPost:
		book := &repository.Book{}

		err := json.NewDecoder(r.Body).Decode(book)
		if err != nil {
			log.Println("Error decode input:", r.Body, " - ", err)
			http.Error(w, "Wrong JSON request", http.StatusBadRequest)
			return
		}
		log.Printf("Getting book: %+v\n", book)

		if book.Title == nil || *book.Title == "" {
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

func (h *Hundler) bookId(w http.ResponseWriter, r *http.Request) {
	log.Printf("Call link %s, method: %s", r.RequestURI, r.Method)

	param := r.PathValue("id")
	bookId, err := strconv.Atoi(param)
	if err != nil {
		log.Printf("Error converting %s to int: %v\n", param, err)
		http.Error(w, "Bad book id in URL", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		log.Println("Processing method GET")
		book, err := h.repository.BookWorker.GetById(bookId)
		log.Println(book, err)
		if err != nil {
			log.Printf("Error getting book %d: %v", bookId, err)
			http.Error(w, "Error getting book", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(book)
		if err != nil {
			log.Printf("Error encoding book %d: %v", bookId, err)
			http.Error(w, "Error encoding book", http.StatusInternalServerError)
			return
		}
		return
	case http.MethodPut:
		log.Println("Processing method PUT")
		book := &repository.Book{}

		err := json.NewDecoder(r.Body).Decode(book)
		if err != nil {
			log.Println("Error decode input:", r.Body, " - ", err)
			http.Error(w, "Wrong JSON request", http.StatusBadRequest)
			return
		}
		log.Printf("Getting book: %+v\n", book)

		if book.Title == nil || *book.Title == "" {
			http.Error(w, "Title is mandatory field", http.StatusNotAcceptable)
			return
		}
		bkId := int32(bookId)
		book.ID = &bkId

		rowsUpdate, err := h.repository.BookWorker.Update(*book)
		if err != nil {
			log.Printf("Error updating book %d: %v", bookId, err)
			http.Error(w, "Error updating book", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		response := fmt.Sprintf("Count updated book: %v", rowsUpdate)
		_, err = w.Write([]byte(response))
		if err != nil {
			log.Println("Error write answer:", err)
		}
	case http.MethodDelete:
		log.Println("Processing method DELETE")
		rowsDel, err := h.repository.BookWorker.Delete(bookId)
		if err != nil {
			log.Printf("Error deleting book %d: %v", bookId, err)
			http.Error(w, "Error deleting book", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		response := fmt.Sprintf("Count deleted book: %v", rowsDel)
		_, err = w.Write([]byte(response))
		if err != nil {
			log.Println("Error write answer:", err)
		}
	default:
		log.Println("Unsupported method")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
