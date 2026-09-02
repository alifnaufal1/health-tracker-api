package web

type WebResponse struct {
	Code    int    `json:"code"`
	Status  bool   `json:"status"`
	Data    any    `json:"data"`
	Message string `json:"message"`
}