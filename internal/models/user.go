package models

import (
	"time"

	commonModels "github.com/Lumina-Enterprise-Solutions/prism-common-libs/pkg/models"
	"github.com/google/uuid"
)

// CreateUserRequest represents the request payload for creating a user
type CreateUserRequest struct {
	Email     string   `json:"email" binding:"required,email"`
	FirstName string   `json:"first_name" binding:"required,min=2,max=50"`
	LastName  string   `json:"last_name" binding:"required,min=2,max=50"`
	Password  string   `json:"password" binding:"required,min=8"`
	Status    string   `json:"status" binding:"omitempty,oneof=active inactive pending"`
	RoleIDs   []string `json:"role_ids" binding:"omitempty"`
}

// UpdateUserRequest represents the request payload for updating a user
type UpdateUserRequest struct {
	FirstName *string  `json:"first_name" binding:"omitempty,min=2,max=50"`
	LastName  *string  `json:"last_name" binding:"omitempty,min=2,max=50"`
	Status    *string  `json:"status" binding:"omitempty,oneof=active inactive pending"`
	RoleIDs   []string `json:"role_ids" binding:"omitempty"`
}

// UpdateProfileRequest represents the request payload for updating user profile
type UpdateProfileRequest struct {
	FirstName *string `json:"first_name" binding:"omitempty,min=2,max=50"`
	LastName  *string `json:"last_name" binding:"omitempty,min=2,max=50"`
}

// UserResponse represents the response payload for user data
type UserResponse struct {
	ID        uuid.UUID           `json:"id"`
	Email     string              `json:"email"`
	FirstName string              `json:"first_name"`
	LastName  string              `json:"last_name"`
	Status    string              `json:"status"`
	Roles     []commonModels.Role `json:"roles"`
	CreatedAt time.Time           `json:"created_at"`
	UpdatedAt time.Time           `json:"updated_at"`
}

// UserQueryRequest represents the request payload for querying users
type UserQueryRequest struct {
	Page    int      `form:"page" binding:"omitempty,min=1"`
	Limit   int      `form:"limit" binding:"omitempty,min=1,max=100"`
	Status  string   `form:"status" binding:"omitempty,oneof=active inactive pending"`
	Role    string   `form:"role" binding:"omitempty"`
	Sort    string   `form:"sort" binding:"omitempty"`
	Search  string   `form:"search" binding:"omitempty"`
	RoleIDs []string `form:"role_ids" binding:"omitempty"`
}

// UserListResponse represents the response payload for user list
type UserListResponse struct {
	Users      []UserResponse `json:"users"`
	Total      int64          `json:"total"`
	Page       int            `json:"page"`
	Limit      int            `json:"limit"`
	TotalPages int            `json:"total_pages"`
}

// ToUserResponse converts a User model to UserResponse
func ToUserResponse(u commonModels.User) UserResponse {
	return UserResponse{
		ID:        u.ID,
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Status:    u.Status,
		Roles:     u.Roles,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
