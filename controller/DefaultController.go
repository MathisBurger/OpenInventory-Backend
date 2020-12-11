package controller

import (
	"fmt"
	"github.com/MathisBurger/OpenInventory-Backend/models"
	"net/http"
)

func DefaultController(writer http.ResponseWriter, request *http.Request) {
	response, err := models.GetJsonResponse("API online", "alert alert-success", "ok", "None", 200)
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(writer, string(response))
}
