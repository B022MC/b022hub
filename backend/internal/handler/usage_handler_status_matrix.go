package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/B022MC/b022hub/internal/pkg/response"
	middleware2 "github.com/B022MC/b022hub/internal/server/middleware"
	"github.com/B022MC/b022hub/internal/service"
	"github.com/gin-gonic/gin"
)

const (
	userOpsStatusMatrixTimeRange90m = "90m"
	userOpsStatusMatrixTimeRange24h = "24h"
)

// GetStatusMatrix returns the real-traffic status matrix for the current user scope.
// GET /api/v1/ops/status-matrix
func (h *UsageHandler) GetStatusMatrix(c *gin.Context) {
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
	}
	if h.opsService == nil {
		response.Error(c, http.StatusServiceUnavailable, "Ops service not available")
		return
	}
	if err := h.opsService.RequireMonitoringEnabled(c.Request.Context()); err != nil {
		response.ErrorFrom(c, err)
		return
	}

	filter, err := parseUserOpsStatusMatrixFilter(c)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	role, _ := middleware2.GetUserRoleFromContext(c)
	if !strings.EqualFold(strings.TrimSpace(role), "admin") {
		if h.apiKeyService == nil {
			response.Error(c, http.StatusServiceUnavailable, "API key service not available")
			return
		}
		availableGroups, err := h.apiKeyService.GetAvailableGroups(c.Request.Context(), subject.UserID)
		if err != nil {
			response.ErrorFrom(c, err)
			return
		}
		scopeStatusMatrixFilterToGroups(filter, availableGroups)
	}

	data, err := h.opsService.GetStatusMatrix(c.Request.Context(), filter)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, data)
}

func parseUserOpsStatusMatrixFilter(c *gin.Context) (*service.OpsStatusMatrixFilter, error) {
	if c == nil {
		return nil, fmt.Errorf("invalid request")
	}

	timeRange := strings.TrimSpace(c.Query("time_range"))
	if timeRange == "" {
		timeRange = userOpsStatusMatrixTimeRange90m
	}

	var window time.Duration
	var bucketSeconds int
	switch timeRange {
	case userOpsStatusMatrixTimeRange90m:
		window = 90 * time.Minute
		bucketSeconds = 300
	case userOpsStatusMatrixTimeRange24h:
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

func scopeStatusMatrixFilterToGroups(filter *service.OpsStatusMatrixFilter, groups []service.Group) {
	if filter == nil {
		return
	}
	filter.EnforceGroupScope = true
	if len(groups) == 0 {
		filter.ScopedGroupIDs = nil
		return
	}

	seen := make(map[int64]struct{}, len(groups))
	groupIDs := make([]int64, 0, len(groups))
	for _, group := range groups {
		if group.ID <= 0 {
			continue
		}
		if _, exists := seen[group.ID]; exists {
			continue
		}
		seen[group.ID] = struct{}{}
		groupIDs = append(groupIDs, group.ID)
	}
	filter.ScopedGroupIDs = groupIDs
}
