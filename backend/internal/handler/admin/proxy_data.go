package admin

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/B022MC/b022hub/internal/pkg/response"
	"github.com/B022MC/b022hub/internal/service"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
)

const proxySubscriptionFetchTimeout = 15 * time.Second

// ExportData exports proxy-only data for migration.
func (h *ProxyHandler) ExportData(c *gin.Context) {
	ctx := c.Request.Context()

	selectedIDs, err := parseProxyIDs(c)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	var proxies []service.Proxy
	if len(selectedIDs) > 0 {
		proxies, err = h.getProxiesByIDs(ctx, selectedIDs)
		if err != nil {
			response.ErrorFrom(c, err)
			return
		}
	} else {
		protocol := c.Query("protocol")
		status := c.Query("status")
		search := strings.TrimSpace(c.Query("search"))
		if len(search) > 100 {
			search = search[:100]
		}

		proxies, err = h.listProxiesFiltered(ctx, protocol, status, search)
		if err != nil {
			response.ErrorFrom(c, err)
			return
		}
	}

	dataProxies := make([]DataProxy, 0, len(proxies))
	for i := range proxies {
		p := proxies[i]
		key := buildProxyKey(p.Protocol, p.Host, p.Port, p.Username, p.Password)
		dataProxies = append(dataProxies, DataProxy{
			ProxyKey: key,
			Name:     p.Name,
			Protocol: p.Protocol,
			Host:     p.Host,
			Port:     p.Port,
			Username: p.Username,
			Password: p.Password,
			Status:   p.Status,
		})
	}

	payload := DataPayload{
		ExportedAt: time.Now().UTC().Format(time.RFC3339),
		Proxies:    dataProxies,
		Accounts:   []DataAccount{},
	}

	response.Success(c, payload)
}

// ImportData imports proxy-only data for migration.
func (h *ProxyHandler) ImportData(c *gin.Context) {
	type ProxyImportRequest struct {
		Data            *DataPayload `json:"data"`
		SubscriptionURL string       `json:"subscription_url"`
	}

	var req ProxyImportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	hasData := req.Data != nil
	hasSubscriptionURL := strings.TrimSpace(req.SubscriptionURL) != ""
	if hasData == hasSubscriptionURL {
		response.BadRequest(c, "Provide either data or subscription_url")
		return
	}

	ctx := c.Request.Context()
	result := DataImportResult{}

	var payload DataPayload
	if hasSubscriptionURL {
		subscriptionPayload, err := buildProxyDataPayloadFromSubscription(ctx, req.SubscriptionURL)
		if err != nil {
			response.BadRequest(c, err.Error())
			return
		}
		payload = subscriptionPayload.Payload
		result.ProxyFailed += subscriptionPayload.FailedCount
		result.Errors = append(result.Errors, subscriptionPayload.Errors...)
	} else {
		payload = *req.Data
	}

	if err := validateDataHeader(payload); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	importResult, err := h.importProxyData(ctx, payload)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	result.ProxyCreated += importResult.ProxyCreated
	result.ProxyReused += importResult.ProxyReused
	result.ProxyFailed += importResult.ProxyFailed
	result.AccountCreated += importResult.AccountCreated
	result.AccountFailed += importResult.AccountFailed
	result.Errors = append(result.Errors, importResult.Errors...)

	response.Success(c, result)
}

func (h *ProxyHandler) importProxyData(ctx context.Context, payload DataPayload) (DataImportResult, error) {
	result := DataImportResult{}

	existingProxies, err := h.listProxiesFiltered(ctx, "", "", "")
	if err != nil {
		return result, err
	}

	proxyByKey := make(map[string]service.Proxy, len(existingProxies))
	for i := range existingProxies {
		p := existingProxies[i]
		key := buildProxyKey(p.Protocol, p.Host, p.Port, p.Username, p.Password)
		proxyByKey[key] = p
	}

	latencyProbeIDs := make([]int64, 0, len(payload.Proxies))
	for i := range payload.Proxies {
		item := payload.Proxies[i]
		key := item.ProxyKey
		if key == "" {
			key = buildProxyKey(item.Protocol, item.Host, item.Port, item.Username, item.Password)
		}

		if err := validateDataProxy(item); err != nil {
			result.ProxyFailed++
			result.Errors = append(result.Errors, DataImportError{
				Kind:     "proxy",
				Name:     item.Name,
				ProxyKey: key,
				Message:  err.Error(),
			})
			continue
		}

		normalizedStatus := normalizeProxyStatus(item.Status)
		if existing, ok := proxyByKey[key]; ok {
			result.ProxyReused++
			if normalizedStatus != "" && normalizedStatus != existing.Status {
				if _, err := h.adminService.UpdateProxy(ctx, existing.ID, &service.UpdateProxyInput{Status: normalizedStatus}); err != nil {
					result.Errors = append(result.Errors, DataImportError{
						Kind:     "proxy",
						Name:     item.Name,
						ProxyKey: key,
						Message:  "update status failed: " + err.Error(),
					})
				}
			}
			latencyProbeIDs = append(latencyProbeIDs, existing.ID)
			continue
		}

		created, err := h.adminService.CreateProxy(ctx, &service.CreateProxyInput{
			Name:     defaultProxyName(item.Name),
			Protocol: item.Protocol,
			Host:     item.Host,
			Port:     item.Port,
			Username: item.Username,
			Password: item.Password,
		})
		if err != nil {
			result.ProxyFailed++
			result.Errors = append(result.Errors, DataImportError{
				Kind:     "proxy",
				Name:     item.Name,
				ProxyKey: key,
				Message:  err.Error(),
			})
			continue
		}
		result.ProxyCreated++
		proxyByKey[key] = *created

		if normalizedStatus != "" && normalizedStatus != created.Status {
			if _, err := h.adminService.UpdateProxy(ctx, created.ID, &service.UpdateProxyInput{Status: normalizedStatus}); err != nil {
				result.Errors = append(result.Errors, DataImportError{
					Kind:     "proxy",
					Name:     item.Name,
					ProxyKey: key,
					Message:  "update status failed: " + err.Error(),
				})
			}
		}
		// CreateProxy already triggers a latency probe, avoid double probing here.
	}

	if len(latencyProbeIDs) > 0 {
		ids := append([]int64(nil), latencyProbeIDs...)
		go func() {
			for _, id := range ids {
				_, _ = h.adminService.TestProxy(context.Background(), id)
			}
		}()
	}

	return result, nil
}

