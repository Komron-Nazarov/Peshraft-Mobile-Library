// package services

// import (
// 	"fmt"
// 	"mobile-library/internal/models"
// 	"mobile-library/internal/repositories"
// 	"time"
// )

// type BorrowService struct {
// 	borrowRepo *repositories.BorrowRepository
// 	bookRepo   *repositories.BookRepository
// 	notifRepo  *repositories.NotificationRepository
// }

// func NewBorrowService(br *repositories.BorrowRepository, bk *repositories.BookRepository, nr *repositories.NotificationRepository) *BorrowService {
// 	return &BorrowService{
// 		borrowRepo: br,
// 		bookRepo:   bk,
// 		notifRepo:  nr,
// 	}
// }

// // BorrowBook — логика выдачи книги
// func (s *BorrowService) BorrowBook(userID, bookID uint) (*models.BorrowResponse, error) {
// 	// 1. Проверяем, нет ли уже активной аренды этой же книги
// 	existing, _ := s.borrowRepo.FindByUserAndBook(userID, bookID)
// 	if existing != nil {
// 		return nil, fmt.Errorf("вы уже взяли эту книгу и еще не вернули её")
// 	}

// 	// 2. Проверяем наличие книги
// 	book, err := s.bookRepo.FindByID(bookID)
// 	if err != nil {
// 		return nil, fmt.Errorf("книга не найдена")
// 	}

// 	if book.AvailableCopies <= 0 {
// 		return nil, fmt.Errorf("нет доступных экземпляров")
// 	}

// 	// 3. Создаем запись об аренде
// 	borrow := &models.Borrow{
// 		UserID:     userID,
// 		BookID:     bookID,
// 		BorrowDate: time.Now(),
// 		DueDate:    time.Now().AddDate(0, 0, 14), // На 14 дней
// 		Status:     models.BorrowActive,
// 	}

// 	if err := s.borrowRepo.Create(borrow); err != nil {
// 		return nil, fmt.Errorf("ошибка при оформлении: %w", err)
// 	}

// 	// 4. Уменьшаем количество копий
// 	_ = s.bookRepo.UpdateAvailableCopies(bookID, -1)

// 	// 5. Уведомление
// 	// _ = s.notifRepo.Create(&models.Notification{
// 	// 	UserID:    userID,
// 	// 	Message:   fmt.Sprintf("Вы взяли книгу '%s'. Срок возврата: %s", book.Title, borrow.DueDate.Format("02.01.2006")),
// 	// 	CreatedAt: time.Now(),
// 	// })
// // 	uid := userID // создаем копию
// // _ = s.notifRepo.Create(&models.Notification{
// //     UserID:  &uid, // передаем адрес
// //     Message: fmt.Sprintf("Вы взяли книгу '%s'. Срок возврата: %s", book.Title, borrow.DueDate.Format("2006-01-02")),
// //     CreatedAt: time.Now(),
// // })
// // В методе BorrowBook
// userID := userID // твоя переменная uint
// err = s.notifRepo.Create(&models.Notification{
//     UserID:    &userID, // Добавь & перед переменной
//     Message:   fmt.Sprintf("Вы взяли книгу '%s'...", book.Title),
//     CreatedAt: time.Now(),
// })

// 	resp := borrow.ToResponse(book.Title)
// 	return &resp, nil
// }

// // ReturnBook — логика возврата
// func (s *BorrowService) ReturnBook(userID, bookID uint) (*models.BorrowResponse, error) {
// 	borrow, err := s.borrowRepo.FindByUserAndBook(userID, bookID)
// 	if err != nil {
// 		return nil, fmt.Errorf("активная аренда не найдена")
// 	}

// 	now := time.Now()
// 	status := models.BorrowReturned
// 	if now.After(borrow.DueDate) {
// 		status = models.BorrowOverdue
// 	}

// 	// Обновляем статус и дату возврата
// 	if err := s.borrowRepo.UpdateReturnStatus(borrow.ID, now, string(status)); err != nil {
// 		return nil, err
// 	}

