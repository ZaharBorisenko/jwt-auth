package helpers

import (
	"github.com/ZaharBorisenko/jwt-auth/models"
	"net/http"
	"strconv"
)

func ParseQueryParams(r *http.Request) models.ConfigURLParams {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	sort := r.URL.Query().Get("sort")
	order := r.URL.Query().Get("order")

	validSortParams := map[string]bool{
		"first_name": true,
		"created_at": true,
	}

	if !validSortParams[sort] {
		sort = "created_at"
	}
	if order != "asc" && order != "desc" {
		order = "asc"
	}

	offset := (page - 1) * limit

	return models.ConfigURLParams{
		Offset: offset,
		Limit:  limit,
		Sort:   sort,
		Order:  order,
	}
}
