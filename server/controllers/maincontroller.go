package controllers

import (
	"fmt"
	"net/http"
)

type MainController struct {
}

func (mc *MainController) GetHandler(responseWriter http.ResponseWriter, request *http.Request, parameters map[string]string) {
	responseWriter.Header().Add("Content-Type", "application/json")
	fmt.Fprintf(responseWriter, `{"status":200, "message": "You reached the server"}`)
}

func NewMainController() MainController {
	return MainController{}
}
