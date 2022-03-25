package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/bestetufan/bookstore/domain/book"
	"github.com/bestetufan/bookstore/models"
	"github.com/gorilla/mux"
)

type BookController struct{}

var Book BookController

// CONSTS
const (
	SUCCESS_MESSAGE       = "Operation completed successfully"
	STOCK_COUNT_MESSAGE   = "Not enough stock. Available stock: %d"
	COUNT_PARAM_MESSAGE   = "Count must be greater than zero"
	INVALID_MODEL_MESSAGE = "Invalid model"
)

// UTILS
func respondWithMessage(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"message": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func (c *BookController) List(bookRepo *book.BookRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		books, err := bookRepo.GetAllBooks()
		if err != nil {
			respondWithMessage(w, http.StatusInternalServerError, err.Error())
			return
		}

		respondWithJSON(w, http.StatusOK, books)
	}
}

func (c *BookController) Get(bookRepo *book.BookRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id := params["id"]

		bookId, err := strconv.Atoi(id)
		if err != nil {
			respondWithMessage(w, http.StatusBadRequest, err.Error())
			return
		}

		book, err := bookRepo.GetBookById(bookId)
		if err != nil {
			respondWithMessage(w, http.StatusInternalServerError, err.Error())
			return
		}

		respondWithJSON(w, http.StatusOK, book)
	}
}

func (c *BookController) Create(bookRepo *book.BookRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var book book.Book
		if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
			respondWithMessage(w, http.StatusBadRequest, INVALID_MODEL_MESSAGE)
			return
		}

		newBook, err := bookRepo.CreateBook(&book)
		if err != nil {
			respondWithMessage(w, http.StatusInternalServerError, err.Error())
			return
		}

		respondWithJSON(w, http.StatusCreated, newBook)
	}
}

func (c *BookController) Update(bookRepo *book.BookRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id := params["id"]

		bookId, err := strconv.Atoi(id)
		if err != nil {
			respondWithMessage(w, http.StatusBadRequest, err.Error())
			return
		}

		var book book.Book
		if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
			respondWithMessage(w, http.StatusBadRequest, INVALID_MODEL_MESSAGE)
			return
		}

		book.ID = uint(bookId)
		updatedBook, err := bookRepo.UpdateBook(&book)
		if err != nil {
			respondWithMessage(w, http.StatusInternalServerError, err.Error())
			return
		}

		respondWithJSON(w, http.StatusOK, updatedBook)
	}
}

func (c *BookController) Delete(bookRepo *book.BookRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id := params["id"]

		bookId, err := strconv.Atoi(id)
		if err != nil {
			respondWithMessage(w, http.StatusBadRequest, err.Error())
			return
		}

		err = bookRepo.DeleteBookById(bookId)
		if err != nil {
			respondWithMessage(w, http.StatusInternalServerError, err.Error())
			return
		}

		respondWithMessage(w, http.StatusOK, SUCCESS_MESSAGE)
	}
}

func (c *BookController) Search(bookRepo *book.BookRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		query := params["query"]

		books := bookRepo.FindBooksByQuery(query)

		respondWithJSON(w, http.StatusOK, books)
	}
}

func (c *BookController) Order(bookRepo *book.BookRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var order models.Order
		if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
			respondWithMessage(w, http.StatusBadRequest, INVALID_MODEL_MESSAGE)
			return
		}

		if order.Count <= 0 {
			respondWithMessage(w, http.StatusBadRequest, COUNT_PARAM_MESSAGE)
			return
		}

		book, err := bookRepo.GetBookById(order.BookId)
		if err != nil {
			respondWithMessage(w, http.StatusInternalServerError, err.Error())
			return
		}

		if order.Count > book.StockCount {
			respondWithMessage(w, http.StatusOK, fmt.Sprintf(STOCK_COUNT_MESSAGE, book.StockCount))
			return
		}

		book.StockCount -= order.Count
		if _, err := bookRepo.UpdateBook(book); err != nil {
			respondWithMessage(w, http.StatusInternalServerError, err.Error())
			return
		}

		respondWithMessage(w, http.StatusOK, SUCCESS_MESSAGE)
	}
}
