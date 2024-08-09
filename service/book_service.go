package service

import (
	"RESTful-APIs-with-Go-and-MySQL-Using-the-Repository-Pattern/models"
	"RESTful-APIs-with-Go-and-MySQL-Using-the-Repository-Pattern/repository"
	"context"
	"errors"
	"time"
)

type BookService interface {
	GetAllBooks(ctx context.Context) ([]models.Book, error)
	GetBookByID(ctx context.Context, id int) (*models.Book, error)
	CreateBook(ctx context.Context, book *models.Book) error
	UpdateBook(ctx context.Context, book *models.Book) error
	DeleteBook(ctx context.Context, id int) error
}

type bookService struct {
	repo repository.BookRepository
}

func NewBookService(repo repository.BookRepository) BookService {
	return &bookService{repo: repo}
}

func (s *bookService) GetAllBooks(ctx context.Context) ([]models.Book, error) {
	return s.repo.GetAllBooks(ctx)
}

func (s *bookService) GetBookByID(ctx context.Context, id int) (*models.Book, error) {
	// Check if a book with that ID exists
	book, err := s.repo.GetBookByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if book == nil {
		return nil, errors.New("errBookNotFound")
	}

	return book, nil
}

func (s *bookService) CreateBook(ctx context.Context, book *models.Book) error {
	// Check if a book with the same title already exists
	existingBooks, err := s.repo.FindByTitle(ctx, book.Title)
	if err != nil {
		return err
	}
	if len(existingBooks) > 0 {
		return errors.New("ErrBookExists")
	}

	// Set default values
	book.CreatedAt = time.Now()
	book.UpdatedAt = time.Now()

	return s.repo.CreateBook(ctx, book)
}

func (s *bookService) UpdateBook(ctx context.Context, book *models.Book) error {
	// Check if a book with that ID exists
	existingBook, err := s.repo.GetBookByID(ctx, book.ID)
	if err != nil {
		return err
	}

	if existingBook == nil {
		return errors.New("errBookNotFound")
	}

	book.CreatedAt = existingBook.CreatedAt
	book.UpdatedAt = time.Now()

	// Validation that the year cannot be in the future
	if book.Year > time.Now().Year() {
		return errors.New("errTheYearCannotBeInTheFuture")
	}

	return s.repo.UpdateBook(ctx, book)
}

func (s *bookService) DeleteBook(ctx context.Context, id int) error {
	// Check if a book with that ID exists
	existingBook, err := s.repo.GetBookByID(ctx, id)
	if err != nil {
		return err
	}

	if existingBook == nil {
		return errors.New("errBookNotFound")
	}

	// For example, books older than 10 years should not be deleted
	if time.Now().Year()-existingBook.Year > 10 {
		return errors.New("errBooksOlderThan10Years")
	}

	return s.repo.DeleteBook(ctx, id)
}
