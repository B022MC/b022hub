package dto

import (
	"testing"

	"github.com/B022MC/b022hub/internal/service"
	"github.com/stretchr/testify/require"
)

func TestGroupFromServiceAdmin_MapsDefaultProxyID(t *testing.T) {
	t.Parallel()

	defaultProxyID := int64(88)
	out := GroupFromServiceAdmin(&service.Group{
		ID:             1,
		Name:           "openai-default",
		Platform:       service.PlatformOpenAI,
		Status:         service.StatusActive,
		DefaultProxyID: &defaultProxyID,
	})

	require.NotNil(t, out)
	require.NotNil(t, out.DefaultProxyID)
	require.Equal(t, defaultProxyID, *out.DefaultProxyID)
}
