package device

type DeviceCreateRequest struct {
	DeviceID         string `validate:"required" json:"device_id"`
	DeviceName       string `json:"device_name"`
	ManufacturerName string `json:"manufacturer_name"`
	LocalName        string `json:"serial_number"`
	UserID           string `validate:"required" json:"user_id"`
}

type DeviceUpdateRequest struct {
	DeviceName       string `validate:"required_without_all=LocalName ManufacturerName,omitempty,min=1,max=255" json:"device_name"`
	ManufacturerName string `validate:"required_without_all=LocalName DeviceName,omitempty,min=1,max=255" json:"manufacturer_name"`
	LocalName        string `validate:"required_without_all=ManufacturerName DeviceName,omitempty,min=1,max=255" json:"local_name"`
	UserID           string `json:"user_id"`
}

type DeviceResponse struct {
	DeviceID         string `json:"device_id"`
	DeviceName       string `json:"device_name"`
	ManufacturerName string `json:"manufacturer_name"`
	LocalName        string `json:"local_name"`
	UserID           string `json:"user_id"`
}