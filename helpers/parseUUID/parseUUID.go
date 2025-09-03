package parseUUID

import (
	"fmt"
	"github.com/google/uuid"
	"net/http"
)

func ParseUUID(r *http.Request) (uuid.UUID, error) {
	idStr := r.PathValue("id")
	if idStr == "" {
		return uuid.Nil, fmt.Errorf("ID parameter is required")
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid UUID format")
	}

	return id, nil
}
