package services

import (
	"errors"
	"testing"
	"time"

	"github.com/Lumina-Enterprise-Solutions/prism-common-libs/pkg/models"
	userModels "github.com/Lumina-Enterprise-Solutions/prism-user-service/internal/models"
	"github.com/Lumina-Enterprise-Solutions/prism-user-service/internal/repository"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestUserService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	logger := logrus.New()
	svc := NewUserService(mockRepo, logger)

	tenantID := "default"
	userID := uuid.New()
	roleID := uuid.New()
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("securepassword123"), bcrypt.DefaultCost)

	defaultUser := &models.User{
		BaseModel: models.BaseModel{
			ID:        userID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Email:        "test.user@example.com",
		FirstName:    "Test",
		LastName:     "User",
		PasswordHash: string(hashedPassword),
		Status:       "active",
		Roles: []models.Role{
			{
				BaseModel: models.BaseModel{
					ID:        roleID,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				Name:        "user",
				Permissions: map[string]interface{}{"users": []string{"read", "update_profile"}},
			},
		},
	}

	t.Run("CreateUser", func(t *testing.T) {
		tests := []struct {
			name        string
			req         *userModels.CreateUserRequest
			setupMock   func()
			expectError error
			expectUser  *userModels.UserResponse
		}{
			{
				name: "Success",
				req: &userModels.CreateUserRequest{
					Email:     "test.user@example.com",
					FirstName: "Test",
					LastName:  "User",
					Password:  "securepassword123",
					Status:    "active",
				},
				setupMock: func() {
					mockRepo.EXPECT().GetByEmail(tenantID, "test.user@example.com").Return(nil, nil)
					mockRepo.EXPECT().Create(tenantID, gomock.Any()).Return(nil)
					mockRepo.EXPECT().GetByID(tenantID, gomock.Any()).Return(defaultUser, nil)
				},
				expectUser: &userModels.UserResponse{ID: userID, Email: "test.user@example.com", FirstName: "Test", LastName: "User", Status: "active"},
			},
			{
				name: "UserExists",
				req: &userModels.CreateUserRequest{
					Email:     "test.user@example.com",
					FirstName: "Test",
					LastName:  "User",
					Password:  "securepassword123",
				},
				setupMock: func() {
					mockRepo.EXPECT().GetByEmail(tenantID, "test.user@example.com").Return(defaultUser, nil)
				},
				expectError: ErrUserExists,
			},
			{
				name: "CreateError",
				req: &userModels.CreateUserRequest{
					Email:     "test.user@example.com",
					FirstName: "Test",
					LastName:  "User",
					Password:  "securepassword123",
				},
				setupMock: func() {
					mockRepo.EXPECT().GetByEmail(tenantID, "test.user@example.com").Return(nil, nil)
					mockRepo.EXPECT().Create(tenantID, gomock.Any()).Return(errors.New("db error"))
				},
				expectError: errors.New("db error"),
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				tt.setupMock()
				user, err := svc.CreateUser(tenantID, tt.req)
				if tt.expectError != nil {
					assert.Error(t, err)
					assert.Equal(t, tt.expectError.Error(), err.Error())
					assert.Nil(t, user)
				} else {
					assert.NoError(t, err)
					assert.NotNil(t, user)
					assert.Equal(t, tt.expectUser.Email, user.Email)
					assert.Equal(t, tt.expectUser.FirstName, user.FirstName)
					assert.Equal(t, tt.expectUser.LastName, user.LastName)
					assert.Equal(t, tt.expectUser.Status, user.Status)
				}
			})
		}
	})

	t.Run("GetUser", func(t *testing.T) {
		tests := []struct {
			name        string
			id          uuid.UUID
			setupMock   func()
			expectError error
			expectUser  *userModels.UserResponse
		}{
			{
				name: "Success",
				id:   userID,
				setupMock: func() {
					mockRepo.EXPECT().GetByID(tenantID, userID).Return(defaultUser, nil)
				},
				expectUser: &userModels.UserResponse{ID: userID, Email: "test.user@example.com", FirstName: "Test", LastName: "User", Status: "active"},
			},
			{
				name: "NotFound",
				id:   userID,
				setupMock: func() {
					mockRepo.EXPECT().GetByID(tenantID, userID).Return(nil, nil)
				},
				expectError: ErrUserNotFound,
			},
			{
				name: "Error",
				id:   userID,
				setupMock: func() {
					mockRepo.EXPECT().GetByID(tenantID, userID).Return(nil, errors.New("db error"))
				},
				expectError: errors.New("db error"),
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				tt.setupMock()
				user, err := svc.GetUser(tenantID, tt.id)
				if tt.expectError != nil {
					assert.Error(t, err)
					assert.Equal(t, tt.expectError.Error(), err.Error())
					assert.Nil(t, user)
				} else {
					assert.NoError(t, err)
					assert.NotNil(t, user)
					assert.Equal(t, tt.expectUser.ID, user.ID)
				}
			})
		}
	})

	t.Run("GetUserByEmail", func(t *testing.T) {
		tests := []struct {
			name        string
			email       string
			setupMock   func()
			expectError error
			expectUser  *userModels.UserResponse
		}{
			{
				name:  "Success",
				email: "test.user@example.com",
				setupMock: func() {
					mockRepo.EXPECT().GetByEmail(tenantID, "test.user@example.com").Return(defaultUser, nil)
				},
				expectUser: &userModels.UserResponse{ID: userID, Email: "test.user@example.com", FirstName: "Test", LastName: "User", Status: "active"},
			},
			{
				name:  "NotFound",
				email: "notfound@example.com",
				setupMock: func() {
					mockRepo.EXPECT().GetByEmail(tenantID, "notfound@example.com").Return(nil, nil)
				},
				expectError: ErrUserNotFound,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				tt.setupMock()
				user, err := svc.GetUserByEmail(tenantID, tt.email)
				if tt.expectError != nil {
					assert.Error(t, err)
					assert.Equal(t, tt.expectError.Error(), err.Error())
					assert.Nil(t, user)
				} else {
					assert.NoError(t, err)
					assert.NotNil(t, user)
					assert.Equal(t, tt.expectUser.Email, user.Email)
				}
			})
		}
	})

	t.Run("UpdateUser", func(t *testing.T) {
		tests := []struct {
			name        string
			id          uuid.UUID
			req         *userModels.UpdateUserRequest
			setupMock   func()
			expectError error
			expectUser  *userModels.UserResponse
		}{
			{
				name: "Success",
				id:   userID,
				req: &userModels.UpdateUserRequest{
					FirstName: stringPtr("Updated"),
					LastName:  stringPtr("User"),
					Status:    stringPtr("active"),
				},
				setupMock: func() {
					mockRepo.EXPECT().GetByID(tenantID, userID).Return(defaultUser, nil)
					mockRepo.EXPECT().Update(tenantID, userID, gomock.Any()).Return(nil)
					updatedUser := *defaultUser
					updatedUser.FirstName = "Updated"
					mockRepo.EXPECT().GetByID(tenantID, userID).Return(&updatedUser, nil)
				},
				expectUser: &userModels.UserResponse{ID: userID, Email: "test.user@example.com", FirstName: "Updated", LastName: "User", Status: "active"},
			},
			{
				name: "NotFound",
				id:   userID,
				req: &userModels.UpdateUserRequest{
					FirstName: stringPtr("Updated"),
				},
				setupMock: func() {
					mockRepo.EXPECT().GetByID(tenantID, userID).Return(nil, nil)
				},
				expectError: ErrUserNotFound,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				tt.setupMock()
				user, err := svc.UpdateUser(tenantID, tt.id, tt.req)
				if tt.expectError != nil {
					assert.Error(t, err)
					assert.Equal(t, tt.expectError.Error(), err.Error())
					assert.Nil(t, user)
				} else {
					assert.NoError(t, err)
					assert.NotNil(t, user)
					assert.Equal(t, tt.expectUser.FirstName, user.FirstName)
				}
			})
		}
	})

	t.Run("DeleteUser", func(t *testing.T) {
		tests := []struct {
			name        string
			id          uuid.UUID
			setupMock   func()
			expectError error
		}{
			{
				name: "Success",
				id:   userID,
				setupMock: func() {
					mockRepo.EXPECT().GetByID(tenantID, userID).Return(defaultUser, nil)
					mockRepo.EXPECT().Delete(tenantID, userID).Return(nil)
				},
			},
			{
				name: "NotFound",
				id:   userID,
				setupMock: func() {
					mockRepo.EXPECT().GetByID(tenantID, userID).Return(nil, nil)
				},
				expectError: ErrUserNotFound,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				tt.setupMock()
				err := svc.DeleteUser(tenantID, tt.id)
				if tt.expectError != nil {
					assert.Error(t, err)
					assert.Equal(t, tt.expectError.Error(), err.Error())
				} else {
					assert.NoError(t, err)
				}
			})
		}
	})

	t.Run("ListUsers", func(t *testing.T) {
		tests := []struct {
			name        string
			query       *userModels.UserQueryRequest
			setupMock   func()
			expectError error
			expectResp  *userModels.UserListResponse
		}{
			{
				name: "Success",
				query: &userModels.UserQueryRequest{
					Page:   1,
					Limit:  20,
					Status: "active",
				},
				setupMock: func() {
					mockRepo.EXPECT().List(tenantID, gomock.Any()).Return([]models.User{*defaultUser}, int64(1), nil)
				},
				expectResp: &userModels.UserListResponse{
					Users:      []userModels.UserResponse{{ID: userID, Email: "test.user@example.com", FirstName: "Test", LastName: "User", Status: "active"}},
					Total:      1,
					Page:       1,
					Limit:      20,
					TotalPages: 1,
				},
			},
			{
				name:  "InvalidPage",
				query: &userModels.UserQueryRequest{Page: 0, Limit: 20},
				setupMock: func() {
					mockRepo.EXPECT().List(tenantID, gomock.Any()).Return([]models.User{*defaultUser}, int64(1), nil)
				},
				expectResp: &userModels.UserListResponse{Total: 1, Page: 1, Limit: 20, TotalPages: 1},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				tt.setupMock()
				resp, err := svc.ListUsers(tenantID, tt.query)
				if tt.expectError != nil {
					assert.Error(t, err)
					assert.Equal(t, tt.expectError.Error(), err.Error())
					assert.Nil(t, resp)
				} else {
					assert.NoError(t, err)
					assert.NotNil(t, resp)
					assert.Equal(t, tt.expectResp.Total, resp.Total)
					assert.Equal(t, tt.expectResp.Page, resp.Page)
					assert.Equal(t, tt.expectResp.Limit, resp.Limit)
					if len(tt.expectResp.Users) > 0 {
						assert.Equal(t, tt.expectResp.Users[0].Email, resp.Users[0].Email)
					}
				}
			})
		}
	})

	t.Run("UpdateProfile", func(t *testing.T) {
		tests := []struct {
			name        string
			id          uuid.UUID
			req         *userModels.UpdateProfileRequest
			setupMock   func()
			expectError error
			expectUser  *userModels.UserResponse
		}{
			{
				name: "Success",
				id:   userID,
				req: &userModels.UpdateProfileRequest{
					FirstName: stringPtr("Updated"),
					LastName:  stringPtr("Profile"),
				},
				setupMock: func() {
					mockRepo.EXPECT().GetByID(tenantID, userID).Return(defaultUser, nil)
					mockRepo.EXPECT().Update(tenantID, userID, gomock.Any()).Return(nil)
					updatedUser := *defaultUser
					updatedUser.FirstName = "Updated"
					updatedUser.LastName = "Profile"
					mockRepo.EXPECT().GetByID(tenantID, userID).Return(&updatedUser, nil)
				},
				expectUser: &userModels.UserResponse{ID: userID, Email: "test.user@example.com", FirstName: "Updated", LastName: "Profile", Status: "active"},
			},
			{
				name: "NoUpdates",
				id:   userID,
				req:  &userModels.UpdateProfileRequest{},
				setupMock: func() {
					mockRepo.EXPECT().GetByID(tenantID, userID).Return(defaultUser, nil)
				},
				expectUser: &userModels.UserResponse{ID: userID, Email: "test.user@example.com", FirstName: "Test", LastName: "User", Status: "active"},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				tt.setupMock()
				user, err := svc.UpdateProfile(tenantID, tt.id, tt.req)
				if tt.expectError != nil {
					assert.Error(t, err)
					assert.Equal(t, tt.expectError.Error(), err.Error())
					assert.Nil(t, user)
				} else {
					assert.NoError(t, err)
					assert.NotNil(t, user)
					assert.Equal(t, tt.expectUser.FirstName, user.FirstName)
					assert.Equal(t, tt.expectUser.LastName, user.LastName)
				}
			})
		}
	})
}

// Helper function to create a string pointer
func stringPtr(s string) *string {
	return &s
}
