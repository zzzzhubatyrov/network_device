package repository

import (
	"fmt"
	"network/internal/models"

	"gorm.io/gorm"
)

type DeviceRepository struct {
	db *gorm.DB
}

func NewDeviceRepository(db *gorm.DB) *DeviceRepository {
	return &DeviceRepository{
		db: db,
	}
}

func (r *DeviceRepository) CreateRouter(router *models.Router) error {
	return r.db.Create(router).Error
}

func (r *DeviceRepository) GetRouterByID(id uint) (*models.Router, error) {
	var router models.Router
	if err := r.db.Preload("Ports").First(&router, id).Error; err != nil {
		return nil, err
	}
	return &router, nil
}

func (r *DeviceRepository) GetRouterByIP(ip string) (*models.Router, error) {
	var router models.Router
	if err := r.db.Preload("Ports").Where("ip_address = ?", ip).First(&router).Error; err != nil {
		return nil, err
	}
	return &router, nil
}

func (r *DeviceRepository) IsIPTaken(ip string) bool {
	var count int64
	r.db.Model(&models.Router{}).Where("ip_address = ?", ip).Count(&count)
	return count > 0
}

func (r *DeviceRepository) UpdatePortStatus(routerID uint, portNumber int, status string) error {
	router, err := r.GetRouterByID(routerID)
	if err != nil {
		return err
	}

	for i, port := range router.Ports {
		if port.PortNumber == portNumber {
			router.Ports[i].Status = status
			return r.UpdateRouter(router)
		}
	}

	return fmt.Errorf("port %d not found", portNumber)
}

func (r *DeviceRepository) ConnectRouter(routerID uint) error {
	return r.db.Model(&models.Router{}).
		Where("id = ?", routerID).
		Update("connected", true).Error
}

func (r *DeviceRepository) DisconnectRouter(routerID uint) error {
	return r.db.Model(&models.Router{}).
		Where("id = ?", routerID).
		Update("connected", false).Error
}

func (r *DeviceRepository) GetAllRouters() ([]models.Router, error) {
	var routers []models.Router
	if err := r.db.Preload("Ports").Find(&routers).Error; err != nil {
		return nil, err
	}
	return routers, nil
}

func (r *DeviceRepository) UpdateRouterConfig(routerID uint, updates map[string]interface{}) error {
	return r.db.Model(&models.Router{}).
		Where("id = ?", routerID).
		Updates(updates).Error
}

func (r *DeviceRepository) UpdateRouter(router *models.Router) error {
	return r.db.Updates(router).Error
}

func (r *DeviceRepository) ConnectionExists(routerFromID, routerToID uint) bool {
	var count int64
	r.db.Model(&models.RouterConnection{}).
		Where("(router_from_id = ? AND router_to_id = ?) OR (router_from_id = ? AND router_to_id = ?)",
			routerFromID, routerToID, routerToID, routerFromID).
		Count(&count)
	return count > 0
}

func (r *DeviceRepository) CreateConnection(connection *models.RouterConnection) error {
	return r.db.Create(connection).Error
}

func (r *DeviceRepository) GetAllConnections() ([]models.RouterConnection, error) {
	var connections []models.RouterConnection
	err := r.db.Find(&connections).Error
	return connections, err
}

func (r *DeviceRepository) GetConnectionByID(id uint) (*models.RouterConnection, error) {
	var connection models.RouterConnection
	err := r.db.First(&connection, id).Error
	return &connection, err
}

func (r *DeviceRepository) GetConnectionsByRouterIP(ip string) ([]models.RouterConnection, error) {
	var router models.Router
	if err := r.db.Where("ip_address = ?", ip).First(&router).Error; err != nil {
		return nil, err
	}

	var connections []models.RouterConnection
	err := r.db.Where("router_from_id = ? OR router_to_id = ?", router.ID, router.ID).Find(&connections).Error
	return connections, err
}
