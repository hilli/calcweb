package main

import (
	"fmt"
	"net/http"
	"webcalc/calc"
)

func main() {
	http.HandleFunc("/calculate", calc.CalculatorHandler)
	http.Handle("/", http.FileServer(http.Dir("./static")))

	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", nil)
}
