package usagestats

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestComputeCacheHitRate(t *testing.T) {
	t.Run("returns zero when denominator is empty", func(t *testing.T) {
		require.Zero(t, ComputeCacheHitRate(0, 0))
	})

	t.Run("returns zero when there are no cache hits", func(t *testing.T) {
		require.Zero(t, ComputeCacheHitRate(128, 0))
	})

	t.Run("calculates normalized cache hit rate", func(t *testing.T) {
		require.InDelta(t, 0.2, ComputeCacheHitRate(80, 20), 0.000001)
	})
}
