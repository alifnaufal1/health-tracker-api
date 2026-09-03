package repository

import (
	"health-tracker-api/internal/helper"
	"health-tracker-api/internal/model/domain"

	"github.com/gofiber/fiber/v3"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type DeviceRepository interface {
	Save(ctx fiber.Ctx, tx *gorm.DB, device *domain.Device) (*domain.Device, error)
	Update(ctx fiber.Ctx, tx *gorm.DB, device *domain.Device) (*domain.Device, error)
	Delete(ctx fiber.Ctx, tx *gorm.DB, deviceId string) error
	FindById(ctx fiber.Ctx, tx *gorm.DB, deviceId string) (*domain.Device, error)
	FindAll(ctx fiber.Ctx, tx *gorm.DB) (*[]domain.Device, error)
}

type DeviceRepositoryImpl struct {
	log *logrus.Logger
}

func NewDeviceRepository(log *logrus.Logger) DeviceRepository {
	return &DeviceRepositoryImpl{log: log}
}

func (r *DeviceRepositoryImpl) Save(c fiber.Ctx, tx *gorm.DB, device *domain.Device) (*domain.Device, error) {
	log := helper.LoggerWithRequestID(c, r.log).WithField("device_name", device.DeviceName)

	log.Debug("Inserting new device into database")
	
	result := gorm.WithResult()
	err := gorm.G[domain.Device](tx, result).Create(c, device)
	if err != nil {
		log.WithField("error", err.Error()).Error("Database insert failed")
		return nil, err
	}
	
	log.WithField("device_id", device.ID).Debug("Device inserted successfully")
	
	return device, nil
}

func (r *DeviceRepositoryImpl) Update(c fiber.Ctx, tx *gorm.DB, device *domain.Device) (*domain.Device, error) {
	log := helper.LoggerWithRequestID(c, r.log).WithField("device_id", device.ID)

	log.Debug("Updating device into database")
	
	result := gorm.WithResult()
	_, err := gorm.G[domain.Device](tx, result).Where("id = ?", device.ID).Updates(c, *device)
	if err != nil {
		log.WithField("error", err.Error()).Error("Database update failed")
		return nil, err
	}
	
	log.WithField("device_id", device.ID).Debug("Device updated successfully")
	
	return device, nil
}

func (r *DeviceRepositoryImpl) Delete(c fiber.Ctx, tx *gorm.DB, deviceID string) error {
	log := helper.LoggerWithRequestID(c, r.log).WithField("device_id", deviceID)
	
	log.Debug("Deleting device into database")
	
	result := gorm.WithResult()
	_, err := gorm.G[domain.Device](tx, result).Where("id = ?", deviceID).Delete(c)
	if err != nil {
		log.WithField("error", err.Error()).Error("Database delete failed")
		return err
	}
	
	log.WithField("device_id", deviceID).Debug("Device deleted successfully")
	
	return nil
}

func (r *DeviceRepositoryImpl) FindById(c fiber.Ctx, tx *gorm.DB, deviceID string) (*domain.Device, error) {
	log := helper.LoggerWithRequestID(c, r.log).WithField("device_id", deviceID)

	log.Debug("Finding device by id from database")
	
	result := gorm.WithResult()
	devices, err := gorm.G[domain.Device](tx, result).Select("id", "device_name", "serial_number", "firmware", "software_revision", "manufacturer_name").Where("id = ?", deviceID).Find(c)
	if err != nil {
		log.WithField("error", err.Error()).Error("Database find by id failed")
		return nil, err
	}
	
	log.WithField("device_id", deviceID).Debug("Device found by id successfully")
	
	return &devices[0], nil
}

func (r *DeviceRepositoryImpl) FindAll(c fiber.Ctx, tx *gorm.DB) (*[]domain.Device, error) {
	log := helper.LoggerWithRequestID(c, r.log)

	log.Debug("Finding all devices from database")
	
	result := gorm.WithResult()
	devices, err := gorm.G[domain.Device](tx, result).Select("id", "device_name", "serial_number", "firmware", "software_revision", "manufacturer_name").Find(c)
	if err != nil {
		log.WithField("error", err.Error()).Error("Database find all failed")
		return nil, err
	}
	
	log.Debug("All device found successfully")
	
	return &devices, nil
}
