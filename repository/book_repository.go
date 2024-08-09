package repository

import (
	"RESTful-APIs-with-Go-and-MySQL-Using-the-Repository-Pattern/models"
	"context"
	"database/sql"

	"github.com/rs/zerolog/log"
)

type BookRepository interface {
	GetAllBooks(ctx context.Context) ([]models.Book, error)
	GetBookByID(ctx context.Context, id int) (*models.Book, error)
	CreateBook(ctx context.Context, book *models.Book) error
	UpdateBook(ctx context.Context, book *models.Book) error
	DeleteBook(ctx context.Context, id int) error
	FindByTitle(ctx context.Context, title string) ([]models.Book, error)
}

type mysqlBookRepository struct {
	DB *sql.DB
}

func NewMySQLBookRepository(db *sql.DB) BookRepository {
	return &mysqlBookRepository{DB: db}
}

func (r *mysqlBookRepository) GetAllBooks(ctx context.Context) ([]models.Book, error) {
	query := "SELECT * FROM books"
	rows, err := r.DB.QueryContext(ctx, query)
	if err != nil {
		log.Error().Err(err).Msg("[BookRepository] Failed to get all books from database")
		return nil, err
	}
	defer rows.Close()

	var books []models.Book
	for rows.Next() {
		var book models.Book
		if err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Year, &book.CreatedAt, &book.UpdatedAt); err != nil {
			log.Error().Err(err).Msg("[BookRepository] Failed to read book data from query results")
			return nil, err
		}
		books = append(books, book)
	}

	log.Info().Msg("[BookRepository] Successfully got all books from database")
	return books, nil
}

func (r *mysqlBookRepository) GetBookByID(ctx context.Context, id int) (*models.Book, error) {
	var book models.Book
	query := "SELECT * FROM books WHERE id = ?"
	err := r.DB.QueryRowContext(ctx, query, id).
		Scan(&book.ID, &book.Title, &book.Author, &book.Year, &book.CreatedAt, &book.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Warn().Int("id", id).Msg("[BookRepository] Data not found")
			return nil, nil
		}
		log.Error().Err(err).Int("id", id).Msg("[BookRepository] Failed to get data from database")
		return nil, err
	}

	log.Info().Int("id", id).Msg("[BookRepository] Successfully get data from database")
	return &book, nil
}

func (r *mysqlBookRepository) CreateBook(ctx context.Context, book *models.Book) error {
	query := "INSERT INTO books (title, author, year, created_at, updated_at) VALUES (?, ?, ?, ?, ?)"
	result, err := r.DB.ExecContext(ctx, query, book.Title, book.Author, book.Year, book.CreatedAt, book.UpdatedAt)
	if err != nil {
		log.Error().Err(err).Msg("[BookRepository] Failed to save data to database")
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		log.Error().Err(err).Msg("[BookRepository] Failed to get newly created data ID")
		return err
	}
	book.ID = int(id)
	log.Info().Int("id", book.ID).Msg("[BookRepository] Successfully saved data to database")
	return nil
}

func (r *mysqlBookRepository) UpdateBook(ctx context.Context, book *models.Book) error {
	query := "UPDATE books SET title = ?, author = ?, year = ?, updated_at = ? WHERE id = ?"
	_, err := r.DB.ExecContext(ctx, query, book.Title, book.Author, book.Year, book.UpdatedAt, book.ID)
	if err != nil {
		log.Error().Err(err).Int("id", book.ID).Msg("[BookRepository] Failed to update data in database")
		return err
	}

	log.Info().Int("id", book.ID).Msg("[BookRepository] Successfully updated data in database")
	return nil
}

func (r *mysqlBookRepository) DeleteBook(ctx context.Context, id int) error {
	query := "DELETE FROM books WHERE id = ?"
	_, err := r.DB.ExecContext(ctx, query, id)
	if err != nil {
		log.Error().Err(err).Int("id", id).Msg("[BookRepository] Failed to delete data from database")
		return err
	}
	log.Info().Int("id", id).Msg("[BookRepository] Successfully deleted data from database")
	return nil
}

func (r *mysqlBookRepository) FindByTitle(ctx context.Context, title string) ([]models.Book, error) {
	query := "SELECT id, title, author, year, created_at, updated_at FROM books WHERE title = ?"
	rows, err := r.DB.QueryContext(ctx, query, title)
	if err != nil {
		log.Error().Err(err).Msg("[BookRepository] Failed to get all books by title from database")
		return nil, err
	}
	defer rows.Close()

	var books []models.Book
	for rows.Next() {
		var book models.Book
		if err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Year, &book.CreatedAt, &book.UpdatedAt); err != nil {
			log.Error().Err(err).Msg("[BookRepository] Failed to read book data from query results")
			return nil, err
		}
		books = append(books, book)
	}

	if err := rows.Err(); err != nil {
		log.Error().Err(err).Msg("[BookRepository] Failed to get all books by title from database")
		return nil, err
	}

	log.Info().Msg("[BookRepository] Successfully got all books by title from database")
	return books, nil
}
