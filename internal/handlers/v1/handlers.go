package handlers

import (
	"network/internal/service"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{
		services: services,
	}
}

func (h *Handler) InitRoute(app *fiber.App) fiber.Handler {
	api := app.Group("/api/v1")

	api.Post("/routers", h.CreateRouter)
	api.Get("/routers", h.GetAllRouters)

	api.Post("/routers/connect", h.ConnectRouter)
	api.Post("/ping", h.PingIP)
	api.Post("/packet", h.SendPacket)

	api.Post("/routers/configure", h.ConfigureRouter)
	api.Post("/ports/configure", h.ConfigurePort)

	return nil
}
