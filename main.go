package main

import (
	"fmt"
	"net/http"
	"webcalc/calc"
)

func main() {
	http.HandleFunc("/calculate", calc.CalculatorHandler)
	http.Handle("/", http.FileServer(http.Dir("./static")))

	fmt.Println("Server is running on at http://localhost:8080")
	err := http.ListenAndServe("localhost:8080", nil)
	if err != nil {
		panic("Failed to start server: " + err.Error())
	}
}