// 	// Возвращаем книгу в фонд
// 	_ = s.bookRepo.UpdateAvailableCopies(bookID, 1)

// 	// Получаем инфо о книге для ответа
// 	book, _ := s.bookRepo.FindByID(bookID)

// 	resp := borrow.ToResponse(book.Title)
// 	return &resp, nil
// }

// // GetHistory — история для мобилки
// func (s *BorrowService) GetHistory(userID uint) ([]models.BorrowResponse, error) {
// 	records, err := s.borrowRepo.GetUserHistory(userID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	responses := make([]models.BorrowResponse, 0, len(records))
// 	for _, r := range records {
// 		responses = append(responses, r.Borrow.ToResponse(r.Book.Title))
// 	}
// 	return responses, nil
// }

// // CheckOverdueAndNotify — фоновая задача для проверки просрочки
// func (s *BorrowService) CheckOverdueAndNotify() error {
// 	overdueBorrows, err := s.borrowRepo.GetOverdueBorrows()
// 	if err != nil {
// 		return err
// 	}

// 	for _, b := range overdueBorrows {
// 		// _ = s.borrowRepo.MarkOverdue(b.ID)
// 		// _ = s.notifRepo.Create(&models.Notification{
// 		// 	UserID:  b.UserID,
// 		// 	Message: fmt.Sprintf("Внимание! Срок возврата книги '%s' истек", b.Title),
// 		// })
// 		uid := b.UserID
// _ = s.notifRepo.Create(&models.Notification{
//     UserID:  &uid,
//     Message: fmt.Sprintf("Внимание! Срок возврата книги '%s' истек", b.Title),
//     CreatedAt: time.Now(),
// })
// 	}
// 	return nil
// }




// package services

// import (
// 	"fmt"
// 	"mobile-library/internal/models"
// 	"mobile-library/internal/repositories"
// 	"time"
// )

// type BorrowService struct {
// 	borrowRepo *repositories.BorrowRepository
// 	bookRepo   *repositories.BookRepository
// 	notifRepo  *repositories.NotificationRepository
// }

// func NewBorrowService(br *repositories.BorrowRepository, bk *repositories.BookRepository, nr *repositories.NotificationRepository) *BorrowService {
// 	return &BorrowService{
// 		borrowRepo: br,
// 		bookRepo:   bk,
// 		notifRepo:  nr,
// 	}
// }

// // BorrowBook — логика выдачи книги
// func (s *BorrowService) BorrowBook(userID, bookID uint) (*models.BorrowResponse, error) {
// 	// 1. Проверяем, нет ли уже активной аренды этой же книги
// 	existing, _ := s.borrowRepo.FindByUserAndBook(userID, bookID)
// 	if existing != nil {
// 		return nil, fmt.Errorf("вы уже взяли эту книгу и еще не вернули её")
// 	}

// 	// 2. Проверяем наличие книги
// 	book, err := s.bookRepo.FindByID(bookID)
// 	if err != nil {
// 		return nil, fmt.Errorf("книга не найдена")
// 	}

// 	if book.AvailableCopies <= 0 {
// 		return nil, fmt.Errorf("нет доступных экземпляров")
// 	}

// 	// 3. Создаем запись об аренде
// 	borrow := &models.Borrow{
// 		UserID:     userID,
// 		BookID:     bookID,
// 		BorrowDate: time.Now(),
// 		DueDate:    time.Now().AddDate(0, 0, 14), // На 14 дней
// 		Status:     models.BorrowActive,
// 	}

// 	if err := s.borrowRepo.Create(borrow); err != nil {
// 		return nil, fmt.Errorf("ошибка при оформлении: %w", err)
// 	}

// 	// 4. Уменьшаем количество копий
// 	_ = s.bookRepo.UpdateAvailableCopies(bookID, -1)

