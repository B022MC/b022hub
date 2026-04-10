package repository

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/B022MC/b022hub/internal/service"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

func TestOpsRepository_GetStatusMatrix(t *testing.T) {
	db, mock := newSQLMock(t)
	repo := &opsRepository{db: db}

	start := time.Date(2026, 4, 10, 8, 30, 0, 0, time.UTC)
	end := start.Add(90 * time.Minute)
	filter := &service.OpsStatusMatrixFilter{
		StartTime:     start,
		EndTime:       end,
		TimeRange:     "90m",
		BucketSeconds: 300,
		Sort:          service.OpsStatusMatrixSortAvailabilityAsc,
	}

	rowQueryRows := sqlmock.NewRows([]string{
		"platform",
		"group_id",
		"group_name",
		"model",
		"success_count",
		"error_count",
		"excluded_error_count",
		"input_tokens",
		"cache_read_tokens",
		"last_success_at",
		"last_latency_ms",
		"last_real_error_at",
	}).
		AddRow("openai", int64(7), "主通道", "gpt-4.1", int64(2), int64(1), int64(1), int64(80), int64(20), end.Add(-10*time.Minute), int64(1234), end.Add(-5*time.Minute)).
		AddRow("gemini", int64(8), "备用通道", "gemini-2.5-pro", int64(0), int64(1), int64(0), int64(0), int64(0), nil, nil, end.Add(-20*time.Minute))

	bucketRows := sqlmock.NewRows([]string{
		"platform",
		"group_id",
		"group_name",
		"model",
		"bucket_start",
		"success_count",
		"error_count",
		"excluded_error_count",
	}).
		AddRow("openai", int64(7), "主通道", "gpt-4.1", start, int64(1), int64(0), int64(0)).
		AddRow("openai", int64(7), "主通道", "gpt-4.1", start.Add(5*time.Minute), int64(1), int64(1), int64(1)).
		AddRow("gemini", int64(8), "备用通道", "gemini-2.5-pro", start.Add(10*time.Minute), int64(0), int64(1), int64(0))

	mock.ExpectQuery(regexp.QuoteMeta("WITH usage_base AS (")).
		WithArgs(start, end, start, end).
		WillReturnRows(rowQueryRows)
	mock.ExpectQuery(regexp.QuoteMeta("WITH usage_base AS (")).
		WithArgs(start, end, start, end).
		WillReturnRows(bucketRows)

	result, err := repo.GetStatusMatrix(context.Background(), filter)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, 300, result.BucketSeconds)
	require.Equal(t, "90m", result.TimeRange)
	require.Len(t, result.Rows, 2)

	first := result.Rows[0]
	require.Equal(t, "gemini", first.Platform)
	require.Equal(t, "gemini-2.5-pro", first.Model)
	require.NotNil(t, first.Availability)
	require.InDelta(t, 0.0, *first.Availability, 0.0001)
	require.NotNil(t, first.LastCheckedAt)
	require.Nil(t, first.LastSuccessAt)
	require.Nil(t, first.LastLatencyMs)
	require.Len(t, first.Buckets, 18)
	require.Equal(t, service.OpsStatusMatrixBucketStatusDown, first.Buckets[2].Status)

	second := result.Rows[1]
	require.Equal(t, "openai", second.Platform)
	require.Equal(t, "主通道", second.GroupName)
	require.NotNil(t, second.Availability)
	require.InDelta(t, 2.0/3.0, *second.Availability, 0.0001)
	require.NotNil(t, second.CacheHitRate)
	require.InDelta(t, 0.2, *second.CacheHitRate, 0.0001)
	require.Equal(t, int64(1), second.ExcludedErrorCount)
	require.NotNil(t, second.LastCheckedAt)
	require.NotNil(t, second.LastSuccessAt)
	require.NotNil(t, second.LastLatencyMs)
	require.Equal(t, int64(1234), *second.LastLatencyMs)
	require.Equal(t, service.OpsStatusMatrixBucketStatusOK, second.Buckets[0].Status)
	require.Equal(t, service.OpsStatusMatrixBucketStatusWarn, second.Buckets[1].Status)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestOpsBucketExprForTimestampColumnUsesQualifiedColumn(t *testing.T) {
	require.Equal(
		t,
		"to_timestamp(floor(extract(epoch from e.created_at) / 300) * 300)",
		opsBucketExprForTimestampColumn(300, "e.created_at"),
	)
	require.Equal(
		t,
		"date_trunc('hour', e.created_at)",
		opsBucketExprForTimestampColumn(3600, "e.created_at"),
	)
}

func TestOpsRepository_GetStatusMatrixFiltersAndExcludedOnlyRows(t *testing.T) {
	db, mock := newSQLMock(t)
	repo := &opsRepository{db: db}

	start := time.Date(2026, 4, 10, 0, 0, 0, 0, time.UTC)
	end := start.Add(24 * time.Hour)
	groupID := int64(7)
	filter := &service.OpsStatusMatrixFilter{
		StartTime:     start,
		EndTime:       end,
		TimeRange:     "24h",
		BucketSeconds: 3600,
		Platform:      "openai",
		GroupID:       &groupID,
		Query:         "通道",
		Sort:          service.OpsStatusMatrixSortLastCheckedDesc,
	}

	rowQueryRows := sqlmock.NewRows([]string{
		"platform",
		"group_id",
		"group_name",
		"model",
		"success_count",
		"error_count",
		"excluded_error_count",
		"input_tokens",
		"cache_read_tokens",
		"last_success_at",
		"last_latency_ms",
		"last_real_error_at",
	}).
		AddRow("openai", groupID, "通道-A", "gpt-4.1", int64(0), int64(0), int64(2), int64(0), int64(0), nil, nil, nil).
		AddRow("openai", groupID, "其他", "gpt-4.1-mini", int64(1), int64(0), int64(0), int64(10), int64(0), end.Add(-time.Hour), int64(400), nil)

	bucketRows := sqlmock.NewRows([]string{
		"platform",
		"group_id",
		"group_name",
		"model",
		"bucket_start",
		"success_count",
		"error_count",
		"excluded_error_count",
	}).
		AddRow("openai", groupID, "通道-A", "gpt-4.1", start.Add(2*time.Hour), int64(0), int64(0), int64(2)).
		AddRow("openai", groupID, "其他", "gpt-4.1-mini", start.Add(3*time.Hour), int64(1), int64(0), int64(0))

	mock.ExpectQuery(regexp.QuoteMeta("WITH usage_base AS (")).
		WithArgs(start, end, groupID, "openai", start, end, groupID, "openai").
		WillReturnRows(rowQueryRows)
	mock.ExpectQuery(regexp.QuoteMeta("WITH usage_base AS (")).
		WithArgs(start, end, groupID, "openai", start, end, groupID, "openai").
		WillReturnRows(bucketRows)

	result, err := repo.GetStatusMatrix(context.Background(), filter)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Len(t, result.Rows, 1)
	require.Equal(t, "通道-A", result.Rows[0].GroupName)
	require.Nil(t, result.Rows[0].Availability)
	require.Equal(t, int64(2), result.Rows[0].ExcludedErrorCount)
	require.Len(t, result.Rows[0].Buckets, 24)
	require.Equal(t, service.OpsStatusMatrixBucketStatusNoData, result.Rows[0].Buckets[2].Status)
	require.Equal(t, int64(2), result.Rows[0].Buckets[2].ExcludedErrorCount)

	require.NoError(t, mock.ExpectationsWereMet())
}
