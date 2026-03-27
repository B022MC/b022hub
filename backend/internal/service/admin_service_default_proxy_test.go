//go:build unit

package service

import (
	"context"
	"testing"

	"github.com/B022MC/b022hub/internal/pkg/pagination"
	"github.com/stretchr/testify/require"
)

type bindGroupsCall struct {
	accountID int64
	groupIDs  []int64
}

type accountRepoStubForDefaultProxy struct {
	accountRepoStub

	nextID          int64
	created         *Account
	updated         []*Account
	accountsByID    map[int64]*Account
	bindGroupsCalls []bindGroupsCall
	bulkUpdateIDs   []int64
}

func newAccountRepoStubForDefaultProxy() *accountRepoStubForDefaultProxy {
	return &accountRepoStubForDefaultProxy{
		nextID:       1,
		accountsByID: make(map[int64]*Account),
	}
}

func (s *accountRepoStubForDefaultProxy) Create(_ context.Context, account *Account) error {
	cloned := cloneAccountForDefaultProxyTest(account)
	if cloned.ID == 0 {
		cloned.ID = s.nextID
		s.nextID++
	}
	account.ID = cloned.ID
	s.created = cloned
	s.accountsByID[cloned.ID] = cloneAccountForDefaultProxyTest(cloned)
	return nil
}

func (s *accountRepoStubForDefaultProxy) GetByID(_ context.Context, id int64) (*Account, error) {
	account, ok := s.accountsByID[id]
	if !ok {
		return nil, ErrAccountNotFound
	}
	return cloneAccountForDefaultProxyTest(account), nil
}

func (s *accountRepoStubForDefaultProxy) GetByIDs(_ context.Context, ids []int64) ([]*Account, error) {
	out := make([]*Account, 0, len(ids))
	for _, id := range ids {
		if account, ok := s.accountsByID[id]; ok {
			out = append(out, cloneAccountForDefaultProxyTest(account))
		}
	}
	return out, nil
}

func (s *accountRepoStubForDefaultProxy) Update(_ context.Context, account *Account) error {
	cloned := cloneAccountForDefaultProxyTest(account)
	s.updated = append(s.updated, cloned)
	s.accountsByID[cloned.ID] = cloneAccountForDefaultProxyTest(cloned)
	return nil
}

func (s *accountRepoStubForDefaultProxy) BindGroups(_ context.Context, accountID int64, groupIDs []int64) error {
	s.bindGroupsCalls = append(s.bindGroupsCalls, bindGroupsCall{
		accountID: accountID,
		groupIDs:  append([]int64(nil), groupIDs...),
	})
	account, ok := s.accountsByID[accountID]
	if ok {
		account.GroupIDs = append([]int64(nil), groupIDs...)
	}
	return nil
}

func (s *accountRepoStubForDefaultProxy) BulkUpdate(_ context.Context, ids []int64, _ AccountBulkUpdate) (int64, error) {
	s.bulkUpdateIDs = append([]int64(nil), ids...)
	return int64(len(ids)), nil
}

type proxyRepoStubForDefaultProxy struct {
	proxy *Proxy
	err   error
}

func (s *proxyRepoStubForDefaultProxy) Create(context.Context, *Proxy) error {
	panic("unexpected Create call")
}

func (s *proxyRepoStubForDefaultProxy) GetByID(context.Context, int64) (*Proxy, error) {
	if s.err != nil {
		return nil, s.err
	}
	if s.proxy == nil {
		return nil, ErrProxyNotFound
	}
	return s.proxy, nil
}

func (s *proxyRepoStubForDefaultProxy) ListByIDs(context.Context, []int64) ([]Proxy, error) {
	panic("unexpected ListByIDs call")
}

func (s *proxyRepoStubForDefaultProxy) Update(context.Context, *Proxy) error {
	panic("unexpected Update call")
}

func (s *proxyRepoStubForDefaultProxy) Delete(context.Context, int64) error {
	panic("unexpected Delete call")
}

func (s *proxyRepoStubForDefaultProxy) List(context.Context, pagination.PaginationParams) ([]Proxy, *pagination.PaginationResult, error) {
	panic("unexpected List call")
}

func (s *proxyRepoStubForDefaultProxy) ListWithFilters(context.Context, pagination.PaginationParams, string, string, string) ([]Proxy, *pagination.PaginationResult, error) {
	panic("unexpected ListWithFilters call")
}

