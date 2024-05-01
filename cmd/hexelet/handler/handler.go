package handler

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"mmskazak/shorturl/cmd/hexelet/order"
)

// OrderCreatorGetter Абстрактное хранилище
type OrderCreatorGetter interface {
	CreateOrder(order order.Order) (string, error)
	GetOrder(id string) (order.Order, error)
}

// Запросы и ответы
type (
	CreateOrderRequest struct {
		UserID     int64   `json:"user_id"`
		ProductIDs []int64 `json:"product_ids"`
	}

	CreateOrderResponse struct {
		ID string `json:"id"`
	}

	GetOrderResponse struct {
		ID         string  `json:"id"`
		UserID     int64   `json:"user_id"`
		ProductIDs []int64 `json:"product_ids"`
	}
)

// OrderHandler Обработчик
type OrderHandler struct {
	Storage OrderCreatorGetter
}

func (h *OrderHandler) CreateOrder(c *fiber.Ctx) error {
	var request CreateOrderRequest
	if err := c.BodyParser(&request); err != nil {
		return fmt.Errorf("body parser: %w", err)
	}

	orderItem := order.Order{
		ID:         uuid.New().String(),
		UserID:     request.UserID,
		ProductIDs: request.ProductIDs,
	}

	id, err := h.Storage.CreateOrder(orderItem)
	if err != nil {
		return fmt.Errorf("create order: %w", err)
	}

	return c.JSON(CreateOrderResponse{
		ID: id,
	})
}

func (h *OrderHandler) GetOrder(c *fiber.Ctx) error {
	id := c.Params("id")

	orderItem, err := h.Storage.GetOrder(id)
	if err != nil {
		return fmt.Errorf("get order: %w", err)
	}

	return c.JSON(GetOrderResponse(orderItem))
}
