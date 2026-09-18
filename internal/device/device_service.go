package device

import (
	"errors"
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
}

type DeviceServiceImpl struct {
	deviceRepository DeviceRepository
	Validate *validator.Validate
	log *logrus.Logger
}

func NewDeviceService(deviceRepository DeviceRepository, validate *validator.Validate, log *logrus.Logger) DeviceService {
	return &DeviceServiceImpl{
		deviceRepository: deviceRepository,
		Validate: validate,
		log: log,
	}
}

func (s *DeviceServiceImpl) Create(c fiber.Ctx, request DeviceCreateRequest) (*DeviceResponse, error) {
	log := helper.LoggerWithRequestID(c, s.log)

	log.WithField("request", request).Info("Received create device request")
	
	err := s.Validate.Struct(request)
	if err != nil {
		log.WithField("error", err.Error()).Warn("Validation failed for create device request")
		return nil, &apperror.ValidationError{Message: err.Error()}
	}
	
	tx := database.DB.Begin()
	defer helper.CommitOrRollback(tx)

	var parsedUserID *uuid.UUID = nil
	if request.UserID != "" {
		parsedUUID, err := uuid.Parse(request.UserID)
		
		if err != nil {
			log.WithField("error", err.Error()).Warn("Fail to parse user_id")
			return nil, errors.New("user_id is not valid")
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
		log.Warn("Create device process failed at repository layer")
		return nil, errors.New(err.Error())
	}

	log.WithField("device_id", device.DeviceID).Info("device created successfully")

	return toDeviceResponse(device), nil
}

func (s *DeviceServiceImpl) Update(c fiber.Ctx, request DeviceUpdateRequest, deviceID string) (*DeviceResponse, error) {
	log := helper.LoggerWithRequestID(c, s.log).WithField("device_id", deviceID)
	
	log.Info("Received update device request")

	err := s.Validate.Struct(request)
	if err != nil {
		log.WithField("error", err.Error()).Warn("Validation failed for update device request")
		return nil, &apperror.ValidationError{Message: err.Error()}
	}

	tx := database.DB.Begin()
	defer helper.CommitOrRollback(tx)

	device, err := s.deviceRepository.FindById(c, tx, deviceID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &apperror.NotFoundError{Message: "device not found"}
		}
		return nil, errors.New(err.Error())
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
		log.Warn("Update device process failed at repository layer")
		return nil, errors.New(err.Error())
	}

	log.WithField("device_id", device.DeviceID).Info("device updated successfully")

	return toDeviceResponse(device), nil
}

func (s *DeviceServiceImpl) Delete(c fiber.Ctx, deviceID string) error {
	log := helper.LoggerWithRequestID(c, s.log).WithField("device_id", deviceID)
	
	log.Info("Received delete device request")

	if deviceID == "" {
		log.Warn("Validation failed for delete device request")
		return &apperror.ValidationError{Message: "device_id required"}
	}

	tx := database.DB.Begin()
	defer helper.CommitOrRollback(tx)

	device, err := s.deviceRepository.FindById(c, tx, deviceID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &apperror.NotFoundError{Message: "device not found"}
		}
		return errors.New(err.Error())
	}

	err = s.deviceRepository.Delete(c, tx, deviceID)
	if err != nil {
		log.Warn("Delete device process failed at repository layer")
		return errors.New(err.Error())
	}

	log.WithField("device_id", device.DeviceID).Info("device deleted successfully")

	return nil
}

func (s *DeviceServiceImpl) FindByID(c fiber.Ctx, deviceID string) (*DeviceResponse, error) {
	log := helper.LoggerWithRequestID(c, s.log).WithField("device_id", deviceID)
	
	log.Info("Received find device by id request")

	if deviceID == "" {
		log.Warn("Validation failed for find device by id request")
		return nil, &apperror.ValidationError{Message: "device_id required"}
	}

	tx := database.DB.Begin()
	defer helper.CommitOrRollback(tx)

	device, err := s.deviceRepository.FindById(c, tx, deviceID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &apperror.NotFoundError{Message: "device not found"}
		}
		return nil, errors.New(err.Error())
	}

	log.WithField("device_id", device.DeviceID).Info("device found successfully")

	return toDeviceResponse(device), nil
}

func (s *DeviceServiceImpl) FindAll(c fiber.Ctx) (*[]DeviceResponse, error) {
	log := helper.LoggerWithRequestID(c, s.log)

	log.Info("Received find all device request")

	tx := database.DB.Begin()
	defer helper.CommitOrRollback(tx)

	devices, err := s.deviceRepository.FindAll(c, tx)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	log.Info("devices found successfully")

	return response.ToResponses(devices, toDeviceResponse), nil
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