package service

import "context"

func (s *OpsService) GetStatusMatrix(ctx context.Context, filter *OpsStatusMatrixFilter) (*OpsStatusMatrixResponse, error) {
	if s == nil || s.opsRepo == nil {
		return nil, ErrOpsDisabled
	}
	if filter == nil {
		return nil, ErrOpsDisabled
	}
	return s.opsRepo.GetStatusMatrix(ctx, filter)
}
