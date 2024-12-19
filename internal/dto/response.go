package dto

type ResponseError struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}
type Response[T any] struct {
	Response *T             `json:"response"`
	Error    *ResponseError `json:"error"`
}
