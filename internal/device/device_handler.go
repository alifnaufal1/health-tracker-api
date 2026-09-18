package device

import (
	"health-tracker-api/pkg/response"

	"github.com/gofiber/fiber/v3"
	"github.com/sirupsen/logrus"
)

type DeviceHandler interface {
	Create(c fiber.Ctx) error
	Update(c fiber.Ctx) error
	Delete(c fiber.Ctx) error
	FindByID(c fiber.Ctx) error
	FindAll(c fiber.Ctx) error
}

type DeviceHandlerImpl struct {
	DeviceService DeviceService
	log *logrus.Logger
}

func NewDeviceHandler(deviceService DeviceService, log *logrus.Logger) DeviceHandler {
	return &DeviceHandlerImpl{
		DeviceService: deviceService,
		log: log,
	}
}

func (h *DeviceHandlerImpl) Create(c fiber.Ctx) error {
	deviceCreateRequest := new(DeviceCreateRequest)
	err := c.Bind().Body(deviceCreateRequest)
	if err != nil {
		return err
	}
	
	createdDevice, err := h.DeviceService.Create(c, *deviceCreateRequest)
	if err != nil {
		return err
	}

	return response.Success(c, createdDevice, "success create new device")
}

func (h *DeviceHandlerImpl) Update(c fiber.Ctx) error {
	deviceUpdateRequest := new(DeviceUpdateRequest)
	err := c.Bind().Body(deviceUpdateRequest)
	if err != nil {
		return err
	}

	updatedUser, err := h.DeviceService.Update(c, *deviceUpdateRequest, c.Params("id"))
	if err != nil {
		return err
	}

	return response.Success(c, updatedUser, "success update device")
}

func (h *DeviceHandlerImpl) Delete(c fiber.Ctx) error {
	err := h.DeviceService.Delete(c, c.Params("id"))
	if err != nil {
		return err
	}

	return response.Success(c, nil, "success delete device")
}

func (h *DeviceHandlerImpl) FindByID(c fiber.Ctx) error {	
	user, err := h.DeviceService.FindByID(c, c.Params("id"))
	if err != nil {
		return err
	}

	return response.Success(c, user, "success find device")
}

func (h *DeviceHandlerImpl) FindAll(c fiber.Ctx) error {	
	users, err := h.DeviceService.FindAll(c)
	if err != nil {
		return err
	}

	return response.Success(c, users, "success find all devices")
}
