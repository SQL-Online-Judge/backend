package restapi

type errorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