// 	// 5. Уведомление
// 	uidCopy := userID // Создаем копию для передачи адреса
// 	_ = s.notifRepo.Create(&models.Notification{
// 		UserID:    &uidCopy,
// 		Message:   fmt.Sprintf("Вы взяли книгу '%s'.", book.Title),
// 		CreatedAt: time.Now(),
// 	})

// 	resp := borrow.ToResponse(book.Title)
// 	return &resp, nil
// }

// // ReturnBook — логика возврата
// func (s *BorrowService) ReturnBook(userID, bookID uint) (*models.BorrowResponse, error) {
// 	borrow, err := s.borrowRepo.FindByUserAndBook(userID, bookID)
// 	if err != nil {
// 		return nil, fmt.Errorf("активная аренда не найдена")
// 	}

// 	now := time.Now()
// 	status := models.BorrowReturned
// 	if now.After(borrow.DueDate) {
// 		status = models.BorrowOverdue
// 	}

// 	// Обновляем статус и дату возврата
// 	if err := s.borrowRepo.UpdateReturnStatus(borrow.ID, now, string(status)); err != nil {
// 		return nil, err
// 	}

// 	// Возвращаем книгу в фонд
// 	_ = s.bookRepo.UpdateAvailableCopies(bookID, 1)

// 	// Получаем инфо о книге для ответа
// 	book, _ := s.bookRepo.FindByID(bookID)

// 	resp := borrow.ToResponse(book.Title)
// 	return &resp, nil
// }

// // GetHistory — история для мобилки
// func (s *BorrowService) GetHistory(userID uint) ([]models.BorrowResponse, error) {
// 	records, err := s.borrowRepo.GetUserHistory(userID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	responses := make([]models.BorrowResponse, 0, len(records))
// 	for _, r := range records {
// 		responses = append(responses, r.Borrow.ToResponse(r.Book.Title))
// 	}
// 	return responses, nil
// }

// // CheckOverdueAndNotify — фоновая задача для проверки просрочки
// func (s *BorrowService) CheckOverdueAndNotify() error {
// 	overdueBorrows, err := s.borrowRepo.GetOverdueBorrows()
// 	if err != nil {
// 		return err
// 	}

// 	for _, b := range overdueBorrows {
// 		uid := b.UserID
// 		_ = s.notifRepo.Create(&models.Notification{
// 			UserID:    &uid,
// 			Message:   fmt.Sprintf("Внимание! Срок возврата книги '%s' истек", b.Title),
// 			CreatedAt: time.Now(),
// 		})
// 	}
// 	return nil
// }


package services

import (
	"fmt"
	"mobile-library/internal/models"
	"mobile-library/internal/repositories"
	"time"
)

type BorrowService struct {
	borrowRepo *repositories.BorrowRepository
	bookRepo   *repositories.BookRepository
	notifRepo  *repositories.NotificationRepository
}

func NewBorrowService(br *repositories.BorrowRepository, bk *repositories.BookRepository, nr *repositories.NotificationRepository) *BorrowService {
	return &BorrowService{
		borrowRepo: br,
		bookRepo:   bk,
		notifRepo:  nr,
	}
}

// BorrowBook — логика выдачи книги
func (s *BorrowService) BorrowBook(userID, bookID uint) (*models.BorrowResponse, error) {
	// 1. Проверяем, нет ли уже активной аренды этой же книги
	existing, _ := s.borrowRepo.FindByUserAndBook(userID, bookID)
	if existing != nil {
		return nil, fmt.Errorf("вы уже взяли эту книгу и еще не вернули её")
	}

	// 2. Проверяем наличие книги (передаемuserID)
	book, err := s.bookRepo.FindByID(bookID, userID)
	if err != nil {
		return nil, fmt.Errorf("книга не найдена")
	}

	if book.AvailableCopies <= 0 {
		return nil, fmt.Errorf("нет доступных экземпляров")
	}

	// 3. Создаем запись об аренде
	borrow := &models.Borrow{
		UserID:     userID,
		BookID:     bookID,
		BorrowDate: time.Now(),
		DueDate:    time.Now().AddDate(0, 0, 14), // На 14 дней
		Status:     models.BorrowActive,
	}

	if err := s.borrowRepo.Create(borrow); err != nil {
		return nil, fmt.Errorf("ошибка при оформлении: %w", err)
	}

	// 4. Уменьшаем количество копий
	_ = s.bookRepo.UpdateAvailableCopies(bookID, -1)

	// 5. Уведомление
	uidCopy := userID // Создаем копию для передачи адреса
	_ = s.notifRepo.Create(&models.Notification{
		UserID:    &uidCopy,
		Message:   fmt.Sprintf("Вы взяли книгу '%s'.", book.Title),
		CreatedAt: time.Now(),
	})

	resp := borrow.ToResponse(book.Title)
	return &resp, nil
}

