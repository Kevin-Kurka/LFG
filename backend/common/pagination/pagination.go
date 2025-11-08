package pagination

import (
	"fmt"
	"net/http"
	"strconv"
)

// Params holds pagination parameters
type Params struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
	Offset   int `json:"-"`
}

// Response holds pagination metadata for responses
type Response struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalItems int64 `json:"total_items"`
	TotalPages int   `json:"total_pages"`
}

const (
	DefaultPage     = 1
	DefaultPageSize = 20
	MaxPageSize     = 100
)

// ParseParams extracts pagination parameters from HTTP request
func ParseParams(r *http.Request) Params {
	page := parseIntParam(r.URL.Query().Get("page"), DefaultPage)
	pageSize := parseIntParam(r.URL.Query().Get("page_size"), DefaultPageSize)

	// Validate bounds
	if page < 1 {
		page = DefaultPage
	}
	if pageSize < 1 {
		pageSize = DefaultPageSize
	}
	if pageSize > MaxPageSize {
		pageSize = MaxPageSize
	}

	offset := (page - 1) * pageSize

	return Params{
		Page:     page,
		PageSize: pageSize,
		Offset:   offset,
	}
}

// NewResponse creates a pagination response
func NewResponse(params Params, totalItems int64) Response {
	totalPages := 0
	if params.PageSize > 0 {
		totalPages = int((totalItems + int64(params.PageSize) - 1) / int64(params.PageSize))
	}

	return Response{
		Page:       params.Page,
		PageSize:   params.PageSize,
		TotalItems: totalItems,
		TotalPages: totalPages,
	}
}

// GetSQLLimit returns SQL LIMIT clause
func (p Params) GetSQLLimit() string {
	return fmt.Sprintf("LIMIT %d OFFSET %d", p.PageSize, p.Offset)
}

// parseIntParam parses integer query parameter with default fallback
func parseIntParam(param string, defaultValue int) int {
	if param == "" {
		return defaultValue
	}
	val, err := strconv.Atoi(param)
	if err != nil {
		return defaultValue
	}
	return val
}

// PaginatedResponse wraps data with pagination metadata
type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Pagination Response    `json:"pagination"`
}

// NewPaginatedResponse creates a new paginated response
func NewPaginatedResponse(data interface{}, params Params, totalItems int64) PaginatedResponse {
	return PaginatedResponse{
		Data:       data,
		Pagination: NewResponse(params, totalItems),
	}
}
