package model

// ResponseData is a generic wrapper for successful API responses.
type ResponseData[T any] struct {
	Data    T      `json:"data"`
	Message string `json:"message"`
}

// CreditResponse is returned by the balance endpoint.
type CreditResponse struct {
	Balance  float64 `json:"balance"`
	Currency string  `json:"currency"`
}
