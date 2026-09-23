package device

import (
	"health-tracker-api/pkg/context"
	"health-tracker-api/pkg/response"

	"github.com/gofiber/fiber/v3"
	"github.com/sirupsen/logrus"
)

type DeviceHandler interface {
	Create(c fiber.Ctx) error
	Update(c fiber.Ctx) error
	Delete(c fiber.Ctx) error
	GetByUserID(c fiber.Ctx) error
	GetByID(c fiber.Ctx) error
	GetAll(c fiber.Ctx) error
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

	authUser, err := context.GetAuthUser(c)
	if err != nil {
		return err
	}

	deviceCreateRequest.UserID = authUser.UserID.String()
	
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

	authUser, err := context.GetAuthUser(c)
	if err != nil {
		return err
	}

	deviceUpdateRequest.UserID = authUser.UserID.String()

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

func (h *DeviceHandlerImpl) GetByUserID(c fiber.Ctx) error {
	authUser, err := context.GetAuthUser(c)
	if err != nil {
		return err
	}
	
	user, err := h.DeviceService.GetByUserID(c, authUser.UserID.String())
	if err != nil {
		return err
	}

	return response.Success(c, user, "success find device")
}

func (h *DeviceHandlerImpl) GetByID(c fiber.Ctx) error {	
	user, err := h.DeviceService.GetByID(c, c.Params("id"))
	if err != nil {
		return err
	}

	return response.Success(c, user, "success find device")
}

func (h *DeviceHandlerImpl) GetAll(c fiber.Ctx) error {	
	users, err := h.DeviceService.GetAll(c)
	if err != nil {
		return err
	}

	return response.Success(c, users, "success find all devices")
}
