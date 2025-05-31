package repository

import (
	"errors"

	"github.com/Lumina-Enterprise-Solutions/prism-common-libs/pkg/database"
	commonModels "github.com/Lumina-Enterprise-Solutions/prism-common-libs/pkg/models"
	userModels "github.com/Lumina-Enterprise-Solutions/prism-user-service/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(tenantID string, user *commonModels.User) error
	GetByID(tenantID string, id uuid.UUID) (*commonModels.User, error)
	GetByEmail(tenantID string, email string) (*commonModels.User, error)
	Update(tenantID string, id uuid.UUID, updates map[string]interface{}) error
	Delete(tenantID string, id uuid.UUID) error
	List(tenantID string, query *userModels.UserQueryRequest) ([]commonModels.User, int64, error)
}

type userRepository struct {
	db *database.PostgresDB
}

func NewUserRepository(db *database.PostgresDB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(tenantID string, user *commonModels.User) error {
	db := r.db.WithTenant(tenantID)
	return db.Create(user).Error
}

func (r *userRepository) GetByID(tenantID string, id uuid.UUID) (*commonModels.User, error) {
	var user commonModels.User
	db := r.db.WithTenant(tenantID)

	err := db.Preload("Roles").Where("id = ?", id).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) GetByEmail(tenantID string, email string) (*commonModels.User, error) {
	var user commonModels.User
	db := r.db.WithTenant(tenantID)

	err := db.Preload("Roles").Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) Update(tenantID string, id uuid.UUID, updates map[string]interface{}) error {
	db := r.db.WithTenant(tenantID)
	return db.Model(&commonModels.User{}).Where("id = ?", id).Updates(updates).Error
}

func (r *userRepository) Delete(tenantID string, id uuid.UUID) error {
	db := r.db.WithTenant(tenantID)
	return db.Where("id = ?", id).Delete(&commonModels.User{}).Error
}

func (r *userRepository) List(tenantID string, query *userModels.UserQueryRequest) ([]commonModels.User, int64, error) {
	var users []commonModels.User
	var total int64

	db := r.db.WithTenant(tenantID)

	// Build query
	queryBuilder := db.Model(&commonModels.User{}).Preload("Roles")

	// Apply filters
	if query.Status != "" {
		queryBuilder = queryBuilder.Where("status = ?", query.Status)
	}

	if query.Search != "" {
		searchTerm := "%" + query.Search + "%"
		queryBuilder = queryBuilder.Where(
			"first_name ILIKE ? OR last_name ILIKE ? OR email ILIKE ?",
			searchTerm, searchTerm, searchTerm,
		)
	}

	if len(query.RoleIDs) > 0 {
		roleUUIDs := make([]uuid.UUID, len(query.RoleIDs))
		for i, roleID := range query.RoleIDs {
			if parsed, err := uuid.Parse(roleID); err == nil {
				roleUUIDs[i] = parsed
			}
		}
		queryBuilder = queryBuilder.Joins("JOIN user_roles ON users.id = user_roles.user_id").
			Where("user_roles.role_id IN ?", roleUUIDs)
	}

	// Count total records
	if err := queryBuilder.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply sorting
	if query.Sort != "" {
		queryBuilder = r.applySorting(queryBuilder, query.Sort)
	} else {
		queryBuilder = queryBuilder.Order("created_at DESC")
	}

	// Apply pagination
	if query.Page > 0 && query.Limit > 0 {
		offset := (query.Page - 1) * query.Limit
		queryBuilder = queryBuilder.Offset(offset).Limit(query.Limit)
	}

	err := queryBuilder.Find(&users).Error
	return users, total, err
}

func (r *userRepository) applySorting(db *gorm.DB, sort string) *gorm.DB {
	switch sort {
	case "email:asc":
		return db.Order("email ASC")
	case "email:desc":
		return db.Order("email DESC")
	case "created_at:asc":
		return db.Order("created_at ASC")
	case "created_at:desc":
		return db.Order("created_at DESC")
	case "first_name:asc":
		return db.Order("first_name ASC")
	case "first_name:desc":
		return db.Order("first_name DESC")
	case "last_name:asc":
		return db.Order("last_name ASC")
	case "last_name:desc":
		return db.Order("last_name DESC")
	default:
		return db.Order("created_at DESC")
	}
}
