package admin

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/B022MC/b022hub/internal/pkg/response"
	"github.com/B022MC/b022hub/internal/service"
	"github.com/gin-gonic/gin"
)

const (
	opsStatusMatrixTimeRange90m = "90m"
	opsStatusMatrixTimeRange24h = "24h"
)

func (h *OpsHandler) GetDashboardStatusMatrix(c *gin.Context) {
	if h.opsService == nil {
		response.Error(c, http.StatusServiceUnavailable, "Ops service not available")
		return
	}
	if err := h.opsService.RequireMonitoringEnabled(c.Request.Context()); err != nil {
		response.ErrorFrom(c, err)
		return
	}

	filter, err := parseOpsStatusMatrixFilter(c)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	data, err := h.opsService.GetStatusMatrix(c.Request.Context(), filter)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, data)
}

func parseOpsStatusMatrixFilter(c *gin.Context) (*service.OpsStatusMatrixFilter, error) {
	if c == nil {
		return nil, fmt.Errorf("invalid request")
	}

	timeRange := strings.TrimSpace(c.Query("time_range"))
	if timeRange == "" {
		timeRange = opsStatusMatrixTimeRange90m
	}

	var window time.Duration
	var bucketSeconds int
	switch timeRange {
	case opsStatusMatrixTimeRange90m:
		window = 90 * time.Minute
		bucketSeconds = 300
	case opsStatusMatrixTimeRange24h:
		window = 24 * time.Hour
		bucketSeconds = 3600
	default:
		return nil, fmt.Errorf("invalid time_range")
	}

	sortMode := service.OpsStatusMatrixSort(strings.TrimSpace(c.Query("sort")))
	if sortMode == "" {
		sortMode = service.OpsStatusMatrixSortAvailabilityAsc
	}
	if !sortMode.IsValid() {
		return nil, fmt.Errorf("invalid sort")
	}

	filter := &service.OpsStatusMatrixFilter{
		StartTime:     time.Now().UTC().Add(-window),
		EndTime:       time.Now().UTC(),
		TimeRange:     timeRange,
		BucketSeconds: bucketSeconds,
		Platform:      strings.TrimSpace(c.Query("platform")),
		Query:         strings.TrimSpace(c.Query("q")),
		Sort:          sortMode,
	}

	if v := strings.TrimSpace(c.Query("group_id")); v != "" {
		id, err := strconv.ParseInt(v, 10, 64)
		if err != nil || id <= 0 {
			return nil, fmt.Errorf("invalid group_id")
		}
		filter.GroupID = &id
	}

	return filter, nil
}
