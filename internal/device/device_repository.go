package device

import (
	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

type DeviceRepository interface {
	Save(c fiber.Ctx, tx *gorm.DB, device *Device) (*Device, error)
	Update(c fiber.Ctx, tx *gorm.DB, device *Device) (*Device, error)
	Delete(c fiber.Ctx, tx *gorm.DB, deviceId string) error
	FindById(c fiber.Ctx, tx *gorm.DB, deviceId string) (*Device, error)
	FindAll(c fiber.Ctx, tx *gorm.DB) (*[]Device, error)
}

type DeviceRepositoryImpl struct {}

func NewDeviceRepository() DeviceRepository {
	return &DeviceRepositoryImpl{}
}

func (r *DeviceRepositoryImpl) Save(c fiber.Ctx, tx *gorm.DB, device *Device) (*Device, error) {
	result := gorm.WithResult()
	err := gorm.G[Device](tx, result).Create(c, device)
	if err != nil {
		return nil, err
	}
	return device, nil
}

func (r *DeviceRepositoryImpl) Update(c fiber.Ctx, tx *gorm.DB, device *Device) (*Device, error) {
	result := gorm.WithResult()
	_, err := gorm.G[Device](tx, result).Where("device_id = ?", device.DeviceID).Updates(c, *device)
	if err != nil {
		return nil, err
	}
	return device, nil
}

func (r *DeviceRepositoryImpl) Delete(c fiber.Ctx, tx *gorm.DB, deviceID string) error {
	result := gorm.WithResult()
	_, err := gorm.G[Device](tx, result).Where("device_id = ?", deviceID).Delete(c)
	if err != nil {
		return err
	}
	return nil
}

func (r *DeviceRepositoryImpl) FindById(c fiber.Ctx, tx *gorm.DB, deviceID string) (*Device, error) {
	result := gorm.WithResult()
	device, err := gorm.G[Device](tx, result).Where("device_id = ?", deviceID).First(c)
	if err != nil {
		return nil, err
	}
	return &device, nil
}

func (r *DeviceRepositoryImpl) FindAll(c fiber.Ctx, tx *gorm.DB) (*[]Device, error) {
	result := gorm.WithResult()
	devices, err := gorm.G[Device](tx, result).Find(c)
	if err != nil {
		return nil, err
	}
	return &devices, nil
}
