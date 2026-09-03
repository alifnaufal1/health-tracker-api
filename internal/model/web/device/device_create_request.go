package web

type DeviceCreateRequest struct {
	DeviceName       string `validate:"required,min=1,max=255" json:"device_name"`
	SerialNumber     string `validate:"required,min=1,max=255" json:"serial_number"`
	FirmwareRevision string `validate:"required,min=1,max=255" json:"firmware_revision"`
	SoftwareRevision string `validate:"required,min=1,max=255" json:"software_revision"`
	ManufacturerName string `validate:"required,min=1,max=255" json:"manufacturer_name"`
	UserID           string `validate:"min=1,max=255" json:"user_id"`
}