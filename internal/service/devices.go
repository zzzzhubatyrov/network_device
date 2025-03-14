package service

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"
	"network/internal/models"
	"network/internal/repository"
	"strconv"
	"time"
)

type DeviceService struct {
	repo *repository.DeviceRepository
}

func NewDeviceService(repo *repository.DeviceRepository) *DeviceService {
	return &DeviceService{
		repo: repo,
	}
}

// generateIP генерирует случайный IP-адрес в сети 192.168.0.0/16
func (s *DeviceService) generateIP() string {
	for i := 0; i < 100; i++ { // Максимум 100 попыток
		// Инициализируем генератор случайных чисел
		rand.Seed(time.Now().UnixNano())

		// Генерируем третий и четвертый октеты
		thirdOctet := rand.Intn(256)  // 0-255
		fourthOctet := rand.Intn(256) // 0-255

		ip := fmt.Sprintf("192.168.%d.%d", thirdOctet, fourthOctet)
		if s.validateIP(ip) {
			return ip
		}
	}
	return "192.168.0.1" // Fallback IP если все попытки неудачны
}

// validateIP проверяет, не занят ли IP-адрес
func (s *DeviceService) validateIP(ip string) bool {
	return !s.repo.IsIPTaken(ip)
}

// getLocalIP получает локальный IP адрес
func (s *DeviceService) getLocalIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}
	return "", fmt.Errorf("local IP not found")
}

// CreateRouter создает роутер с базовым IP
func (s *DeviceService) CreateRouter(req *models.CreateRouterRequest) (*models.Router, error) {
	// Создаем роутер с временным IP
	ip := s.generateIP()

	// Создаем стандартные порты (80 и 443 TCP)
	defaultPorts := []models.Port{
		{
			Number:     80,
			Protocol:   "tcp",
			Status:     "up",
			PortNumber: 1,
		},
		{
			Number:     443,
			Protocol:   "tcp",
			Status:     "down",
			PortNumber: 2,
		},
	}

	// Добавляем дополнительные порты из запроса
	ports := defaultPorts
	for _, portReq := range req.Ports {
		// Проверяем, не дублируется ли порт
		isDuplicate := false
		for _, existingPort := range ports {
			if existingPort.Number == portReq.Number && existingPort.Protocol == portReq.Protocol {
				isDuplicate = true
				break
			}
		}
		if !isDuplicate {
			port := models.Port{
				Number:   portReq.Number,
				Protocol: portReq.Protocol,
				Status:   "down",
			}
			ports = append(ports, port)
		}
	}

	router := &models.Router{
		Name:      req.Name,
		IPAddress: ip,
		Status:    "active",
		Ports:     ports,
		Connected: false,
	}

	if err := s.repo.CreateRouter(router); err != nil {
		return nil, err
	}

	return router, nil
}

func (s *DeviceService) ConnectRouter(req *models.ConnectRouterRequest) (*models.ConnectRouterResponse, error) {
	// Получаем роутер по IP
	router, err := s.repo.GetRouterByIP(req.IPAddress)
	if err != nil {
		return nil, fmt.Errorf("router not found: %w", err)
	}

	// Проверяем, не подключен ли уже роутер
	if router.Connected {
		return nil, fmt.Errorf("router is already connected")
	}

	// Получаем локальный IP
	localIP, err := s.getLocalIP()
	if err != nil {
		return nil, fmt.Errorf("failed to get local IP: %w", err)
	}

	// Подключаем роутер (теперь только меняем флаг connected)
	if err := s.repo.DisconnectRouter(router.ID); err != nil {
		return nil, fmt.Errorf("failed to connect router: %w", err)
	}

	return &models.ConnectRouterResponse{
		RouterID:  router.ID,
		Name:      router.Name,
		IPAddress: router.IPAddress, // Используем IP роутера
		LocalIP:   localIP,          // Локальный IP компьютера
		Status:    router.Status,
		Connected: true,
	}, nil
}

func (s *DeviceService) PingIP(req *models.PingRequest) (*models.PingResult, error) {
	// Создаем результат пинга
	result := &models.PingResult{
		IPAddress: req.IPAddress,
		Status:    "failed",
	}

	// Пингуем IP-адрес
	start := time.Now()
	conn, err := net.DialTimeout("ip:icmp", req.IPAddress, 5*time.Second)
	if err != nil {
		result.Latency = 0
		return result, nil
	}
	defer conn.Close()

	// Вычисляем задержку
	result.Latency = float64(time.Since(start).Milliseconds())
	result.Status = "success"

	return result, nil
}