type proxySubscriptionImportPayload struct {
	Payload     DataPayload
	FailedCount int
	Errors      []DataImportError
}

type clashSubscriptionDocument struct {
	Proxies []map[string]any `yaml:"proxies"`
}

func buildProxyDataPayloadFromSubscription(ctx context.Context, rawURL string) (proxySubscriptionImportPayload, error) {
	subscriptionURL, err := parseSubscriptionURL(rawURL)
	if err != nil {
		return proxySubscriptionImportPayload{}, err
	}

	body, err := fetchSubscriptionBody(ctx, subscriptionURL)
	if err != nil {
		return proxySubscriptionImportPayload{}, err
	}

	var doc clashSubscriptionDocument
	if err := yaml.Unmarshal(body, &doc); err != nil {
		return proxySubscriptionImportPayload{}, fmt.Errorf("failed to parse clash subscription: %w", err)
	}
	if len(doc.Proxies) == 0 {
		return proxySubscriptionImportPayload{}, fmt.Errorf("subscription contains no proxies")
	}

	result := proxySubscriptionImportPayload{
		Payload: DataPayload{
			Type:       dataType,
			Version:    dataVersion,
			ExportedAt: time.Now().UTC().Format(time.RFC3339),
			Proxies:    make([]DataProxy, 0, len(doc.Proxies)),
			Accounts:   []DataAccount{},
		},
	}

	unsupportedCounts := make(map[string]int)
	for i := range doc.Proxies {
		item, supported, err := proxyFromClashNode(doc.Proxies[i])
		if err != nil {
			result.FailedCount++
			result.Errors = append(result.Errors, DataImportError{
				Kind:    "proxy",
				Name:    strings.TrimSpace(stringValue(doc.Proxies[i]["name"])),
				Message: err.Error(),
			})
			continue
		}
		if !supported {
			nodeType := strings.ToLower(strings.TrimSpace(stringValue(doc.Proxies[i]["type"])))
			if nodeType == "" {
				nodeType = "unknown"
			}
			unsupportedCounts[nodeType]++
			continue
		}
		result.Payload.Proxies = append(result.Payload.Proxies, item)
	}

	if len(unsupportedCounts) > 0 {
		unsupportedTypes := make([]string, 0, len(unsupportedCounts))
		for nodeType := range unsupportedCounts {
			unsupportedTypes = append(unsupportedTypes, nodeType)
		}
		sort.Strings(unsupportedTypes)
		for _, nodeType := range unsupportedTypes {
			skipped := unsupportedCounts[nodeType]
			result.FailedCount += skipped
			result.Errors = append(result.Errors, DataImportError{
				Kind:    "proxy",
				Message: fmt.Sprintf("subscription node type %q is not supported yet; only http/https/socks5/socks5h can be imported (%d skipped)", nodeType, skipped),
			})
		}
	}

	if len(result.Payload.Proxies) == 0 && result.FailedCount == 0 {
		return proxySubscriptionImportPayload{}, fmt.Errorf("subscription contains no importable proxies")
	}

	return result, nil
}

func parseSubscriptionURL(rawURL string) (*url.URL, error) {
	trimmed := strings.TrimSpace(rawURL)
	if trimmed == "" {
		return nil, fmt.Errorf("subscription_url is required")
	}
	parsed, err := url.Parse(trimmed)
	if err != nil {
		return nil, fmt.Errorf("invalid subscription_url: %w", err)
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return nil, fmt.Errorf("subscription_url must start with http:// or https://")
	}
	if parsed.Host == "" {
		return nil, fmt.Errorf("subscription_url host is required")
	}
	return parsed, nil
}

