// package models

// import "time"

// type User struct {
// 	ID        uint      `json:"id" gorm:"primaryKey"`
// 	Name      string    `json:"name" gorm:"not null"`
// 	Email     string    `json:"email" gorm:"unique;not null"`
// 	Phone     string    `json:"phone" gorm:"not null"`
// 	Password  string    `json:"-" gorm:"not null"`
// 	DateOfBirth string    `json:"date_of_birth"`
// 	CreatedAt time.Time `json:"created_at"`
// }

// type RegisterRequest struct {
// 	Name            string `json:"name" binding:"required,min=2,max=100"`
// 	Email           string `json:"email" binding:"required,email"`
// 	Phone           string `json:"phone" binding:"required,min=10,max=20"`
// 	Password        string `json:"password" binding:"required,min=6,max=100"`
// 	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=Password"` // Валидация: должно быть равно полю Password
// 	DateOfBirth     string `json:"date_of_birth" binding:"required" example:"2008-04-12"`
// }

// type LoginRequest struct {
// 	Email    string `json:"email" binding:"required,email"`
// 	Password string `json:"password" binding:"required"`
// }

// type UserResponse struct {
// 	ID        uint      `json:"id"`
// 	Name      string    `json:"name"`
// 	Email     string    `json:"email"`
// 	Phone     string    `json:"phone"`
// 	DateOfBirth string    `json:"date_of_birth"`
// 	CreatedAt time.Time `json:"created_at"`
// }

// type UpdateProfileRequest struct {
// 	Name  string `json:"name" binding:"omitempty,min=2,max=100"`
// 	Phone string `json:"phone" binding:"omitempty,min=10,max=20"`
// }

// type ChangePasswordRequest struct {
// 	OldPassword string `json:"old_password" binding:"required,min=6"`
// 	NewPassword string `json:"new_password" binding:"required,min=6,max=100"`
// }

// func (u *User) ToResponse() UserResponse {
// 	return UserResponse{
// 		ID:        u.ID,
// 		Name:      u.Name,
// 		Email:     u.Email,
// 		Phone:     u.Phone,
// 		CreatedAt: u.CreatedAt,
// 	}
// }




package models

import "time"

type User struct {
    ID          uint      `json:"id" gorm:"primaryKey"`
    Name        string    `json:"name" gorm:"not null"`
    Email       string    `json:"email" gorm:"unique;not null"`
    Phone       string    `json:"phone" gorm:"not null"`
    Password    string    `json:"-" gorm:"not null"`
    DateOfBirth string    `json:"date_of_birth"`
    Role           string    `json:"role" gorm:"default:'user'"`           // "user" или "admin"
    JobPosition    string    `json:"job_position" gorm:"default:'Student'"` // "Student", "Volunteer" и т.д.
    IsPendingAdmin bool      `json:"is_pending_admin" gorm:"default:false"` // Для запроса прав админа
    CreatedAt   time.Time `json:"created_at"`
}

type RegisterRequest struct {
    Name            string `json:"name" binding:"required,min=2,max=100"`
    Email           string `json:"email" binding:"required,email"`
    Phone           string `json:"phone" binding:"required,min=10,max=20"`
    Password        string `json:"password" binding:"required,min=6,max=100"`
    ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=Password"`
    DateOfBirth     string `json:"date_of_birth" binding:"required" example:"2008-04-12"`
}

type LoginRequest struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required"`
}

type UserResponse struct {
   ID          uint      `json:"id"`
    Name        string    `json:"name"`
    Email       string    `json:"email"`
    Phone       string    `json:"phone"`
    Role        string    `json:"role"`
    JobPosition string    `json:"job_position"`
    DateOfBirth string    `json:"date_of_birth"`
    CreatedAt   time.Time `json:"created_at"`
}

type UpdateProfileRequest struct {
    Name        string `json:"name" binding:"omitempty,min=2,max=100"`
    Phone       string `json:"phone" binding:"omitempty,min=10,max=20"`
    DateOfBirth string `json:"date_of_birth" binding:"omitempty"`
}

type ChangePasswordRequest struct {
    OldPassword string `json:"old_password" binding:"required,min=6"`
    NewPassword string `json:"new_password" binding:"required,min=6,max=100"`
}

func (u *User) ToResponse() UserResponse {
    return UserResponse{
        ID:          u.ID,
        Name:        u.Name,
        Email:       u.Email,
        Phone:       u.Phone,
        Role:        u.Role,
        JobPosition: u.JobPosition,
        DateOfBirth: u.DateOfBirth,
        CreatedAt:   u.CreatedAt,
    }
}
