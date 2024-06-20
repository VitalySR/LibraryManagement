package handler

import (
	"encoding/json"
	"fmt"
	"io"
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
	mux.HandleFunc("/authors", h.authors)
	mux.HandleFunc("/authors/{id}", h.authorsId)
	return mux
}

func (h *Hundler) books(w http.ResponseWriter, r *http.Request) {
	log.Printf("Call link %s, method: %s", r.RequestURI, r.Method)

	switch r.Method {
	case http.MethodGet:
		log.Println("Processing method GET")
		books, err := h.repository.BookWorker.GetAll()
		if err != nil {
			errorResult(&w, "Error getting books", err, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(books)
		if err != nil {
			errorResult(&w, "Error encoding books", err, http.StatusInternalServerError)
			return
		}
	case http.MethodPost:
		book := &repository.Book{}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			errorResult(&w, "Error reading body", err, http.StatusInternalServerError)
		}
		log.Println("Input body:", string(body))
		//err := json.NewDecoder(r.Body).Decode(book)
		err = json.Unmarshal(body, book)
		if err != nil {
			errorResult(&w, "Error decoding input JSON", err, http.StatusBadRequest)
			return
		}
		log.Printf("Getting book: %+v\n", book)

		if book.Title == nil || *book.Title == "" {
			http.Error(w, "title is mandatory field", http.StatusNotAcceptable)
			return
		}

		id, err := h.repository.BookWorker.Create(*book)
		if err != nil {
			errorResult(&w, "Error creating book", err, http.StatusInternalServerError)
			return
		}
		log.Printf("Created book with id %d", id)

		w.WriteHeader(http.StatusCreated)
		_, err = fmt.Fprint(w, "Book created successfully with id = ", id)
		if err != nil {
			log.Println("Error write answer:", err)
		}
	default:
		errorResult(&w, "Method not allowed", nil, http.StatusMethodNotAllowed)
	}
}

func (h *Hundler) bookId(w http.ResponseWriter, r *http.Request) {
	log.Printf("Call link %s, method: %s", r.RequestURI, r.Method)

	param := r.PathValue("id")
	id, err := strconv.Atoi(param)
	if err != nil {
		log.Printf("Error converting %s to int: %v\n", param, err)
		http.Error(w, "Bad book id in URL", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		log.Println("Processing method GET")
		book, err := h.repository.BookWorker.GetById(id)
		log.Println(book, err)
		if err != nil {
			errorResult(&w, "Error getting book", err, http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(book)
		if err != nil {
			errorResult(&w, "Error encoding book", err, http.StatusInternalServerError)
			return
		}
		return
	case http.MethodPut:
		log.Println("Processing method PUT")
		book := &repository.Book{}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			errorResult(&w, "Error reading body", err, http.StatusInternalServerError)
		}
		log.Println("Input body:", string(body))
		//err := json.NewDecoder(r.Body).Decode(book)
		err = json.Unmarshal(body, book)
		if err != nil {
			errorResult(&w, "Error decoding input JSON", err, http.StatusBadRequest)
			return
		}
		log.Printf("Getting book: %+v\n", book)

		if book.Title == nil || *book.Title == "" {
			http.Error(w, "title is mandatory field", http.StatusNotAcceptable)
			return
		}
		bkId := int32(id)
		book.ID = &bkId

		rowsUpdate, err := h.repository.BookWorker.Update(*book)
		if err != nil {
			errorResult(&w, "Error updating book", err, http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		response := "Updated book successfully"
		if rowsUpdate == 0 {
			response = "No rows updated"
		}
		_, err = w.Write([]byte(response))
		if err != nil {
			log.Println("Error write answer:", err)
		}
	case http.MethodDelete:
		log.Println("Processing method DELETE")
		rowsDel, err := h.repository.BookWorker.Delete(id)
		if err != nil {
			errorResult(&w, "Error deleting book", err, http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		response := "Deleted book successfully"
		if rowsDel == 0 {
			response = "No rows deleted"
		}
		_, err = w.Write([]byte(response))
		if err != nil {
			log.Println("Error write answer:", err)
		}
	default:
		errorResult(&w, "Method not allowed", nil, http.StatusMethodNotAllowed)
	}
}

func (h *Hundler) authors(w http.ResponseWriter, r *http.Request) {
	log.Printf("Call link %s, method: %s", r.RequestURI, r.Method)

	switch r.Method {
	case http.MethodGet:
		log.Println("Processing method GET")
		authors, err := h.repository.AuthorWorker.GetAll()
		if err != nil {
			errorResult(&w, "Error getting authors", err, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(authors)
		if err != nil {
			errorResult(&w, "Error encoding authors to JSON", err, http.StatusInternalServerError)
			return
		}
	case http.MethodPost:
		log.Println("Processing method POST")
		author := &repository.Author{}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			errorResult(&w, "Error reading body", err, http.StatusInternalServerError)
		}
		log.Println("Input body:", string(body))
		//err = json.NewDecoder(r.Body).Decode(author)
		err = json.Unmarshal(body, author)
		if err != nil {
			errorResult(&w, "Error decoding input JSON", err, http.StatusBadRequest)
			return
		}
		log.Printf("Getting author: %+v\n", author)

		if author.FirstName == nil || *author.FirstName == "" || author.LastName == nil || *author.LastName == "" {
			http.Error(w, "first_name and last_name are mandatory fields", http.StatusNotAcceptable)
			return
		}

		id, err := h.repository.AuthorWorker.Create(*author)
		if err != nil {
			errorResult(&w, "Error creating author", err, http.StatusInternalServerError)
			return
		}
		log.Println("Created author with id =", id)

		w.WriteHeader(http.StatusCreated)
		_, err = fmt.Fprint(w, "Author created successfully with id = ", id)
		if err != nil {
			log.Println("Error write answer:", err)
		}
	default:
		errorResult(&w, "Method not allowed", nil, http.StatusMethodNotAllowed)
	}
}

func (h *Hundler) authorsId(w http.ResponseWriter, r *http.Request) {
	log.Printf("Call link %s, method: %s", r.RequestURI, r.Method)

	param := r.PathValue("id")
	id, err := strconv.Atoi(param)
	if err != nil {
		log.Printf("Error converting %s to int: %v\n", param, err)
		http.Error(w, "Bad author id in URL", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		log.Println("Processing method GET")
		author, err := h.repository.AuthorWorker.GetById(id)
		if err != nil {
			errorResult(&w, "Error getting author", err, http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(author)
		if err != nil {
			errorResult(&w, "Error encoding author to JSON", err, http.StatusInternalServerError)
			return
		}
		return
	case http.MethodPut:
		log.Println("Processing method PUT")
		author := &repository.Author{}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			errorResult(&w, "Error reading body", err, http.StatusInternalServerError)
		}
		log.Println("Input body:", string(body))
		//err := json.NewDecoder(r.Body).Decode(author)
		err = json.Unmarshal(body, author)
		if err != nil {
			errorResult(&w, "Error decoding input JSON", err, http.StatusBadRequest)
			return
		}
		log.Printf("Getting book: %+v\n", author)

		if author.FirstName == nil || *author.FirstName == "" || author.LastName == nil || *author.LastName == "" {
			http.Error(w, "first_name and last_name are mandatory fields", http.StatusNotAcceptable)
			return
		}
		id32 := int32(id)
		author.ID = &id32

		rowsUpdate, err := h.repository.AuthorWorker.Update(*author)
		if err != nil {
			errorResult(&w, "Error updating author", err, http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		response := "Updated author successfully"
		if rowsUpdate == 0 {
			response = "No rows updated"
		}
		_, err = w.Write([]byte(response))
		if err != nil {
			log.Println("Error write answer:", err)
		}
	case http.MethodDelete:
		log.Println("Processing method DELETE")
		rowsDel, err := h.repository.AuthorWorker.Delete(id)
		if err != nil {
			errorResult(&w, "Error deleting author", err, http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		response := "Deleted author successfully"
		if rowsDel == 0 {
			response = "No rows deleted"
		}
		_, err = w.Write([]byte(response))
		if err != nil {
			log.Println("Error write answer:", err)
		}
	default:
		errorResult(&w, "Method not allowed", nil, http.StatusMethodNotAllowed)
	}
}

func errorResult(w *http.ResponseWriter, msg string, err error, stat int) {
	log.Println(msg, ". Error:", err)
	http.Error(*w, msg, stat)
}
