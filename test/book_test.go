package test

import (
	"RESTful-APIs-with-Go-and-MySQL-Using-the-Repository-Pattern/handler"
	"RESTful-APIs-with-Go-and-MySQL-Using-the-Repository-Pattern/repository"
	"RESTful-APIs-with-Go-and-MySQL-Using-the-Repository-Pattern/service"
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

func connectDatabase(ctx context.Context) (*sql.DB, error) {
	db, err := sql.Open("mysql", "root:example@tcp(mysql:3306)/book_db_test?parseTime=true")
	if err != nil {
		log.Error().Err(err).Msg("Failed to open connection to database")
		return nil, err
	}

	// Optional: Test the database connection with context
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Optional: Test the database connection
	if err := db.PingContext(ctx); err != nil {
		log.Error().Err(err).Msg("Failed to connect to database")
		return nil, err
	}

	log.Info().Msg("Connection to database successful")
	return db, nil
}

func truncateData(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, "TRUNCATE books")
	return err
}

func TestIntegrationCreateBook(t *testing.T) {
	//  Preparing the context
	ctx := context.Background()

	db, err := connectDatabase(ctx)
	if err != nil {
		t.Fatalf("Error setting up database: %v", err)
	}
	defer db.Close()

	// Initialize repositories, services, and handlers
	bookRepository := repository.NewMySQLBookRepository(db)
	bookService := service.NewBookService(bookRepository)
	bookHandler := handler.NewBookHandler(bookService)

	// Initialize the router
	router := gin.Default()
	router.POST("/books", bookHandler.CreateBook)

	// Define test for case Successfully Created Data
	t.Run("Successfully Created Data", func(t *testing.T) {
		err := truncateData(ctx, db)
		if err != nil {
			t.Fatalf("Error clearing books table: %v", err)
		}

		requestBody, _ := json.Marshal(map[string]interface{}{
			"title":  "Go Programming - From Beginner to Professional",
			"author": "Samantha Coyle",
			"year":   2024,
		})

		req := httptest.NewRequest(http.MethodPost, "/books", bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		var response map[string]interface{}
		json.Unmarshal(rec.Body.Bytes(), &response)

		expectedData := response["data"]

		expectedResponse := map[string]interface{}{
			"code":    float64(http.StatusCreated),
			"message": "Successfully created data",
			"data":    expectedData,
		}

		// If you want to show log the result expected data actual data turn on this line below
		// log.Info().Msgf("Expected => %v", expectedResponse)
		// log.Info().Msgf("Actual => %v", response)

		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Equal(t, expectedResponse, response)
	})

	// Define test for case Duplicate Title
	t.Run("Duplicate Title", func(t *testing.T) {
		// Assume there is already a book with this title in the database
		requestBody, _ := json.Marshal(map[string]interface{}{
			"title":  "Go Programming - From Beginner to Professional",
			"author": "Samantha Coyle",
			"year":   2024,
		})

		req := httptest.NewRequest(http.MethodPost, "/books", bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		var response map[string]interface{}
		json.Unmarshal(rec.Body.Bytes(), &response)

		expectedResponse := map[string]interface{}{
			"code":    float64(http.StatusBadRequest),
			"message": "The book with the same title already exists.",
			"errors":  nil,
		}

		// If you want to show log the result expected data actual data turn on this line below
		// log.Info().Msgf("Expected => %v", expectedResponse)
		// log.Info().Msgf("Actual => %v", response)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, expectedResponse, response)
	})

	// Define test for case Validation Error
	t.Run("Validation Error", func(t *testing.T) {
		// Assume all fields are blank
		requestBody, _ := json.Marshal(map[string]interface{}{
			"title":  "",
			"author": "",
			"year":   0,
		})

		req := httptest.NewRequest(http.MethodPost, "/books", bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		var response map[string]interface{}
		json.Unmarshal(rec.Body.Bytes(), &response)

		expectedErrors := response["errors"]

		expectedResponse := map[string]interface{}{
			"code":    float64(http.StatusUnprocessableEntity),
			"message": "validation error",
			"errors":  expectedErrors,
		}

		// If you want to show log the result expected data actual data turn on this line below
		// log.Info().Msgf("Expected => %v", expectedResponse)
		// log.Info().Msgf("Actual => %v", response)

		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
		assert.Equal(t, expectedResponse, response)
	})
}

func TestIntegrationGetAllBooks(t *testing.T) {
	//  Preparing the context
	ctx := context.Background()

	db, err := connectDatabase(ctx)
	if err != nil {
		t.Fatalf("Error setting up database: %v", err)
	}

	// Initialize repositories, services, and handlers
	bookRepository := repository.NewMySQLBookRepository(db)
	bookService := service.NewBookService(bookRepository)
	bookHandler := handler.NewBookHandler(bookService)

	// Initialize the router
	router := gin.Default()
	router.GET("/books", bookHandler.GetAllBooks)

	// Define test for case Successful Get All of Data
	t.Run("Successfully Got All Of data", func(t *testing.T) {

		req := httptest.NewRequest(http.MethodGet, "/books", nil)
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		var response map[string]interface{}
		json.Unmarshal(rec.Body.Bytes(), &response)

		expectedData := response["data"]

		expectedResponse := map[string]interface{}{
			"code":    float64(http.StatusOK),
			"message": "Successfully got all data.",
			"data":    expectedData,
		}

		// If you want to show log the result expected data actual data turn on this line below
		// log.Info().Msgf("Expected => %v", expectedResponse)
		// log.Info().Msgf("Actual => %v", response)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, expectedResponse, response)
	})

}
func TestIntegrationGetBookByID(t *testing.T) {
	//  Preparing the context
	ctx := context.Background()

	db, err := connectDatabase(ctx)
	if err != nil {
		t.Fatalf("Error setting up database: %v", err)
	}

	// Initialize repositories, services, and handlers
	bookRepository := repository.NewMySQLBookRepository(db)
	bookService := service.NewBookService(bookRepository)
	bookHandler := handler.NewBookHandler(bookService)

	// Initialize the router
	router := gin.Default()
	router.GET("/books/:id", bookHandler.GetBookByID)

	// Define test for case Successful Get Data
	t.Run("Successful Get The Data", func(t *testing.T) {
		// Assume this ID is exists in your test database
		bookID := 1

		req := httptest.NewRequest(http.MethodGet, "/books/"+fmt.Sprint(bookID), nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		var response map[string]interface{}
		json.Unmarshal(rec.Body.Bytes(), &response)

		expectedData := response["data"]

		expectedResponse := map[string]interface{}{
			"code":    float64(http.StatusOK),
			"message": "Successfully got the data.",
			"data":    expectedData,
		}

		// If you want to show log the result expected data actual data turn on this line below
		// log.Info().Msgf("Expected => %v", expectedResponse)
		// log.Info().Msgf("Actual => %v", response)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, expectedResponse, response)
	})

	// Define test for case Data Not Found
	t.Run("Data Not Found", func(t *testing.T) {
		// Assume this ID is not exists in your test database
		bookID := 0

		req := httptest.NewRequest(http.MethodGet, "/books/"+fmt.Sprint(bookID), nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		var response map[string]interface{}
		json.Unmarshal(rec.Body.Bytes(), &response)

		expectedResponse := map[string]interface{}{
			"code":    float64(http.StatusNotFound),
			"message": "Data with that ID does not exist.",
			"errors":  nil,
		}

		// If you want to show log the result expected data actual data turn on this line below
		// log.Info().Msgf("Expected => %v", expectedResponse)
		// log.Info().Msgf("Actual => %v", response)

		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.Equal(t, expectedResponse, response)
	})

	// Define test for case ID Invalid
	t.Run("ID Invalid", func(t *testing.T) {
		// Assume this ID is not a number
		bookID := "ff"

		req := httptest.NewRequest(http.MethodGet, "/books/"+fmt.Sprint(bookID), nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		var response map[string]interface{}
		json.Unmarshal(rec.Body.Bytes(), &response)

		expectedResponse := map[string]interface{}{
			"code":    float64(http.StatusBadRequest),
			"message": "ID must be a valid number.",
			"errors":  nil,
		}

		// If you want to show log the result expected data actual data turn on this line below
		// log.Info().Msgf("Expected => %v", expectedResponse)
		// log.Info().Msgf("Actual => %v", response)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, expectedResponse, response)
	})
}

func TestIntegrationUpdateBook(t *testing.T) {
	//  Preparing the context
	ctx := context.Background()

	db, err := connectDatabase(ctx)
	if err != nil {
		t.Fatalf("Error setting up database: %v", err)
	}
	defer db.Close()

	// Initialize repositories, services, and handlers
	bookRepository := repository.NewMySQLBookRepository(db)
	bookService := service.NewBookService(bookRepository)
	bookHandler := handler.NewBookHandler(bookService)

	// Initialize the router
	router := gin.Default()
	router.PUT("/books/:id", bookHandler.UpdateBook)

	// Define test for case Successfully Updated Data
	t.Run("Successfully Updated Data", func(t *testing.T) {
		requestBody, _ := json.Marshal(map[string]interface{}{
			"title":  "Go Programming - From Beginner to Professional 2nd Edition",
			"author": "Samantha Coyle",
			"year":   2024,
		})

		// Assume this ID is exists in your test database
		bookID := 1

		req := httptest.NewRequest(http.MethodPut, "/books/"+fmt.Sprint(bookID), bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		var response map[string]interface{}
		json.Unmarshal(rec.Body.Bytes(), &response)

		expectedData := response["data"]

		expectedResponse := map[string]interface{}{
			"code":    float64(http.StatusOK),
			"message": "Successfully updated data.",
			"data":    expectedData,
		}

		// If you want to show log the result expected data actual data turn on this line below
		// log.Info().Msgf("Expected => %v", expectedResponse)
		// log.Info().Msgf("Actual => %v", response)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, expectedResponse, response)
	})

	// Define test for case Data Not Found
	t.Run("Data Not Found", func(t *testing.T) {
		requestBody, _ := json.Marshal(map[string]interface{}{
			"title":  "Go Programming - From Beginner to Professional 2nd Edition",
			"author": "Samantha Coyle",
			"year":   2024,
		})

		// Assume this ID is not exists in your test database
		bookID := 0

		req := httptest.NewRequest(http.MethodPut, "/books/"+fmt.Sprint(bookID), bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		var response map[string]interface{}
		json.Unmarshal(rec.Body.Bytes(), &response)

		expectedResponse := map[string]interface{}{
			"code":    float64(http.StatusNotFound),
			"message": "Data with that ID does not exist, cannot update.",
			"errors":  nil,
		}

		// If you want to show log the result expected data actual data turn on this line below
		// log.Info().Msgf("Expected => %v", expectedResponse)
		// log.Info().Msgf("Actual => %v", response)

		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.Equal(t, expectedResponse, response)
	})

	// Define test for case ID Invalid
	t.Run("ID Invalid", func(t *testing.T) {
		requestBody, _ := json.Marshal(map[string]interface{}{
			"title":  "Go Programming - From Beginner to Professional 2nd Edition",
			"author": "Samantha Coyle",
			"year":   2024,
		})

		// Assume this ID is not a number
		bookID := "ff"

		req := httptest.NewRequest(http.MethodPut, "/books/"+fmt.Sprint(bookID), bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		var response map[string]interface{}
		json.Unmarshal(rec.Body.Bytes(), &response)

		expectedResponse := map[string]interface{}{
			"code":    float64(http.StatusBadRequest),
			"message": "ID must be a valid number.",
			"errors":  nil,
		}

		// If you want to show log the result expected data actual data turn on this line below
		// log.Info().Msgf("Expected => %v", expectedResponse)
		// log.Info().Msgf("Actual => %v", response)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, expectedResponse, response)
	})

	// Define test for case Validation Error
	t.Run("Validation Error", func(t *testing.T) {
		// Assume all fields are blank
		requestBody, _ := json.Marshal(map[string]interface{}{
			"title":  "",
			"author": "",
			"year":   0,
		})

		// Assume this ID is exists in your test database
		bookID := 1

		req := httptest.NewRequest(http.MethodPut, "/books/"+fmt.Sprint(bookID), bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		var response map[string]interface{}
		json.Unmarshal(rec.Body.Bytes(), &response)

		expectedErrors := response["errors"]

		expectedResponse := map[string]interface{}{
			"code":    float64(http.StatusUnprocessableEntity),
			"message": "validation error",
			"errors":  expectedErrors,
		}

		// If you want to show log the result expected data actual data turn on this line below
		// log.Info().Msgf("Expected => %v", expectedResponse)
		// log.Info().Msgf("Actual => %v", response)

		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
		assert.Equal(t, expectedResponse, response)
	})

	// Define test for case The Year Cannot Be In The Future
	t.Run("The Year Cannot Be In The Future", func(t *testing.T) {
		requestBody, _ := json.Marshal(map[string]interface{}{
			"title":  "Go Programming - From Beginner to Professional 2nd Edition",
			"author": "Samantha Coyle",
			"year":   2045, // Assume that the year field is situated above the current year.
		})

		// Assume this ID is exists in your test database
		bookID := 1

		req := httptest.NewRequest(http.MethodPut, "/books/"+fmt.Sprint(bookID), bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		var response map[string]interface{}
		json.Unmarshal(rec.Body.Bytes(), &response)

		expectedResponse := map[string]interface{}{
			"code":    float64(http.StatusBadRequest),
			"message": "Year of publication cannot be in the future, cannot update.",
			"errors":  nil,
		}

		// If you want to show log the result expected data actual data turn on this line below
		// log.Info().Msgf("Expected => %v", expectedResponse)
		// log.Info().Msgf("Actual => %v", response)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, expectedResponse, response)
	})
}

func TestIntegrationDeleteBook(t *testing.T) {
	//  Preparing the context
	ctx := context.Background()

	db, err := connectDatabase(ctx)
	if err != nil {
		t.Fatalf("Error setting up database: %v", err)
	}
	defer db.Close()

	// Initialize repositories, services, and handlers
	bookRepository := repository.NewMySQLBookRepository(db)
	bookService := service.NewBookService(bookRepository)
	bookHandler := handler.NewBookHandler(bookService)

	// Initialize the router
	router := gin.Default()
	router.DELETE("/books/:id", bookHandler.DeleteBook)

	// Define test for case Successfully Deleted Data
	t.Run("Successfully Deleted Data", func(t *testing.T) {
		// Assume this ID is exists in your test database
		bookID := 1

		req := httptest.NewRequest(http.MethodDelete, "/books/"+fmt.Sprint(bookID), nil)
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		var response map[string]interface{}
		json.Unmarshal(rec.Body.Bytes(), &response)

		expectedResponse := map[string]interface{}{
			"code":    float64(http.StatusOK),
			"message": "Successfully deleted data.",
			"data":    nil,
		}

		// If you want to show log the result expected data actual data turn on this line below
		// log.Info().Msgf("Expected => %v", expectedResponse)
		// log.Info().Msgf("Actual => %v", response)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, expectedResponse, response)
	})

	// Define test for case Data Not Found
	t.Run("Data Not Found", func(t *testing.T) {
		// Assume this ID is not exists in your test database
		bookID := 0

		req := httptest.NewRequest(http.MethodDelete, "/books/"+fmt.Sprint(bookID), nil)
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		var response map[string]interface{}
		json.Unmarshal(rec.Body.Bytes(), &response)

		expectedResponse := map[string]interface{}{
			"code":    float64(http.StatusNotFound),
			"message": "Data with that ID does not exist, you cannot delete it.",
			"errors":  nil,
		}

		// If you want to show log the result expected data actual data turn on this line below
		// log.Info().Msgf("Expected => %v", expectedResponse)
		// log.Info().Msgf("Actual => %v", response)

		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.Equal(t, expectedResponse, response)
	})

	// Define test for case ID Invalid
	t.Run("ID Invalid", func(t *testing.T) {
		// Assume this ID is not a number
		bookID := "ff"

		req := httptest.NewRequest(http.MethodDelete, "/books/"+fmt.Sprint(bookID), nil)
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		var response map[string]interface{}
		json.Unmarshal(rec.Body.Bytes(), &response)

		expectedResponse := map[string]interface{}{
			"code":    float64(http.StatusBadRequest),
			"message": "ID must be a valid number.",
			"errors":  nil,
		}

		// If you want to show log the result expected data actual data turn on this line below
		// log.Info().Msgf("Expected => %v", expectedResponse)
		// log.Info().Msgf("Actual => %v", response)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, expectedResponse, response)
	})

	// Define test for case Books Older Than 10 Years
	t.Run("Books Older Than 10 Years", func(t *testing.T) {
		// Assume this ID is exists in your test database
		bookID := 2

		req := httptest.NewRequest(http.MethodDelete, "/books/"+fmt.Sprint(bookID), nil)
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		var response map[string]interface{}
		json.Unmarshal(rec.Body.Bytes(), &response)

		expectedResponse := map[string]interface{}{
			"code":    float64(http.StatusBadRequest),
			"message": "Books older than 10 years cannot be deleted.",
			"errors":  nil,
		}

		// If you want to show log the result expected data actual data turn on this line below
		// log.Info().Msgf("Expected => %v", expectedResponse)
		// log.Info().Msgf("Actual => %v", response)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, expectedResponse, response)
	})
}
