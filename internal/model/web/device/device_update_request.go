package web

type DeviceUpdateRequest struct {
	DeviceName       string `validate:"required_without_all=SerialNumber FirmwareRevision SoftwareRevision ManufacturerName,omitempty,min=1,max=255" json:"device_name"`
	SerialNumber     string `validate:"required_without_all=FirmwareRevision SoftwareRevision ManufacturerName DeviceName,omitempty,min=1,max=255" json:"serial_number"`
	FirmwareRevision string `validate:"required_without_all=SerialNumber SoftwareRevision ManufacturerName DeviceName,omitempty,min=1,max=255" json:"firmware_revision"`
	SoftwareRevision string `validate:"required_without_all=SerialNumber FirmwareRevision ManufacturerName DeviceName,omitempty,min=1,max=255" json:"software_revision"`
	ManufacturerName string `validate:"required_without_all=SerialNumber FirmwareRevision SoftwareRevision DeviceName,omitempty,min=1,max=255" json:"manufacturer_name"`
	UserID           string `json:"user_id"`
}