func fetchSubscriptionBody(ctx context.Context, subscriptionURL *url.URL) ([]byte, error) {
	client := &http.Client{Timeout: proxySubscriptionFetchTimeout}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, subscriptionURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to build subscription request: %w", err)
	}
	req.Header.Set("Accept", "application/yaml, text/yaml, text/plain, */*")
	req.Header.Set("User-Agent", "sub2api-proxy-import/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch subscription: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("subscription request failed with status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, 8<<20))
	if err != nil {
		return nil, fmt.Errorf("failed to read subscription body: %w", err)
	}
	if len(body) == 0 {
		return nil, fmt.Errorf("subscription response is empty")
	}
	return body, nil
}

func proxyFromClashNode(node map[string]any) (DataProxy, bool, error) {
	nodeType := strings.ToLower(strings.TrimSpace(stringValue(node["type"])))
	protocol, supported := clashProxyProtocol(nodeType, node)
	if !supported {
		return DataProxy{}, false, nil
	}

	port, err := intValue(node["port"])
	if err != nil {
		return DataProxy{}, true, fmt.Errorf("invalid clash node port for %q: %w", strings.TrimSpace(stringValue(node["name"])), err)
	}

	name := strings.TrimSpace(stringValue(node["name"]))
	host := strings.TrimSpace(stringValue(node["server"]))
	username := strings.TrimSpace(stringValue(node["username"]))
	password := strings.TrimSpace(stringValue(node["password"]))

	item := DataProxy{
		ProxyKey: buildProxyKey(protocol, host, port, username, password),
		Name:     defaultProxyName(name),
		Protocol: protocol,
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
		Status:   service.StatusActive,
	}
	if err := validateDataProxy(item); err != nil {
		return DataProxy{}, true, err
	}
	return item, true, nil
}

func clashProxyProtocol(nodeType string, node map[string]any) (string, bool) {
	switch nodeType {
	case "http":
		if boolValue(node["tls"]) {
			return "https", true
		}
		return "http", true
	case "https":
		return "https", true
	case "socks5", "socks5h":
		return nodeType, true
	default:
		return "", false
	}
}

func stringValue(value any) string {
	switch v := value.(type) {
	case string:
		return v
	case fmt.Stringer:
		return v.String()
	case nil:
		return ""
	default:
		return fmt.Sprint(v)
	}
}

func intValue(value any) (int, error) {
	switch v := value.(type) {
	case int:
		return v, nil
	case int8:
		return int(v), nil
	case int16:
		return int(v), nil
	case int32:
		return int(v), nil
	case int64:
		return int(v), nil
	case uint:
		return int(v), nil
	case uint8:
		return int(v), nil
	case uint16:
		return int(v), nil
	case uint32:
		return int(v), nil
	case uint64:
		return int(v), nil
	case float32:
		return int(v), nil
	case float64:
		return int(v), nil
	case string:
		trimmed := strings.TrimSpace(v)
		if trimmed == "" {
			return 0, fmt.Errorf("empty port")
		}
		port, err := strconv.Atoi(trimmed)
		if err != nil {
			return 0, err
		}
		return port, nil
	default:
		return 0, fmt.Errorf("unsupported port type %T", value)
	}
}

func boolValue(value any) bool {
	switch v := value.(type) {
	case bool:
		return v
	case string:
		return strings.EqualFold(strings.TrimSpace(v), "true")
	default:
		return false
	}
}

func (h *ProxyHandler) getProxiesByIDs(ctx context.Context, ids []int64) ([]service.Proxy, error) {
	if len(ids) == 0 {
		return []service.Proxy{}, nil
	}
	return h.adminService.GetProxiesByIDs(ctx, ids)
}

func parseProxyIDs(c *gin.Context) ([]int64, error) {
	values := c.QueryArray("ids")
	if len(values) == 0 {
		raw := strings.TrimSpace(c.Query("ids"))
		if raw != "" {
			values = []string{raw}
		}
	}
	if len(values) == 0 {
		return nil, nil
	}

	ids := make([]int64, 0, len(values))
	for _, item := range values {
		for _, part := range strings.Split(item, ",") {
			part = strings.TrimSpace(part)
			if part == "" {
				continue
			}
			id, err := strconv.ParseInt(part, 10, 64)
			if err != nil || id <= 0 {
				return nil, fmt.Errorf("invalid proxy id: %s", part)
			}
			ids = append(ids, id)
		}
	}
	return ids, nil
}

func (h *ProxyHandler) listProxiesFiltered(ctx context.Context, protocol, status, search string) ([]service.Proxy, error) {
	page := 1
	pageSize := dataPageCap
	var out []service.Proxy
	for {
		items, total, err := h.adminService.ListProxies(ctx, page, pageSize, protocol, status, search)
		if err != nil {
			return nil, err
		}
		out = append(out, items...)
		if len(out) >= int(total) || len(items) == 0 {
			break
		}
		page++
	}
	return out, nil
}
