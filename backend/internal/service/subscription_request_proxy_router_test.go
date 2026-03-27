package service

import (
	"context"
	"testing"

	"github.com/B022MC/b022hub/internal/pkg/pagination"
)

type subscriptionProxyRepoStub struct {
	proxies []Proxy
}

func (s *subscriptionProxyRepoStub) Create(context.Context, *Proxy) error { panic("not implemented") }
func (s *subscriptionProxyRepoStub) GetByID(context.Context, int64) (*Proxy, error) {
	panic("not implemented")
}
func (s *subscriptionProxyRepoStub) ListByIDs(context.Context, []int64) ([]Proxy, error) {
	panic("not implemented")
}
func (s *subscriptionProxyRepoStub) Update(context.Context, *Proxy) error { panic("not implemented") }
func (s *subscriptionProxyRepoStub) Delete(context.Context, int64) error  { panic("not implemented") }
func (s *subscriptionProxyRepoStub) List(context.Context, pagination.PaginationParams) ([]Proxy, *pagination.PaginationResult, error) {
	panic("not implemented")
}
func (s *subscriptionProxyRepoStub) ListWithFilters(context.Context, pagination.PaginationParams, string, string, string) ([]Proxy, *pagination.PaginationResult, error) {
	panic("not implemented")
}
func (s *subscriptionProxyRepoStub) ListWithFiltersAndAccountCount(context.Context, pagination.PaginationParams, string, string, string) ([]ProxyWithAccountCount, *pagination.PaginationResult, error) {
	panic("not implemented")
}
func (s *subscriptionProxyRepoStub) ListActive(context.Context) ([]Proxy, error) {
	return append([]Proxy(nil), s.proxies...), nil
}
func (s *subscriptionProxyRepoStub) ListActiveWithAccountCount(context.Context) ([]ProxyWithAccountCount, error) {
	panic("not implemented")
}
func (s *subscriptionProxyRepoStub) ExistsByHostPortAuth(context.Context, string, int, string, string) (bool, error) {
	panic("not implemented")
}
func (s *subscriptionProxyRepoStub) CountAccountsByProxyID(context.Context, int64) (int64, error) {
	panic("not implemented")
}
func (s *subscriptionProxyRepoStub) ListAccountSummariesByProxyID(context.Context, int64) ([]ProxyAccountSummary, error) {
	panic("not implemented")
}

func TestFilterSubscriptionScopedProxiesPrefersNamedPool(t *testing.T) {
	input := []Proxy{
		{Name: "default", Protocol: "socks5h", Host: "normal", Port: 1080},
		{Name: "subscription: rotating", Protocol: "socks5h", Host: "sub", Port: 1081},
		{Name: "sub-pool-2", Protocol: "socks5h", Host: "sub2", Port: 1082},
	}

	got := filterSubscriptionScopedProxies(input)
	if len(got) != 2 {
		t.Fatalf("expected 2 scoped proxies, got %d", len(got))
	}
	if got[0].Name != "subscription: rotating" || got[1].Name != "sub-pool-2" {
		t.Fatalf("unexpected scoped proxies: %+v", got)
	}
}

func TestFilterSubscriptionScopedProxiesFallsBackToAll(t *testing.T) {
	input := []Proxy{
		{Name: "default", Protocol: "socks5h", Host: "normal", Port: 1080},
		{Name: "backup", Protocol: "http", Host: "normal2", Port: 8080},
	}

	got := filterSubscriptionScopedProxies(input)
	if len(got) != len(input) {
		t.Fatalf("expected fallback to all proxies, got %d", len(got))
	}
}

func TestSubscriptionRouterListActiveProxyURLsUsesScopedPool(t *testing.T) {
	repo := &subscriptionProxyRepoStub{
		proxies: []Proxy{
			{Name: "default", Protocol: "socks5h", Host: "normal", Port: 1080},
			{Name: "subscription: rotating", Protocol: "socks5h", Host: "sub", Port: 1081},
		},
	}

	router := NewSubscriptionRequestProxyRouter(repo)
	urls, err := router.listActiveProxyURLs(context.Background())
	if err != nil {
		t.Fatalf("listActiveProxyURLs returned error: %v", err)
	}
	if len(urls) != 1 {
		t.Fatalf("expected 1 scoped URL, got %d", len(urls))
	}
	if urls[0] != "socks5h://sub:1081" {
		t.Fatalf("unexpected proxy URL: %q", urls[0])
	}
}
