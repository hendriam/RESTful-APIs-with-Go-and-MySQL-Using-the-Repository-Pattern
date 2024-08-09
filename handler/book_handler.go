package handler

import (
	"RESTful-APIs-with-Go-and-MySQL-Using-the-Repository-Pattern/helper"
	"RESTful-APIs-with-Go-and-MySQL-Using-the-Repository-Pattern/service"
	"net/http"
	"strconv"

	"RESTful-APIs-with-Go-and-MySQL-Using-the-Repository-Pattern/models"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type BookHandler struct {
	Service service.BookService
}

func NewBookHandler(service service.BookService) *BookHandler {
	return &BookHandler{Service: service}
}

func (h *BookHandler) GetAllBooks(c *gin.Context) {
	books, err := h.Service.GetAllBooks(c.Request.Context())
	if err != nil {
		log.Error().Err(err).Msg("[BookHandler] Failed to get data")
		helper.SendErrorResponse(c, http.StatusInternalServerError, "Failed to get data", nil)
		return
	}

	helper.SendSuccessResponse(c, http.StatusOK, "Successfully got all data", books)
	log.Info().Msg("[BookHandler] Successfully got all data")
}

func (h *BookHandler) GetBookByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Error().Err(err).Msg("[BookHandler] Failed to convert ID from URL")
		helper.SendErrorResponse(c, http.StatusBadRequest, "ID must be a valid number", nil)
		return
	}

	book, err := h.Service.GetBookByID(c.Request.Context(), id)
	if err != nil {
		log.Error().Err(err).Int("id", id).Msg("[BookHandler] Failed to get data")
		if err.Error() == "errBookNotFound" {
			helper.SendErrorResponse(c, http.StatusNotFound, "Data with that ID does not exist", nil)
			return
		}
		helper.SendErrorResponse(c, http.StatusInternalServerError, "Failed to get data", nil)
		return
	}

	helper.SendSuccessResponse(c, http.StatusOK, "Successfully got the data", book)
	log.Info().Int("id", id).Msg("[BookHandler] Successfully got the data")
}

func (h *BookHandler) CreateBook(c *gin.Context) {
	var book models.Book
	if err := c.ShouldBindJSON(&book); err != nil {
		log.Error().Err(err).Msg("Failed to process input data")
		helper.HandleValidationError(c, err)
		return
	}

	if err := h.Service.CreateBook(c.Request.Context(), &book); err != nil {
		log.Error().Err(err).Msg("[BookHandler] Failed to create data")
		if err.Error() == "ErrBookExists" {
			helper.SendErrorResponse(c, http.StatusBadRequest, "The book with the same title already exists", nil)
			return
		}
		helper.SendErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	helper.SendSuccessResponse(c, http.StatusCreated, "Successfully created data", book)
	log.Info().Int("id", book.ID).Msg("[BookHandler] Successfully created data")
}

func (h *BookHandler) UpdateBook(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Error().Err(err).Msg("[BookHandler] Failed to convert ID from URL.")
		helper.SendErrorResponse(c, http.StatusBadRequest, "ID must be a valid number.", nil)
		return
	}

	var book models.Book
	if err := c.ShouldBindJSON(&book); err != nil {
		log.Error().Err(err).Msg("[BookHandler] Failed to process data update.")
		helper.HandleValidationError(c, err)
		return
	}

	book.ID = id

	if err := h.Service.UpdateBook(c.Request.Context(), &book); err != nil {
		log.Error().Err(err).Int("id", id).Msg("[BookHandler] Failed to update data.")
		switch err.Error() {
		case "errBookNotFound":
			helper.SendErrorResponse(c, http.StatusNotFound, "Data with that ID does not exist, cannot update", nil)
			return
		case "errTheYearCannotBeInTheFuture":
			helper.SendErrorResponse(c, http.StatusBadRequest, "Year of publication cannot be in the future, cannot update", nil)
			return
		}
		helper.SendErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	helper.SendSuccessResponse(c, http.StatusOK, "Successfully updated data.", book)
	log.Info().Int("id", book.ID).Msg("[BookHandler] Successfully updated data.")
}

func (h *BookHandler) DeleteBook(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Error().Err(err).Msg("[BookHandler] Failed to convert ID from URL.")
		helper.SendErrorResponse(c, http.StatusBadRequest, "ID must be a valid number.", nil)
		return
	}

	if err := h.Service.DeleteBook(c.Request.Context(), id); err != nil {
		log.Error().Err(err).Int("id", id).Msg("[BookHandler] Failed to delete data.")
		switch err.Error() {
		case "errBookNotFound":
			helper.SendErrorResponse(c, http.StatusNotFound, "Data with that ID does not exist, cannot deleted", nil)
			return
		case "errBooksOlderThan10Years":
			helper.SendErrorResponse(c, http.StatusBadRequest, "Books older than 10 years cannot be deleted", nil)
			return
		}
		helper.SendErrorResponse(c, http.StatusInternalServerError, "Failed to delete data.", nil)
		return
	}

	helper.SendSuccessResponse(c, http.StatusOK, "Successfully deleted data.", nil)
	log.Info().Int("id", id).Msg("[BookHandler] Successfully deleted data.")
}
