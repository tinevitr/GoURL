package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
	"gourl/db"
	"gourl/config"
	"gourl/models"
	"gourl/types"
	"gourl/utils"

	"github.com/gin-gonic/gin"
	"github.com/asaskevich/govalidator"
	"github.com/go-redis/redis/v8"
)

type URLHandler struct {
	redisClient *database.RedisClient
	ctx         context.Context
}

func NewURLHandler(redisClient *database.RedisClient) *URLHandler {
	return &URLHandler{
		redisClient: redisClient,
		ctx:         context.Background(),
	}
}

// createShortURL adalah fungsi internal untuk membuat short URL
func (h *URLHandler) createShortURL(originalURL, customSlug string) (*types.URLResponse, error) {
	// Validate URL
	if !govalidator.IsURL(originalURL) {
		return nil, fmt.Errorf("invalid URL format")
	}

	// Generate or validate slug
	var slug string
	var isCustomSlug bool

	if customSlug != "" {
		// Use custom slug
		if !utils.IsValidSlug(customSlug) {
			return nil, fmt.Errorf("slug must be 3-20 characters and contain only letters, numbers, hyphens, or underscores")
		}

		// Check if slug already exists
		exists, err := h.redisClient.SlugExists(customSlug)
		if err != nil {
			return nil, fmt.Errorf("failed to check slug availability")
		}
		if exists {
			return nil, fmt.Errorf("slug already exists")
		}

		slug = customSlug
		isCustomSlug = true
	} else {
		// Generate unique slug
		for i := 0; i < 10; i++ { // Try 10 times
			slug = utils.GenerateSlug()
			exists, err := h.redisClient.SlugExists(slug)
			if err != nil {
				return nil, fmt.Errorf("failed to generate slug")
			}
			if !exists {
				break
			}
			if i == 9 {
				return nil, fmt.Errorf("failed to generate unique slug after 10 attempts")
			}
		}
		isCustomSlug = false
	}

	// Save to Redis
	err := h.redisClient.SaveURL(slug, originalURL)
	if err != nil {
		return nil, fmt.Errorf("failed to save URL")
	}

	// Create response
	cfg := config.LoadConfig()
	response := &types.URLResponse{
		OriginalURL:  originalURL,
		ShortURL:     cfg.BaseURL + "/" + slug,
		Slug:         slug,
		CreatedAt:    time.Now().Unix(),
		ExpiresAt:    time.Now().Add(config.URLExpiration).Unix(),
		ClickCount:   0,
		IsCustomSlug: isCustomSlug,
		LastAccessed: time.Now().Unix(),
	}

	return response, nil
}

// POST /api/create
func (h *URLHandler) CreateShortURL(c *gin.Context) {
	var req types.CreateURLRequest
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse("Invalid request body"))
		return
	}

	response, err := h.createShortURL(req.URL, req.Slug)
	if err != nil {
		statusCode := http.StatusBadRequest
		if err.Error() == "slug already exists" {
			statusCode = http.StatusConflict
		} else if err.Error() == "failed to save URL" || err.Error() == "failed to generate slug" {
			statusCode = http.StatusInternalServerError
		}
		c.JSON(statusCode, types.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, types.SuccessResponse(response))
}

// GET /api/create?url={url}&slug={slug} (optional)
func (h *URLHandler) CreateShortURLViaGet(c *gin.Context) {
	originalURL := c.Query("url")
	if originalURL == "" {
		c.JSON(http.StatusBadRequest, types.ErrorResponse("URL parameter is required"))
		return
	}

	customSlug := c.Query("slug")

	response, err := h.createShortURL(originalURL, customSlug)
	if err != nil {
		statusCode := http.StatusBadRequest
		if err.Error() == "slug already exists" {
			statusCode = http.StatusConflict
		} else if err.Error() == "failed to save URL" || err.Error() == "failed to generate slug" {
			statusCode = http.StatusInternalServerError
		}
		c.JSON(statusCode, types.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, types.SuccessResponse(response))
}

// GET /:slug
func (h *URLHandler) RedirectURL(c *gin.Context) {
	var req models.RedirectRequest
	
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse("Invalid slug"))
		return
	}

	// Get original URL
	originalURL, err := h.redisClient.GetURL(req.Slug)
	if err != nil {
		if err == redis.Nil {
			c.JSON(http.StatusNotFound, types.ErrorResponse("URL not found"))
			return
		}
		c.JSON(http.StatusInternalServerError, types.ErrorResponse("Failed to retrieve URL"))
		return
	}

	// Increment click count
	err = h.redisClient.IncrementClickCount(req.Slug)
	if err != nil {
		// Log error but continue with redirect
		log.Printf("Failed to increment click count for slug %s: %v", req.Slug, err)
	}

	// Redirect to original URL
	c.Redirect(http.StatusFound, originalURL)
}
