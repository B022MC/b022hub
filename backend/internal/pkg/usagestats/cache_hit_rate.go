package usagestats

// ComputeCacheHitRate returns cache_read / (input + cache_read).
// input_tokens in this project are normalized to exclude cache_read_tokens,
// so the denominator represents the full prompt-side token volume.
func ComputeCacheHitRate(inputTokens, cacheReadTokens int64) float64 {
	denominator := inputTokens + cacheReadTokens
	if denominator <= 0 || cacheReadTokens <= 0 {
		return 0
	}
	return float64(cacheReadTokens) / float64(denominator)
}

func (s *DashboardStats) ApplyDerivedMetrics() {
	if s == nil {
		return
	}
	s.TotalTokens = s.TotalInputTokens + s.TotalOutputTokens + s.TotalCacheCreationTokens + s.TotalCacheReadTokens
	s.TodayTokens = s.TodayInputTokens + s.TodayOutputTokens + s.TodayCacheCreationTokens + s.TodayCacheReadTokens
	s.TotalCacheHitRate = ComputeCacheHitRate(s.TotalInputTokens, s.TotalCacheReadTokens)
	s.TodayCacheHitRate = ComputeCacheHitRate(s.TodayInputTokens, s.TodayCacheReadTokens)
}

func (s *UserDashboardStats) ApplyDerivedMetrics() {
	if s == nil {
		return
	}
	s.TotalTokens = s.TotalInputTokens + s.TotalOutputTokens + s.TotalCacheCreationTokens + s.TotalCacheReadTokens
	s.TodayTokens = s.TodayInputTokens + s.TodayOutputTokens + s.TodayCacheCreationTokens + s.TodayCacheReadTokens
	s.TotalCacheHitRate = ComputeCacheHitRate(s.TotalInputTokens, s.TotalCacheReadTokens)
	s.TodayCacheHitRate = ComputeCacheHitRate(s.TodayInputTokens, s.TodayCacheReadTokens)
}

func (s *ModelStat) ApplyDerivedMetrics() {
	if s == nil {
		return
	}
	s.CacheHitRate = ComputeCacheHitRate(s.InputTokens, s.CacheReadTokens)
}
