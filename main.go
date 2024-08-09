package main

import (
	"RESTful-APIs-with-Go-and-MySQL-Using-the-Repository-Pattern/config"
	"RESTful-APIs-with-Go-and-MySQL-Using-the-Repository-Pattern/handler"
	"RESTful-APIs-with-Go-and-MySQL-Using-the-Repository-Pattern/repository"
	"RESTful-APIs-with-Go-and-MySQL-Using-the-Repository-Pattern/service"
	"context"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func main() {
	// Initialize logging
	config.InitializeLogger()

	//  Preparing the context
	ctx := context.Background()

	// Loading database connection from config
	db, err := config.LoadDatabase(ctx)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Initialize repositories, services, and handlers
	bookRepository := repository.NewMySQLBookRepository(db)
	bookService := service.NewBookService(bookRepository)
	bookHandler := handler.NewBookHandler(bookService)

	// Initialize the router
	router := gin.Default()

	// Register routes
	router.GET("/books", bookHandler.GetAllBooks)
	router.GET("/books/:id", bookHandler.GetBookByID)
	router.POST("/books", bookHandler.CreateBook)
	router.PUT("/books/:id", bookHandler.UpdateBook)
	router.DELETE("/books/:id", bookHandler.DeleteBook)

	url := fmt.Sprintf("%s:%s", os.Getenv("HOST"), os.Getenv("PORT"))
	log.Info().Msgf("Server running at http://%s/", url)

	router.Run(url)
}
