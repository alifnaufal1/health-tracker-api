package web

type DeviceResponse struct {
	DeviceID         string `json:"device_id"`
	DeviceName       string `json:"device_name"`
	SerialNumber     string `json:"serial_number"`
	FirmwareRevision string `json:"firmware_revision"`
	SoftwareRevision string `json:"software_revision"`
	ManufacturerName string `json:"manufacturer_name"`
	UserID           string `json:"user_id"`
}