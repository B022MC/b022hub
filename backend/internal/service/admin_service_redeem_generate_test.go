//go:build unit

package service

import (
	"context"
	"net/http"
	"testing"

	infraerrors "github.com/B022MC/b022hub/internal/pkg/errors"
	"github.com/stretchr/testify/require"
)

type redeemGenerateRepoStub struct {
	redeemRepoStub
	created []RedeemCode
}

func (s *redeemGenerateRepoStub) Create(_ context.Context, code *RedeemCode) error {
	s.created = append(s.created, *code)
	return nil
}

func TestGenerateRedeemCodes_RejectsCountsAboveBatchLimit(t *testing.T) {
	repo := &redeemGenerateRepoStub{}
	svc := &adminServiceImpl{redeemCodeRepo: repo}

	codes, err := svc.GenerateRedeemCodes(context.Background(), &GenerateRedeemCodesInput{
		Count: 501,
		Type:  RedeemTypeBalance,
		Value: 10,
	})

	require.Nil(t, codes)
	require.Error(t, err)
	require.Equal(t, http.StatusBadRequest, infraerrors.Code(err))
	require.Empty(t, repo.created)
}
