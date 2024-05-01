package storage

import (
	"errors"
	"mmskazak/shorturl/cmd/hexelet/order"
)

// OrderCreatorGetter Абстрактное хранилище
type OrderCreatorGetter interface {
	CreateOrder(order order.Order) (string, error)
	GetOrder(id string) (order.Order, error)
}

// OrderStorage Хранилище
type OrderStorage struct {
	Orders map[string]order.Order
}

func (o *OrderStorage) CreateOrder(order order.Order) (string, error) {
	o.Orders[order.ID] = order

	return order.ID, nil
}

// Ошибки
var (
	errOrderNotFound = errors.New("order not found")
)

func (o *OrderStorage) GetOrder(id string) (order.Order, error) {
	orderItem, ok := o.Orders[id]
	if !ok {
		return order.Order{}, errOrderNotFound
	}

	return orderItem, nil
}
