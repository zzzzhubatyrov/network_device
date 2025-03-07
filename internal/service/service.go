package service

import "network/internal/repository"

type Service struct {
	Devices *DeviceService
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Devices: NewDeviceService(repos.Devices),
	}
}
