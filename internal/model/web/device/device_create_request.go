package web

type DeviceCreateRequest struct {
	DeviceID         string `validate:"required" json:"device_id"`
	DeviceName       string `json:"device_name"`
	ManufacturerName string `json:"manufacturer_name"`
	LocalName        string `json:"serial_number"`
	UserID           string `validate:"required" json:"user_id"`
}