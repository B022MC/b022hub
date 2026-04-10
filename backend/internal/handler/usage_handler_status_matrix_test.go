package handler

import (
	"testing"

	"github.com/B022MC/b022hub/internal/service"
	"github.com/stretchr/testify/require"
)

func TestScopeStatusMatrixFilterToGroups(t *testing.T) {
	filter := &service.OpsStatusMatrixFilter{}

	scopeStatusMatrixFilterToGroups(filter, []service.Group{
		{ID: 7},
		{ID: 9},
		{ID: 7},
	})

	require.True(t, filter.EnforceGroupScope)
	require.Equal(t, []int64{7, 9}, filter.ScopedGroupIDs)
}

func TestScopeStatusMatrixFilterToGroupsEmpty(t *testing.T) {
	filter := &service.OpsStatusMatrixFilter{}

	scopeStatusMatrixFilterToGroups(filter, nil)

	require.True(t, filter.EnforceGroupScope)
	require.Nil(t, filter.ScopedGroupIDs)
}
