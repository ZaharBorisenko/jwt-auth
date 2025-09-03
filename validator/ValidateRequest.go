package validator

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ZaharBorisenko/jwt-auth/helpers/JSON"
	"github.com/go-playground/validator/v10"
	"net/http"
)

func ValidateRequest(w http.ResponseWriter, r *http.Request, data interface{}) bool {
	if err := json.NewDecoder(r.Body).Decode(data); err != nil {
		JSON.WriteERROR(w, http.StatusBadRequest, "Invalid JSON")
		return false
	}

	if err := ValidateStruct(data); err != nil {
		handleValidationError(w, err)
		return false
	}
	return true
}

func handleValidationError(w http.ResponseWriter, err error) {
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		e := make(map[string]string)
		for _, fieldError := range validationErrors {
			e[fieldError.Field()] = fmt.Sprintf("Field validation failed: %s", fieldError.Tag())
		}
		JSON.WriteJSON(w, http.StatusBadRequest, map[string]interface{}{
			"error":   "Validation failed",
			"details": e,
		})
	}
}
