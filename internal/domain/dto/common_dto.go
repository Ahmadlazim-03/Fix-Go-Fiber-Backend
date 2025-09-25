package dto

// Common response structures
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

// Meta represents pagination metadata
type Meta struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int64 `json:"total_pages"`
}

// PaginationQuery represents query parameters for pagination
type PaginationQuery struct {
	Page   int    `query:"page" validate:"omitempty,min=1"`
	Limit  int    `query:"limit" validate:"omitempty,min=1,max=100"`
	Search string `query:"search"`
}

// GetOffset returns the offset for pagination
func (p *PaginationQuery) GetOffset() int {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.Limit <= 0 {
		p.Limit = 10
	}
	return (p.Page - 1) * p.Limit
}

// GetMeta returns pagination metadata
func (p *PaginationQuery) GetMeta(total int64) *Meta {
	if p.Limit <= 0 {
		p.Limit = 10
	}
	totalPages := (total + int64(p.Limit) - 1) / int64(p.Limit)
	return &Meta{
		Page:       p.Page,
		Limit:      p.Limit,
		Total:      total,
		TotalPages: totalPages,
	}
}