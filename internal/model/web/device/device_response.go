package web

type DeviceResponse struct {
	DeviceID         string `json:"device_id"`
	DeviceName       string `json:"device_name"`
	ManufacturerName string `json:"manufacturer_name"`
	LocalName        string `json:"local_name"`
	UserID           string `json:"user_id"`
}