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
	mux.HandleFunc("/books/{book_id}/authors/{author_id}", h.updateBookAndAuthor)
	return mux
}

func (h *Hundler) books(w http.ResponseWriter, r *http.Request) {
	log.Printf("Call link %s, method: %s", r.RequestURI, r.Method)

	switch r.Method {
	case http.MethodGet:
		log.Println("Processing method GET")
		books, err := h.repository.GetAllBook()
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
			errorResult(&w, "Error decoding input JSON", err, http.StatusBadRequest)
			return
		}
		log.Println("Input body:", string(body))
		//err := json.NewDecoder(r.Body).Decode(book)
		err = json.Unmarshal(body, book)
		if err != nil {
			errorResult(&w, "Error decoding input JSON", err, http.StatusBadRequest)
			return
		}
		log.Printf("Getting book: %+v\n", *book)

		if err = book.Validate(false); err != nil {
			errorResult(&w, err.Error(), err, http.StatusBadRequest)
			return
		}

		id, err := h.repository.CreateBook(book)
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
		book, err := h.repository.GetBookById(id)
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
			errorResult(&w, "Error decoding input JSON", err, http.StatusBadRequest)
			return
		}
		log.Println("Input body:", string(body))
		//err := json.NewDecoder(r.Body).Decode(book)
		err = json.Unmarshal(body, book)
		if err != nil {
			errorResult(&w, "Error decoding input JSON", err, http.StatusBadRequest)
			return
		}
		log.Printf("Getting book: %+v\n", *book)

		if err = book.Validate(false); err != nil {
			errorResult(&w, err.Error(), err, http.StatusBadRequest)
			return
		}

		id32 := int32(id)
		if book.ID != nil && *book.ID != id32 {
			errorResult(&w, "Book id does not match with URL", nil, http.StatusBadRequest)
			return
		}
		book.ID = &id32

		updResult, err := h.repository.UpdateBook(book)
		if err != nil {
			errorResult(&w, "Error updating book", err, http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		response := "Updated book successfully"
		if !updResult {
			response = "No books updated"
		}
		_, err = w.Write([]byte(response))
		if err != nil {
			log.Println("Error write answer:", err)
		}
	case http.MethodDelete:
		log.Println("Processing method DELETE")
		delResult, err := h.repository.DeleteBook(id)
		if err != nil {
			errorResult(&w, "Error deleting book", err, http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		response := "Deleted book successfully"
		if !delResult {
			response = "No books deleted"
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
		authors, err := h.repository.GetAllAuthor()
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
			errorResult(&w, "Error decoding input JSON", err, http.StatusBadRequest)
			return
		}
		log.Println("Input body:", string(body))
		//err = json.NewDecoder(r.Body).Decode(author)
		err = json.Unmarshal(body, author)
		if err != nil {
			errorResult(&w, "Error decoding input JSON", err, http.StatusBadRequest)
			return
		}
		log.Printf("Getting author: %+v\n", *author)

		if err = author.Validate(); err != nil {
			errorResult(&w, err.Error(), err, http.StatusBadRequest)
			return
		}

		id, err := h.repository.CreateAuthor(author)
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
		author, err := h.repository.GetAuthorById(id)
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
			errorResult(&w, "Error decoding input JSON", err, http.StatusBadRequest)
			return
		}
		log.Println("Input body:", string(body))
		//err := json.NewDecoder(r.Body).Decode(author)
		err = json.Unmarshal(body, author)
		if err != nil {
			errorResult(&w, "Error decoding input JSON", err, http.StatusBadRequest)
			return
		}
		log.Printf("Getting author: %+v\n", *author)

		if err = author.Validate(); err != nil {
			errorResult(&w, err.Error(), err, http.StatusBadRequest)
			return
		}

		id32 := int32(id)
		if author.ID != nil && *author.ID != id32 {
			errorResult(&w, "Author id does not match with URL", nil, http.StatusBadRequest)
			return
		}
		author.ID = &id32

		updResult, err := h.repository.UpdateAuthor(author)
		if err != nil {
			errorResult(&w, "Error updating author", err, http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		response := "Updated author successfully"
		if !updResult {
			response = "No authors updated"
		}
		_, err = w.Write([]byte(response))
		if err != nil {
			log.Println("Error write answer:", err)
		}
	case http.MethodDelete:
		log.Println("Processing method DELETE")
		delResult, err := h.repository.DeleteAuthor(id)
		if err != nil {
			errorResult(&w, "Error deleting author", err, http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		response := "Deleted author successfully"
		if !delResult {
			response = "No authors deleted"
		}
		_, err = w.Write([]byte(response))
		if err != nil {
			log.Println("Error write answer:", err)
		}
	default:
		errorResult(&w, "Method not allowed", nil, http.StatusMethodNotAllowed)
	}
}

func (h *Hundler) updateBookAndAuthor(w http.ResponseWriter, r *http.Request) {
	log.Printf("Call link %s, method: %s", r.RequestURI, r.Method)
	if r.Method != http.MethodPut {
		errorResult(&w, "Method not allowed", nil, http.StatusMethodNotAllowed)
		return
	}

	param1 := r.PathValue("book_id")
	bookId, err := strconv.Atoi(param1)
	if err != nil {
		log.Printf("Error converting %s to int: %v\n", param1, err)
		http.Error(w, "Bad book id in URL", http.StatusBadRequest)
		return
	}

	param2 := r.PathValue("author_id")
	authorId, err := strconv.Atoi(param2)
	if err != nil {
		log.Printf("Error converting %s to int: %v\n", param2, err)
		http.Error(w, "Bad author id in URL", http.StatusBadRequest)
		return
	}

	book := &repository.Book{}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		errorResult(&w, "Error decoding input JSON", err, http.StatusBadRequest)
		return
	}
	log.Println("Input body:", string(body))
	//err := json.NewDecoder(r.Body).Decode(book)
	err = json.Unmarshal(body, book)
	if err != nil {
		errorResult(&w, "Error decoding input JSON", err, http.StatusBadRequest)
		return
	}
	log.Printf("Getting book: %+v\n", book)

	if err = book.Validate(true); err != nil {
		errorResult(&w, err.Error(), err, http.StatusBadRequest)
		return
	}

	bookId32 := int32(bookId)
	if book.ID != nil && *book.ID != bookId32 {
		errorResult(&w, "Book id does not match with URL", nil, http.StatusBadRequest)
		return
	}
	book.ID = &bookId32

	authorId32 := int32(authorId)
	if book.Author.ID != nil && *book.Author.ID != authorId32 {
		errorResult(&w, "Author id does not match with URL", nil, http.StatusBadRequest)
		return
	}
	book.Author.ID = &authorId32

	updResult, err := h.repository.UpdateBookAndAuthor(book)
	if err != nil {
		errorResult(&w, "Error updating", err, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	response := "Updated book and author successfully"
	if !updResult {
		response = "No objects updated"
	}
	_, err = w.Write([]byte(response))
	if err != nil {
		log.Println("Error write answer:", err)
	}
}

func errorResult(w *http.ResponseWriter, msg string, err error, stat int) {
	log.Println(msg, ". Error:", err)
	http.Error(*w, msg, stat)
}
