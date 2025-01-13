package handler

import (
	"net/http"
	"strconv"

	"backend-master/internal/model"

	"github.com/gin-gonic/gin"
)

// GetUser handles user retrieval
// @Summary Get user by ID
// @Description Get user information by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} response
// @Failure 400 {object} response
// @Failure 404 {object} response
// @Failure 500 {object} response
// @Security Bearer
// @Router /users/{id} [get]
func (h *Handler) GetUser(c *gin.Context) {
	id := c.Param("id")

	user, err := h.service.GetUser(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, successResponse(user.SanitizeUser()))
}

// UpdateUser handles user update
// @Summary Update user
// @Description Update user information
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param input body model.UserUpdate true "User update data"
// @Success 200 {object} response
// @Failure 400 {object} response
// @Failure 404 {object} response
// @Failure 500 {object} response
// @Security Bearer
// @Router /users/{id} [put]
func (h *Handler) UpdateUser(c *gin.Context) {
	id := c.Param("id")

	var input model.UserUpdate
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := h.service.UpdateUser(c.Request.Context(), id, input); err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, successResponse(nil))
}

// DeleteUser handles user deletion
// @Summary Delete user
// @Description Delete user by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} response
// @Failure 404 {object} response
// @Failure 500 {object} response
// @Security Bearer
// @Router /users/{id} [delete]
func (h *Handler) DeleteUser(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.DeleteUser(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, successResponse(nil))
}

// ListUsers handles user listing
// @Summary List users
// @Description Get paginated list of users
// @Tags users
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Success 200 {object} response
// @Failure 400 {object} response
// @Failure 500 {object} response
// @Security Bearer
// @Router /users [get]
func (h *Handler) ListUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	users, total, err := h.service.ListUsers(c.Request.Context(), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Sanitize user data
	var sanitizedUsers []interface{}
	for _, user := range users {
		sanitizedUsers = append(sanitizedUsers, user.SanitizeUser())
	}

	c.JSON(http.StatusOK, successResponse(gin.H{
		"users": sanitizedUsers,
		"pagination": gin.H{
			"current_page": page,
			"page_size":   pageSize,
			"total_items": total,
			"total_pages": (total + pageSize - 1) / pageSize,
		},
	}))
}
