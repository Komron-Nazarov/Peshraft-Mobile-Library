package services

import (
	"fmt"
	"mobile-library/internal/models"
	"mobile-library/internal/repositories"
)

type BookService struct {
	bookRepo *repositories.BookRepository
}

func NewBookService(bookRepo *repositories.BookRepository) *BookService {
	return &BookService{bookRepo: bookRepo}
}

func (s *BookService) GetAll(page, pageSize int) ([]models.BookResponse, int, error) {
	books, total, err := s.bookRepo.FindAll(page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to fetch books: %w", err)
	}

	responses := make([]models.BookResponse, len(books))
	for i, b := range books {
		responses[i] = b.ToResponse()
	}
	return responses, total, nil
}

func (s *BookService) GetByID(id uint) (*models.BookResponse, error) {
	book, err := s.bookRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("book not found: %w", err)
	}

	resp := book.ToResponse()
	return &resp, nil
}

func (s *BookService) Search(params models.BookSearchParams) ([]models.BookResponse, int, error) {
	books, total, err := s.bookRepo.Search(params)
	if err != nil {
		return nil, 0, fmt.Errorf("search failed: %w", err)
	}

	responses := make([]models.BookResponse, len(books))
	for i, b := range books {
		responses[i] = b.ToResponse()
	}
	return responses, total, nil
}

func (s *BookService) Create(book *models.Book) error {
	if book.Title == "" || book.Author == "" {
		return fmt.Errorf("title and author are required")
	}

	err := s.bookRepo.Create(book)
	if err != nil {
		return fmt.Errorf("failed to create book: %w", err)
	}

	return nil
}


func (s *BookService) Update(id uint, book *models.Book) error {
	if book.Title == "" || book.Author == "" {
		return fmt.Errorf("title and author are required for update")
	}
	return s.bookRepo.Update(id, book)
}

func (s *BookService) Delete(id uint) error {
	return s.bookRepo.Delete(id)
}



// GetCategories возвращает список уникальных категорий
func (s *BookService) GetCategories() ([]string, error) {
	categories, err := s.bookRepo.GetCategories()
	if err != nil {
		return nil, fmt.Errorf("failed to get categories from repo: %w", err)
	}
	return categories, nil
}

func (s *BookService) CreateCategory(name string) error {
    return s.bookRepo.CreateCategory(name)
}

func (s *BookService) UpdateCategory(id string, name string) error {
    return s.bookRepo.UpdateCategory(id, name)
}

func (s *BookService) DeleteCategory(id string) error {
    return s.bookRepo.DeleteCategory(id)
}