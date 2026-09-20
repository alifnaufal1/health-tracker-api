package device

import (
	"errors"
	"fmt"
	"health-tracker-api/pkg/apperror"
	"health-tracker-api/pkg/database"
	"health-tracker-api/pkg/helper"
	"health-tracker-api/pkg/response"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type DeviceService interface {
	Create(c fiber.Ctx, request DeviceCreateRequest) (*DeviceResponse, error)
	Update(c fiber.Ctx, request DeviceUpdateRequest, deviceID string) (*DeviceResponse, error)
	Delete(c fiber.Ctx, deviceID string) error
	FindByID(c fiber.Ctx, deviceID string) (*DeviceResponse, error)
	FindAll(c fiber.Ctx) (*[]DeviceResponse, error)
	IsOwnedByUser(c fiber.Ctx, tx *gorm.DB, deviceID string, userID uuid.UUID) (bool, error)
}

type DeviceServiceImpl struct {
	deviceRepository DeviceRepository
	Validate *validator.Validate
	log *logrus.Entry
}

func NewDeviceService(deviceRepository DeviceRepository, validate *validator.Validate, base *logrus.Logger) DeviceService {
	return &DeviceServiceImpl{
		deviceRepository: deviceRepository,
		Validate: validate,
		log: helper.NewModuleLogger(base, "service", "device") ,
	}
}

func (s *DeviceServiceImpl) Create(c fiber.Ctx, request DeviceCreateRequest) (*DeviceResponse, error) {
	log := helper.LoggerWithRequestID(c, s.log.Logger).WithField("request", request)
	
	log.Debug("received create device request")
	
	err := s.Validate.Struct(request)
	if err != nil {
		log.WithError(err).Warn("failed to validate create device request")
		return nil, &apperror.ValidationError{Message: err.Error()}
	}
	
	tx := database.DB.Begin()
	defer helper.CommitOrRollback(tx)

	var parsedUserID *uuid.UUID = nil
	if request.UserID != "" {
		parsedUUID, err := uuid.Parse(request.UserID)
		if err != nil {
			log.WithError(err).Error("failed to parse user_id")
			return nil, err
		}
		parsedUserID = &parsedUUID
	}
	
	device := &Device{
		DeviceID: request.DeviceID,
		DeviceName: request.DeviceName,
		ManufacturerName: request.ManufacturerName,
		LocalName: request.LocalName,
		UserID: parsedUserID,
	}

	device, err = s.deviceRepository.Save(c, tx, device)
	if err != nil {
		log.WithError(err).Error("failed to save device")
		return nil, fmt.Errorf("save device of user_id %s: %w", request.UserID, err)
	}

	log.WithField("device_id", device.DeviceID).Debug("device created successfully")

	return toDeviceResponse(device), nil
}

func (s *DeviceServiceImpl) Update(c fiber.Ctx, request DeviceUpdateRequest, deviceID string) (*DeviceResponse, error) {
	log := helper.LoggerWithRequestID(c, s.log.Logger).WithField("device_id", deviceID)
	
	log.Debug("received update device request")

	err := s.Validate.Struct(request)
	if err != nil {
		log.WithField("error", err.Error()).Error("Validation failed for update device request")
		return nil, &apperror.ValidationError{Message: err.Error()}
	} 

	tx := database.DB.Begin()
	defer helper.CommitOrRollback(tx)

	device, err := s.deviceRepository.FindById(c, tx, deviceID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn("not found: this device is not found")
			return nil, &apperror.NotFoundError{Message: "device not found"}
		}
		log.WithError(err).Error("failed to find device")
		return nil, err
	}

	if request.DeviceName != "" {
		device.DeviceName = request.DeviceName
	}
	if request.ManufacturerName != "" {
		device.ManufacturerName = request.ManufacturerName
	}
	if request.LocalName != "" {
		device.LocalName = request.LocalName
	}
	if request.UserID != "" {
		parsedUUID, err := uuid.Parse(request.UserID)
		if err != nil {
			log.WithField("error", err.Error()).Warn("Fail to parse user_id")
			return nil, errors.New("user_id is not valid")
		}
		device.UserID = (*uuid.UUID)(&parsedUUID)
	}

	device, err = s.deviceRepository.Update(c, tx, device)
	if err != nil {
		log.WithError(err).Error("failed to update device by device_id")
		return nil, fmt.Errorf("update device by device_id %s: %w", deviceID, err)
	}

	log.WithField("device_id", deviceID).Debug("device updated successfully")

	return toDeviceResponse(device), nil
}

