package dtos

type Response struct {
	Code    int         `json:"code" example:"200"`
	Success bool        `json:"success" example:"true"`
	Message string      `json:"message,omitempty" example:"request berhasil"`
	Data    interface{} `json:"data,omitempty"`
}

type SuccessResponse struct {
	Code    int         `json:"code" example:"200"`
	Success bool        `json:"success" example:"true"`
	Message string      `json:"message,omitempty" example:"get data success"`
	Data    interface{} `json:"data,omitempty"`
}

type ErrorResponse struct {
	Code    int         `json:"code" example:"400"`
	Success bool        `json:"success" example:"false"`
	Message string      `json:"message,omitempty" example:"error"`
	Data    interface{} `json:"data,omitempty"`
}
