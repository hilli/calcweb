package calc

import (
	"fmt"
	"net/http"
	"strconv"
)

func CalculatorHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Could not parse form", http.StatusBadRequest)
		return
	}

	first, _ := strconv.Atoi(r.Form.Get("first"))
	second, _ := strconv.Atoi(r.Form.Get("second"))
	operation := r.Form.Get("operation")

	var result int
	switch operation {
	case "add":
		result = first + second
	case "subtract":
		result = first - second
	case "multiply":
		result = first * second
	case "divide":
		if second != 0 {
			result = first / second
		} else {
			http.Error(w, "Cannot divide by zero", http.StatusBadRequest)
			return
		}
	default:
		http.Error(w, "Invalid operation", http.StatusBadRequest)
		return
	}

	fmt.Fprint(w, result)
}