func (s *DeviceServiceImpl) Delete(c fiber.Ctx, deviceID string) error {
	log := helper.LoggerWithRequestID(c, s.log.Logger).WithField("device_id", deviceID)
	
	log.Debug("received delete device request")

	if deviceID == "" {
		log.Warn("validation failed: device_id required")
		return &apperror.ValidationError{Message: "device_id required"}
	}

	tx := database.DB.Begin()
	defer helper.CommitOrRollback(tx)

	_, err := s.deviceRepository.FindById(c, tx, deviceID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn("not found: this device is not found")
			return &apperror.NotFoundError{Message: "device not found"}
		}
		log.WithError(err).Error("failed to find device")
		return err
	}

	err = s.deviceRepository.Delete(c, tx, deviceID)
	if err != nil {
		log.WithError(err).Error("failed to delete device by device_id")
		return fmt.Errorf("delete device by device_id %s: %w", deviceID, err)
	}

	log.WithField("device_id", deviceID).Debug("device deleted successfully")

	return nil
}

func (s *DeviceServiceImpl) FindByID(c fiber.Ctx, deviceID string) (*DeviceResponse, error) {
	log := helper.LoggerWithRequestID(c, s.log.Logger).WithField("device_id", deviceID)
	
	log.Debug("received find device by device_id request")

	if deviceID == "" {
		log.Warn("validation failed: device_id required")
		return nil, &apperror.ValidationError{Message: "device_id required"}
	}

	tx := database.DB.Begin()
	defer helper.CommitOrRollback(tx)

	device, err := s.deviceRepository.FindById(c, tx, deviceID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn("not found: this device is not found")
			return nil, &apperror.NotFoundError{Message: "device not found"}
		}
		log.WithError(err).Error("failed to find device")
		return nil, err
	}

	log.WithField("device_id", device.DeviceID).Debug("device found successfully")

	return toDeviceResponse(device), nil
}

func (s *DeviceServiceImpl) FindAll(c fiber.Ctx) (*[]DeviceResponse, error) {
	log := helper.LoggerWithRequestID(c, s.log.Logger)

	log.Debug("received find all device request")

	tx := database.DB.Begin()
	defer helper.CommitOrRollback(tx)

	devices, err := s.deviceRepository.FindAll(c, tx)
	if err != nil {
		log.WithError(err).Error("failed to find all device")
		return nil, err
	}

	log.WithField("result_count", len(*devices)).Debug("devices found successfully")

	return response.ToResponses(devices, toDeviceResponse), nil
}

func (s *DeviceServiceImpl) IsOwnedByUser(c fiber.Ctx, tx *gorm.DB, deviceID string, userID uuid.UUID) (bool, error) {
	device, err := s.deviceRepository.FindById(c, tx, deviceID)
	if err != nil {
		return false, err
	}
	return *device.UserID == userID, nil
}

func toDeviceResponse(device *Device) *DeviceResponse {
	var userID string
	if device.UserID != nil {
		userID = device.UserID.String()
	}
	return &DeviceResponse{
		DeviceID:   device.DeviceID,
		DeviceName: device.DeviceName,
		LocalName: device.LocalName,
		ManufacturerName: device.ManufacturerName,
		UserID: userID,
	}
}