func (s *DeviceService) SendPacket(req *models.PacketRequest) (*models.PacketResponse, error) {
	// Если source_ip пустой, проверяем подключенный роутер
	if req.SourceIP == "" {
		return nil, fmt.Errorf("source_ip is required, please connect to a router first")
	}

	// Проверяем существование роутера-отправителя
	sourceRouter, err := s.repo.GetRouterByIP(req.SourceIP)
	if err != nil {
		return nil, fmt.Errorf("source router with IP %s not found", req.SourceIP)
	}

	// Проверяем, что роутер подключен
	if !sourceRouter.Connected {
		return nil, fmt.Errorf("source router is not connected")
	}

	// Проверяем существование роутера-получателя
	destRouter, err := s.repo.GetRouterByIP(req.DestinationIP)
	if err != nil {
		return nil, fmt.Errorf("destination router with IP %s not found", req.DestinationIP)
	}

	// Проверяем, открыт ли порт на роутере-получателе
	portFound := false
	for _, port := range destRouter.Ports {
		if port.Number == req.Port && port.Protocol == req.Protocol {
			portFound = true
			if port.Status != "open" {
				return &models.PacketResponse{
					SourceIP:      req.SourceIP,
					DestinationIP: req.DestinationIP,
					Protocol:      req.Protocol,
					Port:          req.Port,
					Status:        "failed",
					Error:         fmt.Sprintf("port %d is %s", req.Port, port.Status),
				}, nil
			}
			break
		}
	}

	if !portFound {
		return &models.PacketResponse{
			SourceIP:      req.SourceIP,
			DestinationIP: req.DestinationIP,
			Protocol:      req.Protocol,
			Port:          req.Port,
			Status:        "failed",
			Error:         fmt.Sprintf("port %d not found", req.Port),
		}, nil
	}

	response := &models.PacketResponse{
		SourceIP:      req.SourceIP,
		DestinationIP: req.DestinationIP,
		Protocol:      req.Protocol,
		Port:          req.Port,
		Status:        "failed",
	}

	// Эмулируем задержку сети (от 10 до 100 мс)
	rand.Seed(time.Now().UnixNano())
	latency := float64(10+rand.Intn(90)) + rand.Float64()
	response.Latency = latency

	// Эмулируем различное поведение для TCP и UDP
	switch req.Protocol {
	case "tcp":
		return s.handleTCPPacket(req, response)
	case "udp":
		return s.handleUDPPacket(req, response)
	default:
		response.Error = "unsupported protocol"
		return response, nil
	}
}

func (s *DeviceService) handleTCPPacket(_ *models.PacketRequest, response *models.PacketResponse) (*models.PacketResponse, error) {
	// Эмулируем TCP соединение (более надежное)
	rand.Seed(time.Now().UnixNano())
	successRate := 0.95 // 95% успешных соединений

	if rand.Float64() < successRate {
		response.Status = "success"
	} else {
		response.Status = "failed"
		response.Error = "connection refused"
		// Увеличиваем латенцию при ошибке
		response.Latency += float64(100 + rand.Intn(200))
	}

	return response, nil
}

func (s *DeviceService) handleUDPPacket(_ *models.PacketRequest, response *models.PacketResponse) (*models.PacketResponse, error) {
	// Эмулируем UDP соединение (менее надежное, но быстрее)
	rand.Seed(time.Now().UnixNano())
	successRate := 0.85 // 85% успешных пакетов

	if rand.Float64() < successRate {
		response.Status = "success"
		// UDP обычно быстрее TCP
		response.Latency *= 0.8
	} else {
		response.Status = "failed"
		response.Error = "packet lost"
	}

	return response, nil
}

func (s *DeviceService) GetAllRouters() ([]models.Router, error) {
	return s.repo.GetAllRouters()
}

func (s *DeviceService) ConfigureRouter(req *models.ConfigureRouterRequest) (*models.ConfigureResponse, error) {
	router, err := s.repo.GetRouterByID(req.RouterID)
	log.Println(router)
	if err != nil {
		return nil, fmt.Errorf("router not found: %w", err)
	}

	updates := make(map[string]interface{})
	if router.Name != "" {
		updates["name"] = req.Name
	}
	if router.Status != "" {
		updates["status"] = req.Status
	}

	if err := s.repo.UpdateRouterConfig(router.ID, updates); err != nil {
		return nil, err
	}

	return &models.ConfigureResponse{
		Success: true,
		Message: "Router configuration updated",
	}, nil
}

