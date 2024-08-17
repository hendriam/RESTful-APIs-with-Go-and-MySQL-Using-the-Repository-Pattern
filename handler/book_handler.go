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

	helper.SendSuccessResponse(c, http.StatusOK, "Successfully got all data.", books)
	log.Info().Msg("[BookHandler] Successfully got all data.")
}

func (h *BookHandler) GetBookByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Error().Err(err).Msg("[BookHandler] Failed to convert ID from URL")
		helper.SendErrorResponse(c, http.StatusBadRequest, "ID must be a valid number.", nil)
		return
	}

	book, err := h.Service.GetBookByID(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "errBookNotFound" {
			log.Error().Err(err).Int("id", id).Msg("[BookHandler] Data with that ID does not exist.")
			helper.SendErrorResponse(c, http.StatusNotFound, "Data with that ID does not exist.", nil)
			return
		}
		log.Error().Err(err).Int("id", id).Msg("[BookHandler] Failed to get data")
		helper.SendErrorResponse(c, http.StatusInternalServerError, "Failed to get data", nil)
		return
	}

	helper.SendSuccessResponse(c, http.StatusOK, "Successfully got the data.", book)
	log.Info().Int("id", id).Msg("[BookHandler] Successfully got the data.")
}

func (h *BookHandler) CreateBook(c *gin.Context) {
	var book models.Book
	if err := c.ShouldBindJSON(&book); err != nil {
		log.Error().Err(err).Msg("Failed to process input data")
		helper.HandleValidationError(c, err)
		return
	}

	if err := h.Service.CreateBook(c.Request.Context(), &book); err != nil {
		if err.Error() == "ErrBookExists" {
			log.Error().Err(err).Msg("[BookHandler] The book with the same title already exists.")
			helper.SendErrorResponse(c, http.StatusBadRequest, "The book with the same title already exists.", nil)
			return
		}
		log.Error().Err(err).Msg("[BookHandler] Failed to create data")
		helper.SendErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	helper.SendSuccessResponse(c, http.StatusCreated, "Successfully created data", book)
	log.Info().Int("id", book.ID).Msgf("[BookHandler] Successfully created data %v", book)
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
		switch err.Error() {
		case "errBookNotFound":
			log.Error().Err(err).Int("id", id).Msg("[BookHandler] Data with that ID does not exist, cannot update.")
			helper.SendErrorResponse(c, http.StatusNotFound, "Data with that ID does not exist, cannot update.", nil)
			return
		case "errTheYearCannotBeInTheFuture":
			log.Error().Err(err).Int("id", id).Msg("[BookHandler] Year of publication cannot be in the future, cannot update.")
			helper.SendErrorResponse(c, http.StatusBadRequest, "Year of publication cannot be in the future, cannot update.", nil)
			return
		}
		log.Error().Err(err).Int("id", id).Msg("[BookHandler] Failed to update data.")
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
		switch err.Error() {
		case "errBookNotFound":
			log.Error().Err(err).Int("id", id).Msg("[BookHandler] Data with that ID does not exist, you cannot delete it.")
			helper.SendErrorResponse(c, http.StatusNotFound, "Data with that ID does not exist, you cannot delete it.", nil)
			return
		case "errBooksOlderThan10Years":
			log.Error().Err(err).Int("id", id).Msg("[BookHandler] Books older than 10 years cannot be deleted.")
			helper.SendErrorResponse(c, http.StatusBadRequest, "Books older than 10 years cannot be deleted.", nil)
			return
		}
		log.Error().Err(err).Int("id", id).Msg("[BookHandler] Failed to delete data.")
		helper.SendErrorResponse(c, http.StatusInternalServerError, "Failed to delete data.", nil)
		return
	}

	helper.SendSuccessResponse(c, http.StatusOK, "Successfully deleted data.", nil)
	log.Info().Int("id", id).Msg("[BookHandler] Successfully deleted data.")
}
