package service

import (
	"context"
	cryptorand "crypto/rand"
	"math/big"
	mathrand "math/rand"
	"sort"
	"strings"
	"sync"
	"time"
)

const defaultSubscriptionRequestProxyCacheTTL = 10 * time.Second

var subscriptionProxyNamePrefixes = []string{
	"subscription:",
	"subscription-",
	"sub:",
	"sub-",
	"[subscription]",
}

type requestProxyURLContextKey struct{}

// WithRequestProxyURL stores the resolved upstream proxy URL in the request context.
func WithRequestProxyURL(ctx context.Context, proxyURL string) context.Context {
	trimmed := strings.TrimSpace(proxyURL)
	if trimmed == "" {
		return ctx
	}
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, requestProxyURLContextKey{}, trimmed)
}

// RequestProxyURLFromContext returns the request-scoped upstream proxy URL if present.
func RequestProxyURLFromContext(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	proxyURL, _ := ctx.Value(requestProxyURLContextKey{}).(string)
	return strings.TrimSpace(proxyURL)
}

// ResolveUpstreamProxyURL prefers the request-scoped proxy selection, then falls back to the account proxy.
func ResolveUpstreamProxyURL(ctx context.Context, account *Account) string {
	if proxyURL := RequestProxyURLFromContext(ctx); proxyURL != "" {
		return proxyURL
	}
	if account == nil || account.ProxyID == nil || account.Proxy == nil {
		return ""
	}
	return strings.TrimSpace(account.Proxy.URL())
}

// SubscriptionRequestProxyRouter resolves a random active proxy for subscription requests.
type SubscriptionRequestProxyRouter struct {
	proxyRepo ProxyRepository
	cacheTTL  time.Duration
	now       func() time.Time

	mu         sync.RWMutex
	cachedURLs []string
	expiresAt  time.Time
}

func NewSubscriptionRequestProxyRouter(proxyRepo ProxyRepository) *SubscriptionRequestProxyRouter {
	return &SubscriptionRequestProxyRouter{
		proxyRepo: proxyRepo,
		cacheTTL:  defaultSubscriptionRequestProxyCacheTTL,
		now:       time.Now,
	}
}

// BindRandomProxy stores a random active proxy in the request context.
// If no active global proxy is available, the original context is returned unchanged.
func (r *SubscriptionRequestProxyRouter) BindRandomProxy(ctx context.Context) context.Context {
	if r == nil || r.proxyRepo == nil {
		return ctx
	}
	if RequestProxyURLFromContext(ctx) != "" {
		return ctx
	}
	proxyURL, err := r.pickRandomProxyURL(ctx)
	if err != nil || proxyURL == "" {
		return ctx
	}
	return WithRequestProxyURL(ctx, proxyURL)
}

func (r *SubscriptionRequestProxyRouter) pickRandomProxyURL(ctx context.Context) (string, error) {
	urls, err := r.listActiveProxyURLs(ctx)
	if err != nil || len(urls) == 0 {
		return "", err
	}
	if len(urls) == 1 {
		return urls[0], nil
	}

	indexLimit := big.NewInt(int64(len(urls)))
	index, err := cryptorand.Int(cryptorand.Reader, indexLimit)
	if err == nil {
		return urls[index.Int64()], nil
	}

	return urls[mathrand.Intn(len(urls))], nil
}

func (r *SubscriptionRequestProxyRouter) listActiveProxyURLs(ctx context.Context) ([]string, error) {
	now := r.now()

	r.mu.RLock()
	if now.Before(r.expiresAt) {
		cached := append([]string(nil), r.cachedURLs...)
		r.mu.RUnlock()
		return cached, nil
	}
	r.mu.RUnlock()

	r.mu.Lock()
	defer r.mu.Unlock()

	now = r.now()
	if now.Before(r.expiresAt) {
		return append([]string(nil), r.cachedURLs...), nil
	}

	proxies, err := r.proxyRepo.ListActive(ctx)
	if err != nil {
		return nil, err
	}

	proxies = filterSubscriptionScopedProxies(proxies)

	seen := make(map[string]struct{}, len(proxies))
	urls := make([]string, 0, len(proxies))
	for i := range proxies {
		proxyURL := strings.TrimSpace(proxies[i].URL())
		if proxyURL == "" {
			continue
		}
		if _, ok := seen[proxyURL]; ok {
			continue
		}
		seen[proxyURL] = struct{}{}
		urls = append(urls, proxyURL)
	}
	sort.Strings(urls)

	r.cachedURLs = urls
	r.expiresAt = now.Add(r.cacheTTL)
	return append([]string(nil), urls...), nil
}

func filterSubscriptionScopedProxies(proxies []Proxy) []Proxy {
	if len(proxies) == 0 {
		return proxies
	}

	scoped := make([]Proxy, 0, len(proxies))
	for i := range proxies {
		name := strings.ToLower(strings.TrimSpace(proxies[i].Name))
		for _, prefix := range subscriptionProxyNamePrefixes {
			if strings.HasPrefix(name, prefix) {
				scoped = append(scoped, proxies[i])
				break
			}
		}
	}
	if len(scoped) > 0 {
		return scoped
	}
	return proxies
}
