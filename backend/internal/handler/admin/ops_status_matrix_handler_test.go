package admin

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/B022MC/b022hub/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

type statusMatrixRepoStub struct {
	service.OpsRepository
	filter *service.OpsStatusMatrixFilter
	resp   *service.OpsStatusMatrixResponse
	err    error
}

func (s *statusMatrixRepoStub) GetStatusMatrix(_ context.Context, filter *service.OpsStatusMatrixFilter) (*service.OpsStatusMatrixResponse, error) {
	s.filter = filter
	if s.err != nil {
		return nil, s.err
	}
	if s.resp != nil {
		return s.resp, nil
	}
	return &service.OpsStatusMatrixResponse{}, nil
}

func newOpsStatusMatrixTestRouter(handler *OpsHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/status-matrix", handler.GetDashboardStatusMatrix)
	return r
}

func TestParseOpsStatusMatrixFilter(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/?time_range=24h&platform=openai&group_id=7&q=gpt&sort=last_checked_desc", nil)

	filter, err := parseOpsStatusMatrixFilter(c)
	require.NoError(t, err)
	require.Equal(t, "24h", filter.TimeRange)
	require.Equal(t, 3600, filter.BucketSeconds)
	require.Equal(t, "openai", filter.Platform)
	require.NotNil(t, filter.GroupID)
	require.Equal(t, int64(7), *filter.GroupID)
	require.Equal(t, "gpt", filter.Query)
	require.Equal(t, service.OpsStatusMatrixSortLastCheckedDesc, filter.Sort)
}

func TestParseOpsStatusMatrixFilterInvalid(t *testing.T) {
	tests := []string{
		"/?time_range=1h",
		"/?sort=bad",
		"/?group_id=0",
		"/?group_id=abc",
	}

	gin.SetMode(gin.TestMode)
	for _, rawURL := range tests {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, rawURL, nil)

		_, err := parseOpsStatusMatrixFilter(c)
		require.Error(t, err, "url=%s", rawURL)
	}
}

func TestOpsStatusMatrixHandlerSuccess(t *testing.T) {
	now := time.Now().UTC()
	repo := &statusMatrixRepoStub{
		resp: &service.OpsStatusMatrixResponse{
			StartTime:     now.Add(-90 * time.Minute),
			EndTime:       now,
			TimeRange:     "90m",
			BucketSeconds: 300,
			Rows: []*service.OpsStatusMatrixRow{
				{
					Platform:  "openai",
					GroupName: "主通道",
					Model:     "gpt-4.1",
				},
			},
		},
	}
	svc := service.NewOpsService(repo, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
	h := NewOpsHandler(svc)
	r := newOpsStatusMatrixTestRouter(h)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/status-matrix?time_range=24h&platform=openai&group_id=7&q=gpt&sort=availability_desc", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
	require.NotNil(t, repo.filter)
	require.Equal(t, "24h", repo.filter.TimeRange)
	require.Equal(t, "openai", repo.filter.Platform)
	require.NotNil(t, repo.filter.GroupID)
	require.Equal(t, int64(7), *repo.filter.GroupID)
	require.Equal(t, "gpt", repo.filter.Query)
	require.Equal(t, service.OpsStatusMatrixSortAvailabilityDesc, repo.filter.Sort)

	var resp responseEnvelope
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	require.Equal(t, 0, resp.Code)
}

func TestOpsStatusMatrixHandlerInvalidTimeRange(t *testing.T) {
	svc := service.NewOpsService(nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
	h := NewOpsHandler(svc)
	r := newOpsStatusMatrixTestRouter(h)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/status-matrix?time_range=1h", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusBadRequest, w.Code)
}
