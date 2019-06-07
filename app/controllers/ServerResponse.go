package controllers

import "fmt"

// ServerResponse структура ответа сервера на GET запрос
type ServerResponse struct {
	Status       string
	ErrorMessage string
	Data         interface{}
}

// Succes получение структуры ответа, при успешном запросе
func Succes(data interface{}) *ServerResponse {
	response := new(ServerResponse)

	response.Status = "Succes"
	response.Data = data

	return response
}

// Failed получение структуры ответа, при успешном запросе
func Failed(err error) *ServerResponse {
	response := new(ServerResponse)

	fmt.Printf("\n\n%s\n\n", err.Error())

	response.Status = "Failed"
	response.ErrorMessage = err.Error()

	return response
}
