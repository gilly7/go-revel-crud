package entities

type (
	Response struct {
		Success bool        `json:"success"`
		Status  int         `json:"status"`
		Message string      `json:"message"`
		Data    interface{} `json:"data"`
	}
)