// // ReturnBook — логика возврата
// func (s *BorrowService) ReturnBook(userID, bookID uint) (*models.BorrowResponse, error) {
// 	borrow, err := s.borrowRepo.FindByUserAndBook(userID, bookID)
// 	if err != nil {
// 		return nil, fmt.Errorf("активная аренда не найдена")
// 	}

// 	now := time.Now()
// 	status := models.BorrowReturned
// 	if now.After(borrow.DueDate) {
// 		status = models.BorrowOverdue
// 	}

// 	// Обновляем статус и дату возврата
// 	if err := s.borrowRepo.UpdateReturnStatus(borrow.ID, now, string(status)); err != nil {
// 		return nil, err
// 	}

// 	// Возвращаем книгу в фонд
// 	_ = s.bookRepo.UpdateAvailableCopies(bookID, 1)

// 	// Получаем инфо о книге для ответа
// 	book, _ := s.bookRepo.FindByID(bookID, userID)

// 	resp := borrow.ToResponse(book.Title)
// 	return &resp, nil
// }




// ReturnBook — логика возврата
func (s *BorrowService) ReturnBook(userID, bookID uint) (*models.BorrowResponse, error) {
	borrow, err := s.borrowRepo.FindByUserAndBook(userID, bookID)
	if err != nil {
		return nil, fmt.Errorf("активная аренда не найдена")
	}

	now := time.Now()
	status := models.BorrowReturned
	if now.After(borrow.DueDate) {
		status = models.BorrowOverdue
	}

	// Обновляем статус и дату возврата
	if err := s.borrowRepo.UpdateReturnStatus(borrow.ID, now, string(status)); err != nil {
		return nil, err
	}

	// Возвращаем книгу в фонд
	_ = s.bookRepo.UpdateAvailableCopies(bookID, 1)

	// Безопасно получаем инфо о книге для ответа, защищаясь от nil
	bookTitle := "Удаленная книга"
	if book, err := s.bookRepo.FindByID(bookID, userID); err == nil && book != nil {
		bookTitle = book.Title
	}

	resp := borrow.ToResponse(bookTitle)
	return &resp, nil
}


// GetHistory — история для мобилки
func (s *BorrowService) GetHistory(userID uint) ([]models.BorrowResponse, error) {
	records, err := s.borrowRepo.GetUserHistory(userID)
	if err != nil {
		return nil, err
	}

	responses := make([]models.BorrowResponse, 0, len(records))
	for _, r := range records {
		responses = append(responses, r.Borrow.ToResponse(r.Book.Title))
	}
	return responses, nil
}

// CheckOverdueAndNotify — фоновая задача для проверки просрочки
func (s *BorrowService) CheckOverdueAndNotify() error {
	overdueBorrows, err := s.borrowRepo.GetOverdueBorrows()
	if err != nil {
		return err
	}

	for _, b := range overdueBorrows {
		uid := b.UserID
		_ = s.notifRepo.Create(&models.Notification{
			UserID:    &uid,
			Message:   fmt.Sprintf("Внимание! Срок возврата книги '%s' истек", b.Title),
			CreatedAt: time.Now(),
		})
	}
	return nil
}