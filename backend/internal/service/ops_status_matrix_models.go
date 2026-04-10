package service

import "time"

type OpsStatusMatrixSort string

const (
	OpsStatusMatrixSortAvailabilityAsc  OpsStatusMatrixSort = "availability_asc"
	OpsStatusMatrixSortAvailabilityDesc OpsStatusMatrixSort = "availability_desc"
	OpsStatusMatrixSortLastCheckedDesc  OpsStatusMatrixSort = "last_checked_desc"
)

func (s OpsStatusMatrixSort) IsValid() bool {
	switch s {
	case OpsStatusMatrixSortAvailabilityAsc,
		OpsStatusMatrixSortAvailabilityDesc,
		OpsStatusMatrixSortLastCheckedDesc:
		return true
	default:
		return false
	}
}

type OpsStatusMatrixBucketStatus string

const (
	OpsStatusMatrixBucketStatusOK     OpsStatusMatrixBucketStatus = "ok"
	OpsStatusMatrixBucketStatusWarn   OpsStatusMatrixBucketStatus = "warn"
	OpsStatusMatrixBucketStatusDown   OpsStatusMatrixBucketStatus = "down"
	OpsStatusMatrixBucketStatusNoData OpsStatusMatrixBucketStatus = "nodata"
)

type OpsStatusMatrixFilter struct {
	StartTime time.Time
	EndTime   time.Time

	TimeRange     string
	BucketSeconds int

	Platform string
	GroupID  *int64
	Query    string
	Sort     OpsStatusMatrixSort
}

type OpsStatusMatrixBucket struct {
	BucketStart        time.Time                   `json:"bucket_start"`
	BucketEnd          time.Time                   `json:"bucket_end"`
	SuccessCount       int64                       `json:"success_count"`
	ErrorCount         int64                       `json:"error_count"`
	ExcludedErrorCount int64                       `json:"excluded_error_count"`
	Status             OpsStatusMatrixBucketStatus `json:"status"`
}

type OpsStatusMatrixRow struct {
	Platform string `json:"platform"`

	GroupID   *int64 `json:"group_id"`
	GroupName string `json:"group_name"`
	Model     string `json:"model"`

	SuccessCount       int64    `json:"success_count"`
	ErrorCount         int64    `json:"error_count"`
	ExcludedErrorCount int64    `json:"excluded_error_count"`
	Availability       *float64 `json:"availability"`

	LastCheckedAt *time.Time `json:"last_checked_at"`
	LastSuccessAt *time.Time `json:"last_success_at"`
	LastLatencyMs *int64     `json:"last_latency_ms"`

	Buckets []*OpsStatusMatrixBucket `json:"buckets"`
}

type OpsStatusMatrixResponse struct {
	StartTime     time.Time `json:"start_time"`
	EndTime       time.Time `json:"end_time"`
	BucketSeconds int       `json:"bucket_seconds"`
	TimeRange     string    `json:"time_range"`

	Rows []*OpsStatusMatrixRow `json:"rows"`
}
