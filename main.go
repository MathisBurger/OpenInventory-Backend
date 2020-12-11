package main

import (
	"fmt"
	"github.com/MathisBurger/OpenInventory-Backend/controller"
	"net/http"
)

func main() {
	http.HandleFunc("/", controller.DefaultController)
	http.HandleFunc("/info", controller.InformationController)
	fmt.Println("Server running on port 3000")
	http.ListenAndServe(":3000", nil)
}
