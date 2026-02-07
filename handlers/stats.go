package handlers

import (
	"context"
	"net/http"
	"strconv"
	"gourl/db"
	"gourl/config"
	"gourl/types"

	"github.com/gin-gonic/gin"
)

type StatsHandler struct {
	redisClient *database.RedisClient
	ctx         context.Context
}

func NewStatsHandler(redisClient *database.RedisClient) *StatsHandler {
	return &StatsHandler{
		redisClient: redisClient,
		ctx:         context.Background(),
	}
}

// GET /api/stats
func (h *StatsHandler) GetAllStats(c *gin.Context) {
	stats, err := h.redisClient.GetAllStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse("Failed to retrieve stats"))
		return
	}

	// Convert to structured response
	cfg := config.LoadConfig()
	var response []types.StatsResponse
	for _, stat := range stats {
		clickCount, _ := strconv.ParseInt(stat["click_count"], 10, 64)
		createdAt, _ := strconv.ParseInt(stat["created_at"], 10, 64)
		lastAccessed, _ := strconv.ParseInt(stat["last_accessed"], 10, 64)

		item := types.StatsResponse{
			Slug:         stat["slug"],
			OriginalURL:  stat["original_url"],
			ShortURL:     cfg.BaseURL + "/" + stat["slug"],
			ClickCount:   clickCount,
			CreatedAt:    createdAt,
			LastAccessed: lastAccessed,
		}

		response = append(response, item)
	}

	// Buat meta data
	meta := &types.Meta{
		Total: len(response),
	}

	c.JSON(http.StatusOK, types.Response{
		Success: true,
		Data:    response,
		Meta:    meta,
	})
}
