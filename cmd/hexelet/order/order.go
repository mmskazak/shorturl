package order

// Order Модель заказа
type Order struct {
	ID         string
	UserID     int64
	ProductIDs []int64
}
