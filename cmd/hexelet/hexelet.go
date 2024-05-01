package hexelet

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"mmskazak/shorturl/cmd/hexelet/handler"
	"mmskazak/shorturl/cmd/hexelet/order"
	"mmskazak/shorturl/cmd/hexelet/storage"
)

func main() {
	orderHandler := &handler.OrderHandler{
		Storage: &storage.OrderStorage{
			Orders: make(map[string]order.Order),
		},
	}

	webApp := fiber.New()
	webApp.Post("/orders", orderHandler.CreateOrder)
	webApp.Get("/orders/:id", orderHandler.GetOrder)

	logrus.Fatal(webApp.Listen(":80"))
}
