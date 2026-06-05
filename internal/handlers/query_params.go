package handlers

import (
	"net/http"
	"strconv"
	"strings"
)

func getQueryParam[T any](r *http.Request, param string, parse func(string) (T, error), defaultValue *T) *T {
	value := r.URL.Query().Get(param)
	if value == "" {
		return defaultValue
	}
	parsed, err := parse(value)
	if err != nil {
		return defaultValue
	}
	return &parsed
}

func parseString(s string) (string, error) {
	return s, nil
}

func parsePositiveInt(s string) (int, error) {
	n, err := strconv.Atoi(s)
	if err != nil || n <= 0 {
		return 0, err
	}
	return n, nil
}

func parsePositiveFloat(s string) (float64, error) {
	n, err := strconv.ParseFloat(s, 64)
	if err != nil || n < 0 {
		return 0, err
	}
	return n, nil
}

func ExtractPageParam(r *http.Request) int {
	defaultPage := 1
	return *getQueryParam(r, "page", parsePositiveInt, &defaultPage)
}

func ExtractLimitParam(r *http.Request) int {
	defaultLimit := 10
	return *getQueryParam(r, "limit", parsePositiveInt, &defaultLimit)
}

func ExtractSortByParam(r *http.Request) string {
	defaultSortBy := "created_at"
	return *getQueryParam(r, "sort_by", parseString, &defaultSortBy)
}

func ExtractOrderParam(r *http.Request) string {
	defaultOrder := "DESC"
	order := *getQueryParam(r, "order", parseString, &defaultOrder)
	return strings.ToUpper(order)
}

func ExtractMinPriceParam(r *http.Request) *float64 {
	return getQueryParam(r, "min_price", parsePositiveFloat, nil)
}

func ExtractMaxPriceParam(r *http.Request) *float64 {
	return getQueryParam(r, "max_price", parsePositiveFloat, nil)
}

func ExtractStartDateParam(r *http.Request) *string {
	return getQueryParam(r, "start", parseString, nil)
}

func ExtractEndDateParam(r *http.Request) *string {
	return getQueryParam(r, "end", parseString, nil)
}

func ExtractOffsetFromPage(page, limit int) int {
	return (page - 1) * limit
}
