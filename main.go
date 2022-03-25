package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/bestetufan/bookstore/controllers"
	"github.com/bestetufan/bookstore/domain/author"
	"github.com/bestetufan/bookstore/domain/book"
	"github.com/bestetufan/bookstore/helpers"
	"github.com/bestetufan/bookstore/infrastructure"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

var (
	authorRepo *author.AuthorRepository
	bookRepo   *book.BookRepository
)

func init() {
	db := infrastructure.NewPostgresDB("host=localhost user=postgres password=postgres dbname=bookstore port=5432 sslmode=disable")
	authorRepo = author.NewAuthorRepository(db)
	bookRepo = book.NewBookRepository(db)

	// Create database tables
	authorRepo.Migration()
	bookRepo.Migration()

	// Read csv file and generate bookstore slice
	// !! RUNS EVERYTIME, SHOULD COMMENT OUT AFTER FIRST TIME !!
	books, err := helpers.ReadBookCSV("book-data.csv")
	if err != nil {
		fmt.Println("Unable to read csv data!")
		return
	}

	// Feed csv data to db
	bookRepo.InsertSampleData(books)
}

// MWs
func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth_header := r.Header.Get("Authorization")
		if !strings.HasPrefix(auth_header, "Bearer") {
			http.Error(w, "Not Authorized", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(auth_header, "Bearer ")
		if tokenString != "T1E2S3T4" {
			http.Error(w, "Invalid Token", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("URL: %s, METHOD: %s, HEADER: %s", r.RequestURI, r.Method, r.Header)
		next.ServeHTTP(w, r)
	})
}

func main() {
	router := mux.NewRouter()

	// Configure Auth (for all endpoints)
	router.Use(authMiddleware)
	router.Use(loggingMiddleware)

	// Configure routes.
	booksR := router.PathPrefix("/books").Subrouter()
	booksR.HandleFunc("", controllers.Book.List(bookRepo)).Methods(http.MethodGet)
	booksR.HandleFunc("/search/{query}", controllers.Book.Search(bookRepo)).Methods(http.MethodGet)
	booksR.HandleFunc("/order", controllers.Book.Order(bookRepo)).Methods(http.MethodPost)

	bookR := router.PathPrefix("/book").Subrouter()
	bookR.HandleFunc("/{id:[0-9]+}", controllers.Book.Get(bookRepo)).Methods(http.MethodGet)
	bookR.HandleFunc("", controllers.Book.Create(bookRepo)).Methods(http.MethodPost)
	bookR.HandleFunc("/{id:[0-9]+}", controllers.Book.Update(bookRepo)).Methods(http.MethodPut)
	bookR.HandleFunc("/{id:[0-9]+}", controllers.Book.Delete(bookRepo)).Methods(http.MethodDelete)

	// Configure CORS.
	router.Use(handlers.CORS(
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"}),
		handlers.AllowedHeaders([]string{"Accept", "Accept-Language", "Content-Type", "Content-Language", "Origin"}),
		handlers.AllowedOrigins([]string{"*"}),
	))

	// Run our server.
	log.Fatal(http.ListenAndServe(":9000", (router)))
}
