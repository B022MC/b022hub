package service

import (
	"context"
	"fmt"
)

func resolveGroupDefaultProxyID(ctx context.Context, groupRepo GroupRepository, groupIDs []int64) (*int64, error) {
	if len(groupIDs) == 0 {
		return nil, nil
	}
	if groupRepo == nil {
		return nil, fmt.Errorf("group repository not configured")
	}

	for _, groupID := range groupIDs {
		if groupID <= 0 {
			return nil, fmt.Errorf("get group: %w", ErrGroupNotFound)
		}

		group, err := groupRepo.GetByIDLite(ctx, groupID)
		if err != nil {
			return nil, fmt.Errorf("get group: %w", err)
		}
		if group.DefaultProxyID == nil || *group.DefaultProxyID <= 0 {
			continue
		}

		proxyID := *group.DefaultProxyID
		return &proxyID, nil
	}

	return nil, nil
}
