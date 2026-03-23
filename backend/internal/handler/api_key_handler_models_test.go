package handler

import (
	"testing"

	"github.com/B022MC/b022hub/internal/service"
	"github.com/stretchr/testify/require"
)

func TestNormalizeAvailableModelIDs_UsesProvidedModels(t *testing.T) {
	got := normalizeAvailableModelIDs([]string{" gpt-5.4 ", "gpt-5.2", "gpt-5.4", ""}, service.PlatformOpenAI)

	require.Equal(t, []string{"gpt-5.2", "gpt-5.4"}, got)
}

func TestNormalizeAvailableModelIDs_FallsBackToGeminiDefaults(t *testing.T) {
	got := normalizeAvailableModelIDs(nil, service.PlatformGemini)

	require.NotEmpty(t, got)
	require.Contains(t, got, "gemini-2.5-pro")
	for _, modelID := range got {
		require.NotContains(t, modelID, "models/")
	}
}

func TestNormalizeAvailableModelIDs_FallsBackToOpenAIDefaults(t *testing.T) {
	got := normalizeAvailableModelIDs(nil, service.PlatformOpenAI)

	require.NotEmpty(t, got)
	require.Contains(t, got, "gpt-5.4")
}
