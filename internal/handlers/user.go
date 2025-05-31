package handlers

import (
	"net/http"

	"github.com/Lumina-Enterprise-Solutions/prism-common-libs/pkg/utils"
	userModels "github.com/Lumina-Enterprise-Solutions/prism-user-service/internal/models"
	"github.com/Lumina-Enterprise-Solutions/prism-user-service/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type UserHandler struct {
	userService services.UserService
	logger      *logrus.Logger // Change from commonLogger.Logger to *logrus.Logger
}

func NewUserHandler(userService services.UserService, logger *logrus.Logger) *UserHandler { // Update parameter type
	return &UserHandler{
		userService: userService,
		logger:      logger,
	}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var req userModels.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, utils.FormatValidationErrors(err))
		return
	}

	tenantID := h.getTenantID(c)
	user, err := h.userService.CreateUser(tenantID, &req)
	if err != nil {
		if err == services.ErrUserExists {
			utils.ErrorResponse(c, http.StatusConflict, "User already exists", err)
			return
		}
		h.logger.Errorf("Error creating user: %v", err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create user", err)
		return
	}

	utils.SuccessResponse(c, "User created successfully", user)
}

func (h *UserHandler) GetUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid user ID", err)
		return
	}

	tenantID := h.getTenantID(c)
	user, err := h.userService.GetUser(tenantID, id)
	if err != nil {
		if err == services.ErrUserNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "User not found", err)
			return
		}
		h.logger.Errorf("Error fetching user: %v", err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch user", err)
		return
	}

	utils.SuccessResponse(c, "User retrieved successfully", user)
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid user ID", err)
		return
	}

	var req userModels.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, utils.FormatValidationErrors(err))
		return
	}

	tenantID := h.getTenantID(c)
	user, err := h.userService.UpdateUser(tenantID, id, &req)
	if err != nil {
		if err == services.ErrUserNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "User not found", err)
			return
		}
		h.logger.Errorf("Error updating user: %v", err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update user", err)
		return
	}

	utils.SuccessResponse(c, "User updated successfully", user)
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid user ID", err)
		return
	}

	tenantID := h.getTenantID(c)
	err = h.userService.DeleteUser(tenantID, id)
	if err != nil {
		if err == services.ErrUserNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "User not found", err)
			return
		}
		h.logger.Errorf("Error deleting user: %v", err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete user", err)
		return
	}

	utils.SuccessResponse(c, "User deleted successfully", nil)
}

func (h *UserHandler) ListUsers(c *gin.Context) {
	var query userModels.UserQueryRequest
	if err := c.ShouldBindQuery(&query); err != nil {
		utils.ValidationErrorResponse(c, utils.FormatValidationErrors(err))
		return
	}

	tenantID := h.getTenantID(c)
	users, err := h.userService.ListUsers(tenantID, &query)
	if err != nil {
		h.logger.Errorf("Error listing users: %v", err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to list users", err)
		return
	}

	utils.SuccessResponse(c, "Users retrieved successfully", users)
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	userID := h.getUserID(c)
	if userID == uuid.Nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	tenantID := h.getTenantID(c)
	user, err := h.userService.GetUser(tenantID, userID)
	if err != nil {
		if err == services.ErrUserNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "User not found", err)
			return
		}
		h.logger.Errorf("Error fetching user profile: %v", err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch profile", err)
		return
	}

	utils.SuccessResponse(c, "Profile retrieved successfully", user)
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID := h.getUserID(c)
	if userID == uuid.Nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	var req userModels.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, utils.FormatValidationErrors(err))
		return
	}

	tenantID := h.getTenantID(c)
	user, err := h.userService.UpdateProfile(tenantID, userID, &req)
	if err != nil {
		if err == services.ErrUserNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "User not found", err)
			return
		}
		h.logger.Errorf("Error updating profile: %v", err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update profile", err)
		return
	}

	utils.SuccessResponse(c, "Profile updated successfully", user)
}

func (h *UserHandler) getTenantID(c *gin.Context) string {
	if tenantID, exists := c.Get("tenant_id"); exists {
		if tid, ok := tenantID.(string); ok {
			return tid
		}
	}
	return "default"
}

func (h *UserHandler) getUserID(c *gin.Context) uuid.UUID {
	if userID, exists := c.Get("user_id"); exists {
		switch v := userID.(type) {
		case string:
			if id, err := uuid.Parse(v); err == nil {
				return id
			}
		case uuid.UUID:
			return v
		}
	}
	return uuid.Nil
}
