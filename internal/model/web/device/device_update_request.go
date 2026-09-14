package web

type DeviceUpdateRequest struct {
	DeviceName       string `validate:"required_without_all=LocalName ManufacturerName,omitempty,min=1,max=255" json:"device_name"`
	ManufacturerName string `validate:"required_without_all=LocalName DeviceName,omitempty,min=1,max=255" json:"manufacturer_name"`
	LocalName        string `validate:"required_without_all=ManufacturerName DeviceName,omitempty,min=1,max=255" json:"local_name"`
	UserID           string `json:"user_id"`
}