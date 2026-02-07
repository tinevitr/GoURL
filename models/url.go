package models

// RedirectRequest tetap di models untuk Gin binding
type RedirectRequest struct {
	Slug string `uri:"slug" binding:"required"`
}
