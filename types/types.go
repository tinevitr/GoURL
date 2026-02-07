package types

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

type Meta struct {
	Total  int `json:"total,omitempty"`
	Limit  int `json:"limit,omitempty"`
	Offset int `json:"offset,omitempty"`
}

type URLResponse struct {
	OriginalURL  string `json:"original_url"`
	ShortURL     string `json:"short_url"`
	Slug         string `json:"slug"`
	CreatedAt    int64  `json:"created_at"`
	ExpiresAt    int64  `json:"expires_at,omitempty"`
	ClickCount   int64  `json:"click_count"`
	IsCustomSlug bool   `json:"is_custom_slug"`
	LastAccessed int64  `json:"last_accessed,omitempty"`
}

type StatsResponse struct {
	Slug         string `json:"slug"`
	OriginalURL  string `json:"original_url"`
	ShortURL     string `json:"short_url"`
	ClickCount   int64  `json:"click_count"`
	CreatedAt    int64  `json:"created_at"`
	LastAccessed int64  `json:"last_accessed"`
}

type CreateURLRequest struct {
	URL   string `json:"url" binding:"required"`
	Slug  string `json:"slug,omitempty"`
	Title string `json:"title,omitempty"`
}

func SuccessResponse(data interface{}) Response {
	return Response{
		Success: true,
		Data:    data,
	}
}

func SuccessResponseWithMessage(data interface{}, message string) Response {
	return Response{
		Success: true,
		Message: message,
		Data:    data,
	}
}

func ErrorResponse(message string) Response {
	return Response{
		Success: false,
		Error:   message,
	}
}

func ErrorResponseWithCode(message string, code int) Response {
	return Response{
		Success: false,
		Error:   message,
		Meta: &Meta{
			Total: code,
		},
	}
}
