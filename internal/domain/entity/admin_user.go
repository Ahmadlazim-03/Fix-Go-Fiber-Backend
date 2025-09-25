package entity

import (
	"time"
	"gorm.io/gorm"
)

type AdminRole string

const (
	AdminRoleSuperAdmin AdminRole = "super_admin"
	AdminRoleAdmin      AdminRole = "admin"
	AdminRoleModerator  AdminRole = "moderator"
)

type AdminUser struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Username  string         `json:"username" gorm:"unique;not null;size:50"`
	Email     string         `json:"email" gorm:"unique;not null;size:100"`
	Password  string         `json:"-" gorm:"not null"`
	Role      AdminRole      `json:"role" gorm:"type:varchar(20);default:'moderator'"`
	IsActive  bool           `json:"is_active" gorm:"default:true"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

type AdminUserResponse struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Role      AdminRole `json:"role"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (a *AdminUser) ToResponse() *AdminUserResponse {
	return &AdminUserResponse{
		ID:        a.ID,
		Username:  a.Username,
		Email:     a.Email,
		Role:      a.Role,
		IsActive:  a.IsActive,
		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
	}
}

func (a *AdminUser) IsSuperAdmin() bool {
	return a.Role == AdminRoleSuperAdmin
}

func (a *AdminUser) IsAdmin() bool {
	return a.Role == AdminRoleAdmin || a.Role == AdminRoleSuperAdmin
}

func (a *AdminUser) CanModerate() bool {
	return a.Role == AdminRoleModerator || a.IsAdmin()
}

func (AdminUser) TableName() string {
	return "admin_users"
}