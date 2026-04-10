package repository

import (
	"context"
	"database/sql"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/B022MC/b022hub/internal/pkg/usagestats"
	"github.com/B022MC/b022hub/internal/service"
	"github.com/lib/pq"
)

type opsStatusMatrixRowAgg struct {
	Platform           string
	GroupID            *int64
	GroupName          string
	Model              string
	SuccessCount       int64
	ErrorCount         int64
	ExcludedErrorCount int64
	InputTokens        int64
	CacheReadTokens    int64
	LastSuccessAt      *time.Time
	LastLatencyMs      *int64
	LastRealErrorAt    *time.Time
}

type opsStatusMatrixBucketAgg struct {
	Platform           string
	GroupID            *int64
	GroupName          string
	Model              string
	BucketStart        time.Time
	SuccessCount       int64
	ErrorCount         int64
	ExcludedErrorCount int64
}

type opsStatusMatrixRowKey struct {
	Platform string
	GroupID  int64
	HasGroup bool
	Model    string
}

func (r *opsRepository) GetStatusMatrix(ctx context.Context, filter *service.OpsStatusMatrixFilter) (*service.OpsStatusMatrixResponse, error) {
	if r == nil || r.db == nil {
		return nil, fmt.Errorf("nil ops repository")
	}
	if filter == nil {
		return nil, fmt.Errorf("nil filter")
	}
	if filter.StartTime.IsZero() || filter.EndTime.IsZero() {
		return nil, fmt.Errorf("start_time/end_time required")
	}
	if filter.BucketSeconds <= 0 {
		return nil, fmt.Errorf("bucket_seconds required")
	}

	start := filter.StartTime.UTC()
	end := filter.EndTime.UTC()

	rows, err := r.queryStatusMatrixRows(ctx, filter, start, end)
	if err != nil {
		return nil, err
	}
	bucketAggs, err := r.queryStatusMatrixBuckets(ctx, filter, start, end)
	if err != nil {
		return nil, err
	}

	rowMap := make(map[opsStatusMatrixRowKey]*service.OpsStatusMatrixRow, len(rows))
	resultRows := make([]*service.OpsStatusMatrixRow, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		item := &service.OpsStatusMatrixRow{
			Platform:           normalizeStatusMatrixPlatform(row.Platform),
			GroupID:            row.GroupID,
			GroupName:          strings.TrimSpace(row.GroupName),
			Model:              normalizeStatusMatrixModel(row.Model),
			SuccessCount:       row.SuccessCount,
			ErrorCount:         row.ErrorCount,
			ExcludedErrorCount: row.ExcludedErrorCount,
			LastSuccessAt:      cloneTimePtr(row.LastSuccessAt),
			LastLatencyMs:      cloneInt64Ptr(row.LastLatencyMs),
			Buckets:            buildStatusMatrixBuckets(start, end, filter.BucketSeconds),
		}
		if denom := row.SuccessCount + row.ErrorCount; denom > 0 {
			value := float64(row.SuccessCount) / float64(denom)
			item.Availability = &value
		}
		if row.InputTokens+row.CacheReadTokens > 0 {
			value := usagestats.ComputeCacheHitRate(row.InputTokens, row.CacheReadTokens)
			item.CacheHitRate = &value
		}
		item.LastCheckedAt = maxTimePtr(row.LastSuccessAt, row.LastRealErrorAt)

		key := makeStatusMatrixRowKey(item.Platform, item.GroupID, item.Model)
		rowMap[key] = item
		resultRows = append(resultRows, item)
	}

	for _, agg := range bucketAggs {
		if agg == nil {
			continue
		}
		key := makeStatusMatrixRowKey(
			normalizeStatusMatrixPlatform(agg.Platform),
			agg.GroupID,
			normalizeStatusMatrixModel(agg.Model),
		)
		row, ok := rowMap[key]
		if !ok || row == nil {
			row = &service.OpsStatusMatrixRow{
				Platform:           normalizeStatusMatrixPlatform(agg.Platform),
				GroupID:            agg.GroupID,
				GroupName:          strings.TrimSpace(agg.GroupName),
				Model:              normalizeStatusMatrixModel(agg.Model),
				Buckets:            buildStatusMatrixBuckets(start, end, filter.BucketSeconds),
				ExcludedErrorCount: agg.ExcludedErrorCount,
			}
			rowMap[key] = row
			resultRows = append(resultRows, row)
		}

		idx := bucketIndexForWindow(start, end, filter.BucketSeconds, agg.BucketStart)
		if idx < 0 || idx >= len(row.Buckets) {
			continue
		}
		bucket := row.Buckets[idx]
		bucket.SuccessCount = agg.SuccessCount
		bucket.ErrorCount = agg.ErrorCount
		bucket.ExcludedErrorCount = agg.ExcludedErrorCount
		bucket.Status = classifyStatusMatrixBucket(agg.SuccessCount, agg.ErrorCount)
	}

	query := strings.ToLower(strings.TrimSpace(filter.Query))
	if query != "" {
		filtered := make([]*service.OpsStatusMatrixRow, 0, len(resultRows))
		for _, row := range resultRows {
			if row == nil {
				continue
			}
			groupName := strings.ToLower(strings.TrimSpace(row.GroupName))
			model := strings.ToLower(strings.TrimSpace(row.Model))
			if strings.Contains(groupName, query) || strings.Contains(model, query) {
				filtered = append(filtered, row)
			}
		}
		resultRows = filtered
	}

	sortStatusMatrixRows(resultRows, filter.Sort)

	return &service.OpsStatusMatrixResponse{
		StartTime:     start,
		EndTime:       end,
		BucketSeconds: filter.BucketSeconds,
		TimeRange:     filter.TimeRange,
		Rows:          resultRows,
	}, nil
}

