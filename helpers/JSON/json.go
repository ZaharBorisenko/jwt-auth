package JSON

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Error struct {
	Error string
}

func WriteJSON(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		fmt.Printf("could not write json %v", err)
	}
}

func WriteERROR(w http.ResponseWriter, statusCode int, message string) {
	WriteJSON(w, statusCode, Error{Error: message})
}
