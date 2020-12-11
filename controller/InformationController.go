package controller

import (
	"fmt"
	"github.com/MathisBurger/OpenInventory-Backend/models"
	"net/http"
)

func InformationController(writer http.ResponseWriter, request *http.Request) {
	response, err := models.GetInformationResponse()
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(writer, string(response))
}