func (r *opsRepository) queryStatusMatrixRows(ctx context.Context, filter *service.OpsStatusMatrixFilter, start, end time.Time) ([]*opsStatusMatrixRowAgg, error) {
	usageWhere, usageArgs, next := buildStatusMatrixUsageWhere(filter, start, end, 1)
	errorWhere, errorArgs, _ := buildStatusMatrixErrorWhere(filter, start, end, next)

	q := `
WITH usage_base AS (
  SELECT
    COALESCE(NULLIF(g.platform, ''), a.platform, 'unknown') AS platform,
    ul.group_id AS group_id,
    COALESCE(g.name, '') AS group_name,
    COALESCE(NULLIF(ul.model, ''), NULLIF(ul.upstream_model, ''), 'unknown') AS model,
    ul.input_tokens AS input_tokens,
    ul.cache_read_tokens AS cache_read_tokens,
    ul.duration_ms AS duration_ms,
    ul.created_at AS created_at,
    ul.id AS id
  FROM usage_logs ul
  LEFT JOIN groups g ON g.id = ul.group_id
  LEFT JOIN accounts a ON a.id = ul.account_id
  ` + usageWhere + `
),
success_agg AS (
  SELECT
    platform,
    group_id,
    group_name,
    model,
    COUNT(*) AS success_count,
    COALESCE(SUM(input_tokens), 0) AS input_tokens,
    COALESCE(SUM(cache_read_tokens), 0) AS cache_read_tokens,
    MAX(created_at) AS last_success_at
  FROM usage_base
  GROUP BY 1, 2, 3, 4
),
latest_success AS (
  SELECT DISTINCT ON (platform, group_id, model)
    platform,
    group_id,
    model,
    duration_ms
  FROM usage_base ub
  ORDER BY ub.platform, ub.group_id, ub.model, ub.created_at DESC, ub.id DESC
),
error_base AS (
  SELECT
    COALESCE(NULLIF(e.platform, ''), g.platform, 'unknown') AS platform,
    e.group_id AS group_id,
    COALESCE(g.name, '') AS group_name,
    COALESCE(NULLIF(e.model, ''), 'unknown') AS model,
    e.created_at AS created_at,
    (e.is_business_limited OR COALESCE(e.upstream_status_code, e.status_code, 0) IN (429, 529)) AS excluded
  FROM ops_error_logs e
  LEFT JOIN groups g ON g.id = e.group_id
  ` + errorWhere + `
    AND COALESCE(e.status_code, 0) >= 400
),
error_agg AS (
  SELECT
    platform,
    group_id,
    group_name,
    model,
    COUNT(*) FILTER (WHERE NOT excluded) AS error_count,
    COUNT(*) FILTER (WHERE excluded) AS excluded_error_count,
    MAX(created_at) FILTER (WHERE NOT excluded) AS last_real_error_at
  FROM error_base
  GROUP BY 1, 2, 3, 4
)
SELECT
  COALESCE(s.platform, e.platform) AS platform,
  COALESCE(s.group_id, e.group_id) AS group_id,
  COALESCE(NULLIF(s.group_name, ''), NULLIF(e.group_name, ''), '') AS group_name,
    COALESCE(s.model, e.model) AS model,
    COALESCE(s.success_count, 0) AS success_count,
    COALESCE(e.error_count, 0) AS error_count,
    COALESCE(e.excluded_error_count, 0) AS excluded_error_count,
    COALESCE(s.input_tokens, 0) AS input_tokens,
    COALESCE(s.cache_read_tokens, 0) AS cache_read_tokens,
    s.last_success_at AS last_success_at,
    ls.duration_ms AS last_latency_ms,
    e.last_real_error_at AS last_real_error_at
FROM success_agg s
FULL OUTER JOIN error_agg e
  ON s.platform IS NOT DISTINCT FROM e.platform
 AND s.group_id IS NOT DISTINCT FROM e.group_id
 AND s.model = e.model
LEFT JOIN latest_success ls
  ON ls.platform IS NOT DISTINCT FROM COALESCE(s.platform, e.platform)
 AND ls.group_id IS NOT DISTINCT FROM COALESCE(s.group_id, e.group_id)
 AND ls.model = COALESCE(s.model, e.model)`

	args := append(usageArgs, errorArgs...)

	rows, err := r.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	items := make([]*opsStatusMatrixRowAgg, 0, 128)
	for rows.Next() {
		item := &opsStatusMatrixRowAgg{}
		var platform sql.NullString
		var groupID sql.NullInt64
		var groupName sql.NullString
		var model sql.NullString
		var lastSuccessAt sql.NullTime
		var lastLatencyMs sql.NullInt64
		var lastRealErrorAt sql.NullTime

		if err := rows.Scan(
			&platform,
			&groupID,
			&groupName,
			&model,
			&item.SuccessCount,
			&item.ErrorCount,
			&item.ExcludedErrorCount,
			&item.InputTokens,
			&item.CacheReadTokens,
			&lastSuccessAt,
			&lastLatencyMs,
			&lastRealErrorAt,
		); err != nil {
			return nil, err
		}

		item.Platform = normalizeStatusMatrixPlatform(platform.String)
		item.GroupID = nullableInt64Ptr(groupID)
		item.GroupName = strings.TrimSpace(groupName.String)
		item.Model = normalizeStatusMatrixModel(model.String)
		item.LastSuccessAt = nullableTimePtr(lastSuccessAt)
		item.LastLatencyMs = nullableInt64Ptr(lastLatencyMs)
		item.LastRealErrorAt = nullableTimePtr(lastRealErrorAt)
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *opsRepository) queryStatusMatrixBuckets(ctx context.Context, filter *service.OpsStatusMatrixFilter, start, end time.Time) ([]*opsStatusMatrixBucketAgg, error) {
	usageWhere, usageArgs, next := buildStatusMatrixUsageWhere(filter, start, end, 1)
	errorWhere, errorArgs, _ := buildStatusMatrixErrorWhere(filter, start, end, next)
	usageBucketExpr := opsBucketExprForUsage(filter.BucketSeconds)
	errorBucketExpr := opsBucketExprForTimestampColumn(filter.BucketSeconds, "e.created_at")

	q := `
WITH usage_base AS (
  SELECT
    COALESCE(NULLIF(g.platform, ''), a.platform, 'unknown') AS platform,
    ul.group_id AS group_id,
    COALESCE(g.name, '') AS group_name,
    COALESCE(NULLIF(ul.model, ''), NULLIF(ul.upstream_model, ''), 'unknown') AS model,
    ` + usageBucketExpr + ` AS bucket_start
	FROM usage_logs ul
	LEFT JOIN groups g ON g.id = ul.group_id
	LEFT JOIN accounts a ON a.id = ul.account_id
  ` + usageWhere + `
),
usage_buckets AS (
  SELECT
    platform,
    group_id,
    group_name,
    model,
    bucket_start,
    COUNT(*) AS success_count
  FROM usage_base
  GROUP BY 1, 2, 3, 4, 5
),
error_base AS (
  SELECT
    COALESCE(NULLIF(e.platform, ''), g.platform, 'unknown') AS platform,
    e.group_id AS group_id,
    COALESCE(g.name, '') AS group_name,
    COALESCE(NULLIF(e.model, ''), 'unknown') AS model,
    ` + errorBucketExpr + ` AS bucket_start,
    (e.is_business_limited OR COALESCE(e.upstream_status_code, e.status_code, 0) IN (429, 529)) AS excluded
  FROM ops_error_logs e
  LEFT JOIN groups g ON g.id = e.group_id
  ` + errorWhere + `
    AND COALESCE(e.status_code, 0) >= 400
),
error_buckets AS (
  SELECT
    platform,
    group_id,
    group_name,
    model,
    bucket_start,
    COUNT(*) FILTER (WHERE NOT excluded) AS error_count,
    COUNT(*) FILTER (WHERE excluded) AS excluded_error_count
  FROM error_base
  GROUP BY 1, 2, 3, 4, 5
)
SELECT
  COALESCE(u.platform, e.platform) AS platform,
  COALESCE(u.group_id, e.group_id) AS group_id,
  COALESCE(NULLIF(u.group_name, ''), NULLIF(e.group_name, ''), '') AS group_name,
  COALESCE(u.model, e.model) AS model,
  COALESCE(u.bucket_start, e.bucket_start) AS bucket_start,
  COALESCE(u.success_count, 0) AS success_count,
  COALESCE(e.error_count, 0) AS error_count,
  COALESCE(e.excluded_error_count, 0) AS excluded_error_count
FROM usage_buckets u
FULL OUTER JOIN error_buckets e
  ON u.platform IS NOT DISTINCT FROM e.platform
 AND u.group_id IS NOT DISTINCT FROM e.group_id
 AND u.model = e.model
 AND u.bucket_start = e.bucket_start`

	args := append(usageArgs, errorArgs...)

	rows, err := r.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	items := make([]*opsStatusMatrixBucketAgg, 0, 512)
	for rows.Next() {
		item := &opsStatusMatrixBucketAgg{}
		var platform sql.NullString
		var groupID sql.NullInt64
		var groupName sql.NullString
		var model sql.NullString

		if err := rows.Scan(
			&platform,
			&groupID,
			&groupName,
			&model,
			&item.BucketStart,
			&item.SuccessCount,
			&item.ErrorCount,
			&item.ExcludedErrorCount,
		); err != nil {
			return nil, err
		}
		item.Platform = normalizeStatusMatrixPlatform(platform.String)
		item.GroupID = nullableInt64Ptr(groupID)
		item.GroupName = strings.TrimSpace(groupName.String)
		item.Model = normalizeStatusMatrixModel(model.String)
		item.BucketStart = item.BucketStart.UTC()
		items = append(items, item)
	}
	return items, rows.Err()
}

func buildStatusMatrixUsageWhere(filter *service.OpsStatusMatrixFilter, start, end time.Time, startIndex int) (where string, args []any, nextIndex int) {
	platform := ""
	groupID := (*int64)(nil)
	if filter != nil {
		platform = strings.TrimSpace(strings.ToLower(filter.Platform))
		groupID = filter.GroupID
	}

	idx := startIndex
	clauses := make([]string, 0, 5)
	args = make([]any, 0, 5)

	args = append(args, start)
	clauses = append(clauses, fmt.Sprintf("ul.created_at >= $%d", idx))
	idx++
	args = append(args, end)
	clauses = append(clauses, fmt.Sprintf("ul.created_at < $%d", idx))
	idx++

	if groupID != nil && *groupID > 0 {
		args = append(args, *groupID)
		clauses = append(clauses, fmt.Sprintf("ul.group_id = $%d", idx))
		idx++
	}
	if filter != nil && filter.EnforceGroupScope {
		if len(filter.ScopedGroupIDs) == 0 {
			clauses = append(clauses, "1 = 0")
		} else {
			args = append(args, pq.Array(filter.ScopedGroupIDs))
			clauses = append(clauses, fmt.Sprintf("ul.group_id = ANY($%d::bigint[])", idx))
			idx++
		}
	}
	if platform != "" {
		args = append(args, platform)
		clauses = append(clauses, fmt.Sprintf("COALESCE(NULLIF(g.platform, ''), a.platform) = $%d", idx))
		idx++
	}

	where = "WHERE " + strings.Join(clauses, " AND ")
	return where, args, idx
}

func buildStatusMatrixErrorWhere(filter *service.OpsStatusMatrixFilter, start, end time.Time, startIndex int) (where string, args []any, nextIndex int) {
	platform := ""
	groupID := (*int64)(nil)
	if filter != nil {
		platform = strings.TrimSpace(strings.ToLower(filter.Platform))
		groupID = filter.GroupID
	}

	idx := startIndex
	clauses := make([]string, 0, 6)
	args = make([]any, 0, 6)

	args = append(args, start)
	clauses = append(clauses, fmt.Sprintf("e.created_at >= $%d", idx))
	idx++
	args = append(args, end)
	clauses = append(clauses, fmt.Sprintf("e.created_at < $%d", idx))
	idx++

	clauses = append(clauses, "e.is_count_tokens = FALSE")

	if groupID != nil && *groupID > 0 {
		args = append(args, *groupID)
		clauses = append(clauses, fmt.Sprintf("e.group_id = $%d", idx))
		idx++
	}
	if filter != nil && filter.EnforceGroupScope {
		if len(filter.ScopedGroupIDs) == 0 {
			clauses = append(clauses, "1 = 0")
		} else {
			args = append(args, pq.Array(filter.ScopedGroupIDs))
			clauses = append(clauses, fmt.Sprintf("e.group_id = ANY($%d::bigint[])", idx))
			idx++
		}
	}
	if platform != "" {
		args = append(args, platform)
		clauses = append(clauses, fmt.Sprintf("COALESCE(NULLIF(e.platform, ''), g.platform) = $%d", idx))
		idx++
	}

	where = "WHERE " + strings.Join(clauses, " AND ")
	return where, args, idx
}

func buildStatusMatrixBuckets(start, end time.Time, bucketSeconds int) []*service.OpsStatusMatrixBucket {
	if bucketSeconds <= 0 || !start.Before(end) {
		return []*service.OpsStatusMatrixBucket{}
	}

	endMinus := end.Add(-time.Nanosecond)
	if endMinus.Before(start) {
		return []*service.OpsStatusMatrixBucket{}
	}

	first := opsFloorToBucketStart(start, bucketSeconds)
	last := opsFloorToBucketStart(endMinus, bucketSeconds)
	step := time.Duration(bucketSeconds) * time.Second
	out := make([]*service.OpsStatusMatrixBucket, 0, int(last.Sub(first)/step)+1)

	for cursor := first; !cursor.After(last); cursor = cursor.Add(step) {
		out = append(out, &service.OpsStatusMatrixBucket{
			BucketStart: cursor,
			BucketEnd:   cursor.Add(step),
			Status:      service.OpsStatusMatrixBucketStatusNoData,
		})
	}

	return out
}

func bucketIndexForWindow(start, end time.Time, bucketSeconds int, bucketStart time.Time) int {
	if bucketSeconds <= 0 || !start.Before(end) {
		return -1
	}
	first := opsFloorToBucketStart(start, bucketSeconds)
	endMinus := end.Add(-time.Nanosecond)
	last := opsFloorToBucketStart(endMinus, bucketSeconds)
	bucketStart = bucketStart.UTC()
	if bucketStart.Before(first) || bucketStart.After(last) {
		return -1
	}
	step := time.Duration(bucketSeconds) * time.Second
	return int(bucketStart.Sub(first) / step)
}

func classifyStatusMatrixBucket(successCount, errorCount int64) service.OpsStatusMatrixBucketStatus {
	switch {
	case successCount == 0 && errorCount == 0:
		return service.OpsStatusMatrixBucketStatusNoData
	case successCount > 0 && errorCount == 0:
		return service.OpsStatusMatrixBucketStatusOK
	case successCount > 0 && errorCount > 0:
		return service.OpsStatusMatrixBucketStatusWarn
	default:
		return service.OpsStatusMatrixBucketStatusDown
	}
}

func sortStatusMatrixRows(rows []*service.OpsStatusMatrixRow, sortMode service.OpsStatusMatrixSort) {
	if !sortMode.IsValid() {
		sortMode = service.OpsStatusMatrixSortAvailabilityAsc
	}

	sort.SliceStable(rows, func(i, j int) bool {
		left := rows[i]
		right := rows[j]
		if left == nil || right == nil {
			return left != nil
		}

		switch sortMode {
		case service.OpsStatusMatrixSortAvailabilityDesc:
			if cmp := compareAvailability(left.Availability, right.Availability, true); cmp != 0 {
				return cmp < 0
			}
		case service.OpsStatusMatrixSortLastCheckedDesc:
			if cmp := compareTimeDesc(left.LastCheckedAt, right.LastCheckedAt); cmp != 0 {
				return cmp < 0
			}
			if cmp := compareAvailability(left.Availability, right.Availability, false); cmp != 0 {
				return cmp < 0
			}
		default:
			if cmp := compareAvailability(left.Availability, right.Availability, false); cmp != 0 {
				return cmp < 0
			}
		}

		if cmp := compareTimeDesc(left.LastCheckedAt, right.LastCheckedAt); cmp != 0 {
			return cmp < 0
		}
		if left.Platform != right.Platform {
			return left.Platform < right.Platform
		}
		if strings.TrimSpace(left.GroupName) != strings.TrimSpace(right.GroupName) {
			return strings.TrimSpace(left.GroupName) < strings.TrimSpace(right.GroupName)
		}
		return left.Model < right.Model
	})
}

func compareAvailability(left, right *float64, desc bool) int {
	leftValid := left != nil
	rightValid := right != nil
	if leftValid != rightValid {
		if leftValid {
			return -1
		}
		return 1
	}
	if !leftValid {
		return 0
	}
	if *left == *right {
		return 0
	}
	if desc {
		if *left > *right {
			return -1
		}
		return 1
	}
	if *left < *right {
		return -1
	}
	return 1
}

func compareTimeDesc(left, right *time.Time) int {
	leftValid := left != nil
	rightValid := right != nil
	if leftValid != rightValid {
		if leftValid {
			return -1
		}
		return 1
	}
	if !leftValid {
		return 0
	}
	if left.Equal(*right) {
		return 0
	}
	if left.After(*right) {
		return -1
	}
	return 1
}

func makeStatusMatrixRowKey(platform string, groupID *int64, model string) opsStatusMatrixRowKey {
	if groupID == nil {
		return opsStatusMatrixRowKey{
			Platform: platform,
			HasGroup: false,
			Model:    model,
		}
	}
	return opsStatusMatrixRowKey{
		Platform: platform,
		GroupID:  *groupID,
		HasGroup: true,
		Model:    model,
	}
}

func normalizeStatusMatrixPlatform(platform string) string {
	out := strings.TrimSpace(strings.ToLower(platform))
	if out == "" {
		return "unknown"
	}
	return out
}

func normalizeStatusMatrixModel(model string) string {
	out := strings.TrimSpace(model)
	if out == "" {
		return "unknown"
	}
	return out
}

func nullableTimePtr(value sql.NullTime) *time.Time {
	if !value.Valid {
		return nil
	}
	out := value.Time.UTC()
	return &out
}

func nullableInt64Ptr(value sql.NullInt64) *int64 {
	if !value.Valid {
		return nil
	}
	out := value.Int64
	return &out
}

func cloneTimePtr(value *time.Time) *time.Time {
	if value == nil {
		return nil
	}
	out := value.UTC()
	return &out
}

func cloneInt64Ptr(value *int64) *int64 {
	if value == nil {
		return nil
	}
	out := *value
	return &out
}

func maxTimePtr(left, right *time.Time) *time.Time {
	switch {
	case left == nil:
		return cloneTimePtr(right)
	case right == nil:
		return cloneTimePtr(left)
	case left.After(*right):
		return cloneTimePtr(left)
	default:
		return cloneTimePtr(right)
	}
}
