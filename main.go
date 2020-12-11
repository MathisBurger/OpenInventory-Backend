package main

import (
	"fmt"
	"github.com/MathisBurger/OpenInventory-Backend/controller"
	"github.com/MathisBurger/OpenInventory-Backend/installation"
	"net/http"
)

func main() {
	if installation.Install() {
		http.HandleFunc("/", controller.DefaultController)
		http.HandleFunc("/info", controller.InformationController)
		fmt.Println("Server running on port 3000")
		http.ListenAndServe(":3000", nil)
	} else {
		fmt.Println("Please fix errors first to launch webserver")
	}
}
