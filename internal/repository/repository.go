package repository

import "gorm.io/gorm"

type Repository struct {
	Devices *DeviceRepository
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		Devices: NewDeviceRepository(db),
	}
}
