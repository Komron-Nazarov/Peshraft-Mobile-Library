// package repositories

// import (
// 	"mobile-library/internal/models"
// )

// type UserRepository struct {
// 	db *DB
// }

// func NewUserRepository(db *DB) *UserRepository {
// 	return &UserRepository{db: db}
// }

// func (r *UserRepository) Create(user *models.User) error {
// 	return r.db.conn.QueryRow(
// 		`INSERT INTO users (name, email, phone, password, created_at) 
// 		 VALUES ($1, $2, $3, $4, NOW()) RETURNING id, created_at`,
// 		user.Name, user.Email, user.Phone, user.Password,
// 	).Scan(&user.ID, &user.CreatedAt)
// }

// func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
// 	user := &models.User{}
// 	err := r.db.conn.QueryRow(
// 		`SELECT id, name, email, phone, password, created_at FROM users WHERE email = $1`,
// 		email,
// 	).Scan(&user.ID, &user.Name, &user.Email, &user.Phone, &user.Password, &user.CreatedAt)

// 	if err != nil {
// 		return nil, err
// 	}
// 	return user, nil
// }

// func (r *UserRepository) FindByID(id uint) (*models.User, error) {
// 	user := &models.User{}
// 	err := r.db.conn.QueryRow(
// 		`SELECT id, name, email, phone, password, created_at FROM users WHERE id = $1`,
// 		id,
// 	).Scan(&user.ID, &user.Name, &user.Email, &user.Phone, &user.Password, &user.CreatedAt)

// 	if err != nil {
// 		return nil, err
// 	}
// 	return user, nil
// }

// func (r *UserRepository) Update(user *models.User) error {
// 	_, err := r.db.conn.Exec(
// 		`UPDATE users SET name = $1, phone = $2 WHERE id = $3`,
// 		user.Name, user.Phone, user.ID,
// 	)
// 	return err
// }

// func (r *UserRepository) UpdatePassword(id uint, hash string) error {
// 	_, err := r.db.conn.Exec(
// 		`UPDATE users SET password = $1 WHERE id = $2`,
// 		hash, id,
// 	)
// 	return err
// }






package repositories

import (
	"mobile-library/internal/models"
)

type UserRepository struct {
	db *DB
}

func NewUserRepository(db *DB) *UserRepository {
	return &UserRepository{db: db}
}

// 1. Учитываем новые поля при создании (если переданы)
func (r *UserRepository) Create(user *models.User) error {
	return r.db.conn.QueryRow(
		`INSERT INTO users (name, email, phone, password, role, job_position, date_of_birth, is_pending_admin, created_at) 
		 VALUES ($1, $2, $3, $4, COALESCE(NULLIF($5, ''), 'user'), COALESCE(NULLIF($6, ''), 'Student'), $7, $8, NOW()) 
		 RETURNING id, role, job_position, is_pending_admin, created_at`,
		user.Name, user.Email, user.Phone, user.Password, user.Role, user.JobPosition, user.DateOfBirth, user.IsPendingAdmin,
	).Scan(&user.ID, &user.Role, &user.JobPosition, &user.IsPendingAdmin, &user.CreatedAt)
}

// 2. Добавили новые колонки в SELECT и Scan при поиске по Email (для логина)
func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	user := &models.User{}
	err := r.db.conn.QueryRow(
		`SELECT id, name, email, phone, password, role, job_position, date_of_birth, is_pending_admin, created_at 
		 FROM users WHERE email = $1`,
		email,
	).Scan(&user.ID, &user.Name, &user.Email, &user.Phone, &user.Password, &user.Role, &user.JobPosition, &user.DateOfBirth, &user.IsPendingAdmin, &user.CreatedAt)

	if err != nil {
		return nil, err
	}
	return user, nil
}

// 3. Добавили новые колонки в SELECT и Scan при поиске по ID (для мидлвара)
func (r *UserRepository) FindByID(id uint) (*models.User, error) {
	user := &models.User{}
	err := r.db.conn.QueryRow(
		`SELECT id, name, email, phone, password, role, job_position, date_of_birth, is_pending_admin, created_at 
		 FROM users WHERE id = $1`,
		id,
	).Scan(&user.ID, &user.Name, &user.Email, &user.Phone, &user.Password, &user.Role, &user.JobPosition, &user.DateOfBirth, &user.IsPendingAdmin, &user.CreatedAt)

	if err != nil {
		return nil, err
	}
	return user, nil
}

// 4. Обновляем также и дату рождения, если пользователь меняет профиль
func (r *UserRepository) Update(user *models.User) error {
	_, err := r.db.conn.Exec(
		`UPDATE users SET name = $1, phone = $2, date_of_birth = $3 WHERE id = $4`,
		user.Name, user.Phone, user.DateOfBirth, user.ID,
	)
	return err
}

func (r *UserRepository) UpdatePassword(id uint, hash string) error {
	_, err := r.db.conn.Exec(
		`UPDATE users SET password = $1 WHERE id = $2`,
		hash, id,
	)
	return err
}