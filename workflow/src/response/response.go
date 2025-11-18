package response

type Common struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

type ErrorDetails struct {
	Code    int         `json:"code" example:"0"`
	Status  string      `json:"status" example:"error"`
	Message string      `json:"message" example:"string"`
	Errors  interface{} `json:"errors,omitempty" swaggertype:"object,string" example:"additionalProp1:{}"`
}