func (s *proxyRepoStubForDefaultProxy) ListWithFiltersAndAccountCount(context.Context, pagination.PaginationParams, string, string, string) ([]ProxyWithAccountCount, *pagination.PaginationResult, error) {
	panic("unexpected ListWithFiltersAndAccountCount call")
}

func (s *proxyRepoStubForDefaultProxy) ListActive(context.Context) ([]Proxy, error) {
	panic("unexpected ListActive call")
}

func (s *proxyRepoStubForDefaultProxy) ListActiveWithAccountCount(context.Context) ([]ProxyWithAccountCount, error) {
	panic("unexpected ListActiveWithAccountCount call")
}

func (s *proxyRepoStubForDefaultProxy) ExistsByHostPortAuth(context.Context, string, int, string, string) (bool, error) {
	panic("unexpected ExistsByHostPortAuth call")
}

func (s *proxyRepoStubForDefaultProxy) CountAccountsByProxyID(context.Context, int64) (int64, error) {
	panic("unexpected CountAccountsByProxyID call")
}

func (s *proxyRepoStubForDefaultProxy) ListAccountSummariesByProxyID(context.Context, int64) ([]ProxyAccountSummary, error) {
	panic("unexpected ListAccountSummariesByProxyID call")
}

func cloneAccountForDefaultProxyTest(account *Account) *Account {
	if account == nil {
		return nil
	}
	cloned := *account
	if account.ProxyID != nil {
		proxyID := *account.ProxyID
		cloned.ProxyID = &proxyID
	}
	if account.Notes != nil {
		notes := *account.Notes
		cloned.Notes = &notes
	}
	if account.RateMultiplier != nil {
		rate := *account.RateMultiplier
		cloned.RateMultiplier = &rate
	}
	if account.LoadFactor != nil {
		loadFactor := *account.LoadFactor
		cloned.LoadFactor = &loadFactor
	}
	cloned.GroupIDs = append([]int64(nil), account.GroupIDs...)
	return &cloned
}

func TestAdminService_CreateGroup_SetsDefaultProxyID(t *testing.T) {
	t.Parallel()

	groupRepo := &groupRepoStubForAdmin{}
	proxyID := int64(7)
	svc := &adminServiceImpl{
		groupRepo: groupRepo,
		proxyRepo: &proxyRepoStubForDefaultProxy{proxy: &Proxy{ID: proxyID}},
	}

	group, err := svc.CreateGroup(context.Background(), &CreateGroupInput{
		Name:           "openai-default",
		Platform:       PlatformOpenAI,
		RateMultiplier: 1,
		DefaultProxyID: &proxyID,
	})
	require.NoError(t, err)
	require.NotNil(t, group)
	require.NotNil(t, groupRepo.created)
	require.NotNil(t, groupRepo.created.DefaultProxyID)
	require.Equal(t, proxyID, *groupRepo.created.DefaultProxyID)
}

func TestAdminService_UpdateGroup_SetsDefaultProxyID(t *testing.T) {
	t.Parallel()

	existing := &Group{
		ID:       9,
		Name:     "openai-default",
		Platform: PlatformOpenAI,
		Status:   StatusActive,
	}
	groupRepo := &groupRepoStubForAdmin{getByID: existing}
	proxyID := int64(12)
	svc := &adminServiceImpl{
		groupRepo: groupRepo,
		proxyRepo: &proxyRepoStubForDefaultProxy{proxy: &Proxy{ID: proxyID}},
	}

	group, err := svc.UpdateGroup(context.Background(), existing.ID, &UpdateGroupInput{
		DefaultProxyID: &proxyID,
	})
	require.NoError(t, err)
	require.NotNil(t, group)
	require.NotNil(t, groupRepo.updated)
	require.NotNil(t, groupRepo.updated.DefaultProxyID)
	require.Equal(t, proxyID, *groupRepo.updated.DefaultProxyID)
}

func TestAdminService_CreateAccount_InheritsDefaultProxyFromGroup(t *testing.T) {
	t.Parallel()

	proxyID := int64(23)
	accountRepo := newAccountRepoStubForDefaultProxy()
	groupRepo := &groupRepoStubForAdmin{
		getByID: &Group{
			ID:             5,
			Name:           "openai-default",
			Platform:       PlatformOpenAI,
			Status:         StatusActive,
			DefaultProxyID: &proxyID,
		},
	}
	svc := &adminServiceImpl{
		accountRepo: accountRepo,
		groupRepo:   groupRepo,
	}

	account, err := svc.CreateAccount(context.Background(), &CreateAccountInput{
		Name:                  "acc",
		Platform:              PlatformOpenAI,
		Type:                  AccountTypeOAuth,
		Credentials:           map[string]any{"refresh_token": "rt"},
		Concurrency:           1,
		GroupIDs:              []int64{5},
		SkipMixedChannelCheck: true,
	})
	require.NoError(t, err)
	require.NotNil(t, account)
	require.NotNil(t, accountRepo.created)
	require.NotNil(t, accountRepo.created.ProxyID)
	require.Equal(t, proxyID, *accountRepo.created.ProxyID)
	require.Len(t, accountRepo.bindGroupsCalls, 1)
	require.Equal(t, []int64{5}, accountRepo.bindGroupsCalls[0].groupIDs)
}

