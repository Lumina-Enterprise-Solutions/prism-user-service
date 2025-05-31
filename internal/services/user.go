package services

import (
	"errors"
	"math"

	commonModels "github.com/Lumina-Enterprise-Solutions/prism-common-libs/pkg/models"
	userModels "github.com/Lumina-Enterprise-Solutions/prism-user-service/internal/models"
	"github.com/Lumina-Enterprise-Solutions/prism-user-service/internal/repository"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrUserExists      = errors.New("user already exists")
	ErrInvalidPassword = errors.New("invalid password")
	ErrUnauthorized    = errors.New("unauthorized")
)

type UserService interface {
	CreateUser(tenantID string, req *userModels.CreateUserRequest) (*userModels.UserResponse, error)
	GetUser(tenantID string, id uuid.UUID) (*userModels.UserResponse, error)
	GetUserByEmail(tenantID string, email string) (*userModels.UserResponse, error)
	UpdateUser(tenantID string, id uuid.UUID, req *userModels.UpdateUserRequest) (*userModels.UserResponse, error)
	DeleteUser(tenantID string, id uuid.UUID) error
	ListUsers(tenantID string, query *userModels.UserQueryRequest) (*userModels.UserListResponse, error)
	UpdateProfile(tenantID string, userID uuid.UUID, req *userModels.UpdateProfileRequest) (*userModels.UserResponse, error)
}

type userService struct {
	userRepo repository.UserRepository
	logger   *logrus.Logger // Change from commonLogger.Logger to *logrus.Logger
}

func NewUserService(userRepo repository.UserRepository, logger *logrus.Logger) UserService { // Update parameter type
	return &userService{
		userRepo: userRepo,
		logger:   logger,
	}
}

func (s *userService) CreateUser(tenantID string, req *userModels.CreateUserRequest) (*userModels.UserResponse, error) {
	// Check if user already exists
	existingUser, err := s.userRepo.GetByEmail(tenantID, req.Email)
	if err != nil {
		s.logger.Errorf("Error checking existing user: %v", err)
		return nil, err
	}
	if existingUser != nil {
		return nil, ErrUserExists
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Errorf("Error hashing password: %v", err)
		return nil, err
	}

	// Set default status if not provided
	status := req.Status
	if status == "" {
		status = "active"
	}

	// Create user
	user := &commonModels.User{
		BaseModel: commonModels.BaseModel{
			ID: uuid.New(),
		},
		Email:        req.Email,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		PasswordHash: string(hashedPassword),
		Status:       status,
	}

	err = s.userRepo.Create(tenantID, user)
	if err != nil {
		s.logger.Errorf("Error creating user: %v", err)
		return nil, err
	}

	// Get created user with roles
	createdUser, err := s.userRepo.GetByID(tenantID, user.ID)
	if err != nil {
		s.logger.Errorf("Error fetching created user: %v", err)
		return nil, err
	}

	response := userModels.ToUserResponse(*createdUser)
	s.logger.Infof("User created successfully: %s", user.Email)

	return &response, nil
}

func (s *userService) GetUser(tenantID string, id uuid.UUID) (*userModels.UserResponse, error) {
	user, err := s.userRepo.GetByID(tenantID, id)
	if err != nil {
		s.logger.Errorf("Error fetching user: %v", err)
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	response := userModels.ToUserResponse(*user)
	return &response, nil
}

func (s *userService) GetUserByEmail(tenantID string, email string) (*userModels.UserResponse, error) {
	user, err := s.userRepo.GetByEmail(tenantID, email)
	if err != nil {
		s.logger.Errorf("Error fetching user by email: %v", err)
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	response := userModels.ToUserResponse(*user)
	return &response, nil
}

func (s *userService) UpdateUser(tenantID string, id uuid.UUID, req *userModels.UpdateUserRequest) (*userModels.UserResponse, error) {
	// Check if user exists
	user, err := s.userRepo.GetByID(tenantID, id)
	if err != nil {
		s.logger.Errorf("Error fetching user: %v", err)
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	// Build updates map
	updates := make(map[string]interface{})

	if req.FirstName != nil {
		updates["first_name"] = *req.FirstName
	}
	if req.LastName != nil {
		updates["last_name"] = *req.LastName
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}

	// Update user
	err = s.userRepo.Update(tenantID, id, updates)
	if err != nil {
		s.logger.Errorf("Error updating user: %v", err)
		return nil, err
	}

	// Get updated user
	updatedUser, err := s.userRepo.GetByID(tenantID, id)
	if err != nil {
		s.logger.Errorf("Error fetching updated user: %v", err)
		return nil, err
	}

	response := userModels.ToUserResponse(*updatedUser)
	s.logger.Infof("User updated successfully: %s", updatedUser.Email)

	return &response, nil
}

func (s *userService) DeleteUser(tenantID string, id uuid.UUID) error {
	// Check if user exists
	user, err := s.userRepo.GetByID(tenantID, id)
	if err != nil {
		s.logger.Errorf("Error fetching user: %v", err)
		return err
	}
	if user == nil {
		return ErrUserNotFound
	}

	err = s.userRepo.Delete(tenantID, id)
	if err != nil {
		s.logger.Errorf("Error deleting user: %v", err)
		return err
	}

	s.logger.Infof("User deleted successfully: %s", user.Email)
	return nil
}

func (s *userService) ListUsers(tenantID string, query *userModels.UserQueryRequest) (*userModels.UserListResponse, error) {
	// Set defaults
	if query.Page <= 0 {
		query.Page = 1
	}
	if query.Limit <= 0 {
		query.Limit = 20
	}

	users, total, err := s.userRepo.List(tenantID, query)
	if err != nil {
		s.logger.Errorf("Error listing users: %v", err)
		return nil, err
	}

	// Convert to response format
	userResponses := make([]userModels.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = userModels.ToUserResponse(user)
	}

	totalPages := int(math.Ceil(float64(total) / float64(query.Limit)))

	response := &userModels.UserListResponse{
		Users:      userResponses,
		Total:      total,
		Page:       query.Page,
		Limit:      query.Limit,
		TotalPages: totalPages,
	}

	return response, nil
}

func (s *userService) UpdateProfile(tenantID string, userID uuid.UUID, req *userModels.UpdateProfileRequest) (*userModels.UserResponse, error) {
	// Check if user exists
	user, err := s.userRepo.GetByID(tenantID, userID)
	if err != nil {
		s.logger.Errorf("Error fetching user: %v", err)
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	// Build updates map
	updates := make(map[string]interface{})

	if req.FirstName != nil {
		updates["first_name"] = *req.FirstName
	}
	if req.LastName != nil {
		updates["last_name"] = *req.LastName
	}

	if len(updates) == 0 {
		response := userModels.ToUserResponse(*user)
		return &response, nil
	}

	// Update user
	err = s.userRepo.Update(tenantID, userID, updates)
	if err != nil {
		s.logger.Errorf("Error updating user profile: %v", err)
		return nil, err
	}

	// Get updated user
	updatedUser, err := s.userRepo.GetByID(tenantID, userID)
	if err != nil {
		s.logger.Errorf("Error fetching updated user: %v", err)
		return nil, err
	}

	response := userModels.ToUserResponse(*updatedUser)
	s.logger.Infof("User profile updated successfully: %s", updatedUser.Email)

	return &response, nil
}
