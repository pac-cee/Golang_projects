package handlers

import (
    "net/http"
    "wallet-app/internal/services"
    "wallet-app/pkg/utils"

    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
)

type WalletHandler struct {
    service *services.WalletService
}

func NewWalletHandler(service *services.WalletService) *WalletHandler {
    return &WalletHandler{service: service}
}

type CreateWalletRequest struct {
    UserID   string `json:"user_id" binding:"required"`
    Currency string `json:"currency" binding:"required"`
}

func (h *WalletHandler) Create(c *gin.Context) {
    var req CreateWalletRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request payload")
        return
    }

    userID, err := uuid.Parse(req.UserID)
    if err != nil {
        utils.ErrorResponse(c, http.StatusBadRequest, "Invalid user ID")
        return
    }

    wallet, err := h.service.CreateWallet(c.Request.Context(), userID, req.Currency)
    if err != nil {
        utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create wallet")
        return
    }

    c.JSON(http.StatusCreated, wallet)
}