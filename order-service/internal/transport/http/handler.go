package http

import (
	"errors"
	"net/http"
	"order-service/internal/domain"
	"order-service/internal/usecase"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type OrderHandler struct {
	useCase *usecase.OrderUseCase
}

func NewOrderHandler(u *usecase.OrderUseCase) *OrderHandler {
	return &OrderHandler{useCase: u}
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req struct {
		CustomerID string `json:"customer_id" binding:"required"`
		ItemName   string `json:"item_name" binding:"required"`
		Amount     int64  `json:"amount" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order := &domain.Order{
		ID:         uuid.New().String(),
		CustomerID: req.CustomerID,
		ItemName:   req.ItemName,
		Amount:     req.Amount,
	}

	err := h.useCase.CreateOrder(c.Request.Context(), order)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidAmount) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, domain.ErrPaymentServiceDown) {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, order)
}

func (h *OrderHandler) CancelOrder(c *gin.Context) {
	id := c.Param("id")
	err := h.useCase.CancelOrder(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrCannotCancelOrder) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, domain.ErrOrderNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "cancelled"})
}

func (h *OrderHandler) ListOrders(c *gin.Context) {
	minAmount, err := strconv.ParseInt(c.DefaultQuery("min_amount", "0"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid min_amount"})
		return
	}

	maxAmount, err := strconv.ParseInt(c.DefaultQuery("max_amount", "9223372036854775807"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid max_amount"})
		return
	}

	if minAmount > maxAmount {
		c.JSON(http.StatusBadRequest, gin.H{"error": "min_amount cannot be greater than max_amount"})
		return
	}

	orders, err := h.useCase.ListOrdersByAmount(c.Request.Context(), minAmount, maxAmount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, orders)
}
