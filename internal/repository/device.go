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
	if err := r.db.First(&router, id).Error; err != nil {
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
	return r.db.Save(router).Error
}
