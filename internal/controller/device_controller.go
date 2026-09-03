package controller

import (
	"health-tracker-api/internal/helper"
	deviceWeb "health-tracker-api/internal/model/web/device"
	"health-tracker-api/internal/service"

	"github.com/gofiber/fiber/v3"
	"github.com/sirupsen/logrus"
)

type DeviceController interface {
	Create(ctx fiber.Ctx) error
	Update(ctx fiber.Ctx) error
	Delete(ctx fiber.Ctx) error
	FindByID(ctx fiber.Ctx) error
	FindAll(ctx fiber.Ctx) error
}

type DeviceControllerImpl struct {
	DeviceService service.DeviceService
	log *logrus.Logger
}

func NewDeviceController(deviceService service.DeviceService, log *logrus.Logger) DeviceController {
	return &DeviceControllerImpl{
		DeviceService: deviceService,
		log: log,
	}
}

func (c *DeviceControllerImpl) Create(ctx fiber.Ctx) error {
	deviceCreateRequest := new(deviceWeb.DeviceCreateRequest)
	err := ctx.Bind().Body(deviceCreateRequest)
	if err != nil {
		return err
	}
	
	createdDevice, err := c.DeviceService.Create(ctx, *deviceCreateRequest)
	if err != nil {
		return err
	}

	return helper.ToWebResponse(ctx, createdDevice)
}

func (c *DeviceControllerImpl) Update(ctx fiber.Ctx) error {
	deviceUpdateRequest := new(deviceWeb.DeviceUpdateRequest)
	err := ctx.Bind().Body(deviceUpdateRequest)
	if err != nil {
		return err
	}

	updatedUser, err := c.DeviceService.Update(ctx, *deviceUpdateRequest, ctx.Params("id"))
	if err != nil {
		return err
	}

	return helper.ToWebResponse(ctx, updatedUser)
}

func (c *DeviceControllerImpl) Delete(ctx fiber.Ctx) error {
	err := c.DeviceService.Delete(ctx, ctx.Params("id"))
	if err != nil {
		return err
	}

	return helper.ToWebResponse(ctx, nil)
}

func (c *DeviceControllerImpl) FindByID(ctx fiber.Ctx) error {	
	user, err := c.DeviceService.FindByID(ctx, ctx.Params("id"))
	if err != nil {
		return err
	}

	return helper.ToWebResponse(ctx, user)
}

func (c *DeviceControllerImpl) FindAll(ctx fiber.Ctx) error {	
	users, err := c.DeviceService.FindAll(ctx)
	if err != nil {
		return err
	}

	return helper.ToWebResponse(ctx, users)
}