// ConfigurePort configures a port on a router
func (s *DeviceService) ConfigurePort(ctx context.Context, req *models.ConfigurePortRequest) error {
	routerID, err := strconv.ParseUint(req.RouterID, 10, 32)
	if err != nil {
		return fmt.Errorf("invalid router ID: %w", err)
	}

	// Validate port number
	if req.PortNumber < 1 || req.PortNumber > 65535 {
		return fmt.Errorf("invalid port number: %d (must be 1-65535)", req.PortNumber)
	}

	// Validate speed
	switch req.Speed {
	case models.SpeedAuto, models.Speed10, models.Speed100, models.Speed1000, models.Speed10000:
		// Valid speed
	default:
		return fmt.Errorf("invalid speed value: %s", req.Speed)
	}

	// Validate duplex mode
	switch req.DuplexMode {
	case models.DuplexModeAuto, models.DuplexModeFull, models.DuplexModeHalf:
		// Valid duplex mode
	default:
		return fmt.Errorf("invalid duplex mode: %s", req.DuplexMode)
	}

	// Validate status
	if req.Status != "up" && req.Status != "down" {
		return fmt.Errorf("invalid status: %s", req.Status)
	}

	// Validate protocol
	if req.Protocol != "" && req.Protocol != "tcp" && req.Protocol != "udp" {
		return fmt.Errorf("invalid protocol: %s", req.Protocol)
	}

	// Get router
	router, err := s.repo.GetRouterByID(uint(routerID))
	if err != nil {
		return fmt.Errorf("failed to get router: %w", err)
	}

	// Check if port exists
	var port *models.Port
	for i := range router.Ports {
		if router.Ports[i].Number == req.PortNumber { // Используем `Number`, а не `PortNumber`
			port = &router.Ports[i]
			break
		}
	}

	// Если порт не найден — создаем новый
	if port == nil {
		newPort := models.Port{
			RouterID:    uint(routerID),
			Number:      req.PortNumber, // Используем `Number`
			Protocol:    req.Protocol,
			Status:      req.Status,
			Speed:       req.Speed,
			DuplexMode:  req.DuplexMode,
			Description: req.Description,
		}
		router.Ports = append(router.Ports, newPort)
	} else {
		// Обновляем существующий порт
		port.Status = req.Status
		port.Protocol = req.Protocol
		port.Speed = req.Speed
		port.DuplexMode = req.DuplexMode
		port.Description = req.Description
	}

	// Save changes
	if err := s.repo.UpdateRouter(router); err != nil {
		return fmt.Errorf("failed to update router: %w", err)
	}

	return nil
}

func (s *DeviceService) CreateConnection(req *models.CreateConnectionRequest) (*models.CreateConnectionResponse, error) {
	// Проверяем существование первого роутера
	routerFrom, err := s.repo.GetRouterByIP(req.RouterFromIP)
	if err != nil {
		return nil, fmt.Errorf("source router not found: %w", err)
	}

	// Проверяем существование второго роутера
	routerTo, err := s.repo.GetRouterByIP(req.RouterToIP)
	if err != nil {
		return nil, fmt.Errorf("destination router not found: %w", err)
	}

	// Проверяем, не существует ли уже соединение между этими роутерами
	if s.repo.ConnectionExists(routerFrom.ID, routerTo.ID) {
		return nil, fmt.Errorf("connection between these routers already exists")
	}

	// Создаем новое соединение
	connection := &models.RouterConnection{
		RouterFromID: routerFrom.ID,
		RouterToID:   routerTo.ID,
		Status:       "active",
		CreatedAt:    time.Now().Format(time.RFC3339),
	}

	// Сохраняем соединение в базе данных
	if err := s.repo.CreateConnection(connection); err != nil {
		return nil, fmt.Errorf("failed to create connection: %w", err)
	}

	return &models.CreateConnectionResponse{
		ID:           connection.ID,
		RouterFromIP: routerFrom.IPAddress,
		RouterToIP:   routerTo.IPAddress,
		Status:       connection.Status,
		CreatedAt:    connection.CreatedAt,
	}, nil
}

func (s *DeviceService) GetAllConnections() ([]models.ConnectionInfo, error) {
	connections, err := s.repo.GetAllConnections()
	if err != nil {
		return nil, fmt.Errorf("failed to get connections: %w", err)
	}

	connectionInfos := make([]models.ConnectionInfo, 0, len(connections))
	for _, conn := range connections {
		// Получаем информацию о роутерах
		routerFrom, err := s.repo.GetRouterByID(conn.RouterFromID)
		if err != nil {
			return nil, fmt.Errorf("failed to get source router: %w", err)
		}

		routerTo, err := s.repo.GetRouterByID(conn.RouterToID)
		if err != nil {
			return nil, fmt.Errorf("failed to get destination router: %w", err)
		}

		connectionInfos = append(connectionInfos, models.ConnectionInfo{
			ID:           conn.ID,
			RouterFromIP: routerFrom.IPAddress,
			RouterToIP:   routerTo.IPAddress,
			Status:       conn.Status,
			CreatedAt:    conn.CreatedAt,
			FromRouter:   *routerFrom,
			ToRouter:     *routerTo,
		})
	}

	return connectionInfos, nil
}

func (s *DeviceService) GetConnectionsByRouterIP(ip string) ([]models.ConnectionInfo, error) {
	connections, err := s.repo.GetConnectionsByRouterIP(ip)
	if err != nil {
		return nil, fmt.Errorf("failed to get connections: %w", err)
	}

	connectionInfos := make([]models.ConnectionInfo, 0, len(connections))
	for _, conn := range connections {
		routerFrom, err := s.repo.GetRouterByID(conn.RouterFromID)
		if err != nil {
			return nil, fmt.Errorf("failed to get source router: %w", err)
		}

		routerTo, err := s.repo.GetRouterByID(conn.RouterToID)
		if err != nil {
			return nil, fmt.Errorf("failed to get destination router: %w", err)
		}

		connectionInfos = append(connectionInfos, models.ConnectionInfo{
			ID:           conn.ID,
			RouterFromIP: routerFrom.IPAddress,
			RouterToIP:   routerTo.IPAddress,
			Status:       conn.Status,
			CreatedAt:    conn.CreatedAt,
			FromRouter:   *routerFrom,
			ToRouter:     *routerTo,
		})
	}

	return connectionInfos, nil
}