func TestAdminService_UpdateAccount_InheritsDefaultProxyFromAssignedGroup(t *testing.T) {
	t.Parallel()

	proxyID := int64(31)
	accountRepo := newAccountRepoStubForDefaultProxy()
	accountRepo.accountsByID[11] = &Account{
		ID:          11,
		Name:        "acc",
		Platform:    PlatformOpenAI,
		Type:        AccountTypeOAuth,
		Credentials: map[string]any{"refresh_token": "rt"},
		Status:      StatusActive,
		Schedulable: true,
	}
	groupRepo := &groupRepoStubForAdmin{
		getByID: &Group{
			ID:             6,
			Name:           "openai-default",
			Platform:       PlatformOpenAI,
			Status:         StatusActive,
			DefaultProxyID: &proxyID,
		},
	}
	svc := &adminServiceImpl{
		accountRepo: accountRepo,
		groupRepo:   groupRepo,
	}

	groupIDs := []int64{6}
	account, err := svc.UpdateAccount(context.Background(), 11, &UpdateAccountInput{
		GroupIDs:              &groupIDs,
		SkipMixedChannelCheck: true,
	})
	require.NoError(t, err)
	require.NotNil(t, account)
	require.NotEmpty(t, accountRepo.updated)
	require.NotNil(t, accountRepo.updated[0].ProxyID)
	require.Equal(t, proxyID, *accountRepo.updated[0].ProxyID)
	require.Len(t, accountRepo.bindGroupsCalls, 1)
	require.Equal(t, []int64{6}, accountRepo.bindGroupsCalls[0].groupIDs)
}

func TestAdminService_BulkUpdateAccounts_InheritsDefaultProxyFromAssignedGroup(t *testing.T) {
	t.Parallel()

	proxyID := int64(41)
	accountRepo := newAccountRepoStubForDefaultProxy()
	accountRepo.accountsByID[101] = &Account{
		ID:          101,
		Name:        "acc-101",
		Platform:    PlatformOpenAI,
		Type:        AccountTypeOAuth,
		Credentials: map[string]any{"refresh_token": "rt-101"},
		Status:      StatusActive,
		Schedulable: true,
	}
	accountRepo.accountsByID[102] = &Account{
		ID:          102,
		Name:        "acc-102",
		Platform:    PlatformOpenAI,
		Type:        AccountTypeOAuth,
		Credentials: map[string]any{"refresh_token": "rt-102"},
		Status:      StatusActive,
		Schedulable: true,
		ProxyID:     int64PtrDefaultProxyTest(99),
	}
	groupRepo := &groupRepoStubForAdmin{
		getByID: &Group{
			ID:             8,
			Name:           "openai-default",
			Platform:       PlatformOpenAI,
			Status:         StatusActive,
			DefaultProxyID: &proxyID,
		},
	}
	svc := &adminServiceImpl{
		accountRepo: accountRepo,
		groupRepo:   groupRepo,
	}

	groupIDs := []int64{8}
	schedulable := true
	result, err := svc.BulkUpdateAccounts(context.Background(), &BulkUpdateAccountsInput{
		AccountIDs:            []int64{101, 102},
		GroupIDs:              &groupIDs,
		Schedulable:           &schedulable,
		SkipMixedChannelCheck: true,
	})
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, 2, result.Success)
	require.Len(t, accountRepo.bindGroupsCalls, 2)
	require.Len(t, accountRepo.updated, 1)
	require.Equal(t, int64(101), accountRepo.updated[0].ID)
	require.NotNil(t, accountRepo.updated[0].ProxyID)
	require.Equal(t, proxyID, *accountRepo.updated[0].ProxyID)
	require.Equal(t, int64(99), *accountRepo.accountsByID[102].ProxyID)
}

func int64PtrDefaultProxyTest(value int64) *int64 {
	return &value
}
