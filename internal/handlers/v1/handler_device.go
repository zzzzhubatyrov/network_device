package handlers

import (
	"fmt"
	"network/internal/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) CreateRouter(c *fiber.Ctx) error {
	var req models.CreateRouterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	router, err := h.services.Devices.CreateRouter(&req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(router)
}

func (h *Handler) PingIP(c *fiber.Ctx) error {
	var req models.PingRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.IPAddress == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "IP address is required",
		})
	}

	result, err := h.services.Devices.PingIP(&req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(result)
}

func (h *Handler) SendPacket(c *fiber.Ctx) error {
	var req models.PacketRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Валидация
	if req.SourceIP == "" || req.DestinationIP == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Source and destination IP addresses are required",
		})
	}

	if req.Port < 1 || req.Port > 65535 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid port number",
		})
	}

	result, err := h.services.Devices.SendPacket(&req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(result)
}

func (h *Handler) ConnectRouter(c *fiber.Ctx) error {
	var req models.ConnectRouterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	result, err := h.services.Devices.ConnectRouter(&req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(result)
}

func (h *Handler) GetAllRouters(c *fiber.Ctx) error {
	routers, err := h.services.Devices.GetAllRouters()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(routers)
}

func (h *Handler) ConfigureRouter(c *fiber.Ctx) error {
	var req models.ConfigureRouterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	result, err := h.services.Devices.ConfigureRouter(&req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(result)
}

func (h *Handler) ConfigurePort(c *fiber.Ctx) error {
	var req models.ConfigurePortRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Parse router ID
	routerID, err := strconv.ParseUint(req.RouterID, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid router ID",
		})
	}

	// Update request with parsed ID
	req.RouterID = fmt.Sprintf("%d", routerID)

	// Configure port
	err = h.services.Devices.ConfigurePort(c.Context(), &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Port configured successfully",
	})
}

func (h *Handler) CreateRouterConnection(c *fiber.Ctx) error {
	var req models.CreateConnectionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	response, err := h.services.Devices.CreateConnection(&req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}

func (h *Handler) GetAllConnections(c *fiber.Ctx) error {
	connections, err := h.services.Devices.GetAllConnections()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(connections)
}

func (h *Handler) GetConnectionsByRouterIP(c *fiber.Ctx) error {
	ip := c.Query("ip")
	if ip == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "IP address is required",
		})
	}

	connections, err := h.services.Devices.GetConnectionsByRouterIP(ip)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(connections)
}
