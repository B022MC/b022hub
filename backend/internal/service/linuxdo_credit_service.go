package service

import (
	"context"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	dbent "github.com/B022MC/b022hub/ent"
	infraerrors "github.com/B022MC/b022hub/internal/pkg/errors"
)

const (
	linuxDoCreditGatewayBaseURL = "https://credit.linux.do/epay"
	linuxDoCreditCheckoutPath   = "/pay/submit.php"
	linuxDoCreditQueryPath      = "/api.php"
	linuxDoCreditFixedType      = "epay"
	linuxDoCreditSignType       = "MD5"
)

type PaymentOrderRepository interface {
	Create(ctx context.Context, order *PaymentOrder) error
	GetByOutTradeNo(ctx context.Context, outTradeNo string) (*PaymentOrder, error)
	GetByUserAndOutTradeNo(ctx context.Context, userID int64, outTradeNo string) (*PaymentOrder, error)
	ListByUser(ctx context.Context, userID int64, limit int) ([]PaymentOrder, error)
	MarkPaid(ctx context.Context, outTradeNo string, providerTradeNo string, rawProviderPayload string, paidAt time.Time) (*PaymentOrder, bool, error)
}

type LinuxDoCreditService struct {
	entClient            *dbent.Client
	orderRepo            PaymentOrderRepository
	userRepo             UserRepository
	redeemRepo           RedeemCodeRepository
	billingCacheService  *BillingCacheService
	authCacheInvalidator APIKeyAuthCacheInvalidator
	settingService       *SettingService
	httpClient           *http.Client
}

func NewLinuxDoCreditService(
	entClient *dbent.Client,
	orderRepo PaymentOrderRepository,
	userRepo UserRepository,
	redeemRepo RedeemCodeRepository,
	billingCacheService *BillingCacheService,
	authCacheInvalidator APIKeyAuthCacheInvalidator,
	settingService *SettingService,
) *LinuxDoCreditService {
	return &LinuxDoCreditService{
		entClient:            entClient,
		orderRepo:            orderRepo,
		userRepo:             userRepo,
		redeemRepo:           redeemRepo,
		billingCacheService:  billingCacheService,
		authCacheInvalidator: authCacheInvalidator,
		settingService:       settingService,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

type linuxDoCreditQueryResponse struct {
	Code       int    `json:"code"`
	Msg        string `json:"msg"`
	TradeNo    string `json:"trade_no"`
	OutTradeNo string `json:"out_trade_no"`
	Type       string `json:"type"`
	PID        string `json:"pid"`
	Name       string `json:"name"`
	Money      string `json:"money"`
	Status     int    `json:"status"`
}

func (s *LinuxDoCreditService) CreateCheckout(ctx context.Context, userID int64, amount float64) (*LinuxDoCreditCheckout, error) {
	settings, err := s.settingService.GetAllSettings(ctx)
	if err != nil {
		return nil, err
	}

	cfg := linuxDoCreditConfigFromSettings(settings)
	if !cfg.Enabled {
		return nil, ErrLinuxDoCreditDisabled
	}
	if cfg.ClientID == "" || cfg.ClientSecret == "" || cfg.ExchangeRate <= 0 {
		return nil, ErrLinuxDoCreditBadConfig
	}

	amount = normalizeLinuxDoCreditAmount(amount)
	if amount <= 0 {
		return nil, ErrLinuxDoCreditInvalidAmt
	}

	outTradeNo, err := generateLinuxDoCreditOrderNo(userID)
	if err != nil {
		return nil, infraerrors.InternalServer("LINUXDO_CREDIT_ORDER_NO_FAILED", "failed to generate payment order number").WithCause(err)
	}

	title := buildLinuxDoCreditOrderTitle(settings.SiteName)
	order := &PaymentOrder{
		Provider:       PaymentProviderLinuxDoCredit,
		OutTradeNo:     outTradeNo,
		UserID:         userID,
		Title:          title,
		Amount:         amount,
		CreditedAmount: normalizeLinuxDoCreditCredit(amount * cfg.ExchangeRate),
		Status:         PaymentOrderStatusPending,
	}

	if err := s.orderRepo.Create(ctx, order); err != nil {
		return nil, infraerrors.InternalServer("LINUXDO_CREDIT_ORDER_CREATE_FAILED", "failed to create payment order").WithCause(err)
	}

	fields := map[string]string{
		"pid":          cfg.ClientID,
		"type":         linuxDoCreditFixedType,
		"out_trade_no": order.OutTradeNo,
		"name":         order.Title,
		"money":        formatLinuxDoCreditAmount(order.Amount),
		"sign_type":    linuxDoCreditSignType,
	}
	fields["sign"] = signLinuxDoCreditFields(fields, cfg.ClientSecret)

	return &LinuxDoCreditCheckout{
		Order:       order,
		CheckoutURL: linuxDoCreditGatewayBaseURL + linuxDoCreditCheckoutPath,
		Fields:      fields,
	}, nil
}

func (s *LinuxDoCreditService) ListOrders(ctx context.Context, userID int64, limit int) ([]PaymentOrder, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 50 {
		limit = 50
	}
	return s.orderRepo.ListByUser(ctx, userID, limit)
}

func (s *LinuxDoCreditService) GetOrder(ctx context.Context, userID int64, outTradeNo string, syncRemote bool) (*PaymentOrder, error) {
	order, err := s.orderRepo.GetByUserAndOutTradeNo(ctx, userID, strings.TrimSpace(outTradeNo))
	if err != nil {
		return nil, err
	}

	if syncRemote && order.Status == PaymentOrderStatusPending {
		if synced, syncErr := s.SyncOrderStatus(ctx, order.OutTradeNo); syncErr == nil && synced != nil {
			if synced.UserID != userID {
				return nil, ErrLinuxDoCreditOrderDenied
			}
			order = synced
		}
	}

	return order, nil
}

func (s *LinuxDoCreditService) SyncOrderStatus(ctx context.Context, outTradeNo string) (*PaymentOrder, error) {
	order, err := s.orderRepo.GetByOutTradeNo(ctx, strings.TrimSpace(outTradeNo))
	if err != nil {
		return nil, err
	}
	if order.Status == PaymentOrderStatusPaid {
		return order, nil
	}

	cfg, err := s.getLinuxDoCreditConfig(ctx, false)
	if err != nil {
		return nil, err
	}

	result, rawBody, err := s.queryProviderOrder(ctx, cfg, order.OutTradeNo)
	if err != nil {
		return nil, err
	}

	if result.Status == 1 && strings.EqualFold(result.Type, linuxDoCreditFixedType) {
		if err := s.applySuccessfulOrder(ctx, order, result.TradeNo, rawBody); err != nil {
			return nil, err
		}
		return s.orderRepo.GetByOutTradeNo(ctx, order.OutTradeNo)
	}

	return order, nil
}

func (s *LinuxDoCreditService) HandleNotify(ctx context.Context, params url.Values) error {
	cfg, err := s.getLinuxDoCreditConfig(ctx, false)
	if err != nil {
		return err
	}

	flat := flattenLinuxDoCreditValues(params)
	sign := strings.TrimSpace(flat["sign"])
	if sign == "" {
		return ErrLinuxDoCreditBadSign
	}
	if !strings.EqualFold(sign, signLinuxDoCreditFields(flat, cfg.ClientSecret)) {
		return ErrLinuxDoCreditBadSign
	}

	if flat["pid"] != "" && cfg.ClientID != "" && flat["pid"] != cfg.ClientID {
		return infraerrors.BadRequest("LINUXDO_CREDIT_PID_MISMATCH", "linuxdo credit client id mismatch")
	}
	if !strings.EqualFold(flat["type"], linuxDoCreditFixedType) {
		return infraerrors.BadRequest("LINUXDO_CREDIT_TYPE_INVALID", "unsupported linuxdo credit payment type")
	}
	if !strings.EqualFold(flat["trade_status"], "TRADE_SUCCESS") {
		return infraerrors.BadRequest("LINUXDO_CREDIT_TRADE_STATUS_INVALID", "linuxdo credit trade is not successful")
	}

	order, err := s.orderRepo.GetByOutTradeNo(ctx, flat["out_trade_no"])
	if err != nil {
		return err
	}

	amount, err := parseLinuxDoCreditAmount(flat["money"])
	if err != nil {
		return infraerrors.BadRequest("LINUXDO_CREDIT_AMOUNT_INVALID", "invalid linuxdo credit amount").WithCause(err)
	}
	if !sameLinuxDoCreditAmount(order.Amount, amount) {
		return infraerrors.BadRequest("LINUXDO_CREDIT_AMOUNT_MISMATCH", "linuxdo credit amount mismatch")
	}

	rawPayload, err := marshalLinuxDoCreditValues(flat)
	if err != nil {
		return infraerrors.InternalServer("LINUXDO_CREDIT_NOTIFY_PAYLOAD_FAILED", "failed to marshal linuxdo credit payload").WithCause(err)
	}

	return s.applySuccessfulOrder(ctx, order, flat["trade_no"], rawPayload)
}

func (s *LinuxDoCreditService) applySuccessfulOrder(ctx context.Context, order *PaymentOrder, providerTradeNo string, rawPayload string) error {
	if order == nil {
		return ErrPaymentOrderNotFound
	}
	if order.Status == PaymentOrderStatusPaid {
		return nil
	}
	if s.entClient == nil {
		return infraerrors.InternalServer("LINUXDO_CREDIT_TX_UNAVAILABLE", "payment transaction support is unavailable")
	}

	tx, err := s.entClient.Tx(ctx)
	if err != nil {
		return infraerrors.InternalServer("LINUXDO_CREDIT_TX_BEGIN_FAILED", "failed to begin payment transaction").WithCause(err)
	}
	defer func() { _ = tx.Rollback() }()

	txCtx := dbent.NewTxContext(ctx, tx)
	updatedOrder, changed, err := s.orderRepo.MarkPaid(txCtx, order.OutTradeNo, strings.TrimSpace(providerTradeNo), rawPayload, time.Now())
	if err != nil {
		return infraerrors.InternalServer("LINUXDO_CREDIT_ORDER_UPDATE_FAILED", "failed to update payment order").WithCause(err)
	}
	if !changed {
		if err := tx.Commit(); err != nil {
			return infraerrors.InternalServer("LINUXDO_CREDIT_TX_COMMIT_FAILED", "failed to commit payment transaction").WithCause(err)
		}
		return nil
	}

	if err := s.userRepo.UpdateBalance(txCtx, updatedOrder.UserID, updatedOrder.CreditedAmount); err != nil {
		return infraerrors.InternalServer("LINUXDO_CREDIT_BALANCE_UPDATE_FAILED", "failed to update user balance").WithCause(err)
	}

	usedBy := updatedOrder.UserID
	usedAt := time.Now()
	record := &RedeemCode{
		Code:   linuxDoCreditHistoryCode(updatedOrder.OutTradeNo),
		Type:   AdjustmentTypeAdminBalance,
		Value:  updatedOrder.CreditedAmount,
		Status: StatusUsed,
		UsedBy: &usedBy,
		UsedAt: &usedAt,
		Notes:  fmt.Sprintf("LINUX DO Credit order %s", updatedOrder.OutTradeNo),
	}
	if err := s.redeemRepo.Create(txCtx, record); err != nil {
		return infraerrors.InternalServer("LINUXDO_CREDIT_HISTORY_CREATE_FAILED", "failed to create payment history record").WithCause(err)
	}

	if err := tx.Commit(); err != nil {
		return infraerrors.InternalServer("LINUXDO_CREDIT_TX_COMMIT_FAILED", "failed to commit payment transaction").WithCause(err)
	}

	if s.authCacheInvalidator != nil {
		s.authCacheInvalidator.InvalidateAuthCacheByUserID(ctx, updatedOrder.UserID)
	}
	if s.billingCacheService != nil {
		cacheCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = s.billingCacheService.InvalidateUserBalance(cacheCtx, updatedOrder.UserID)
	}

	return nil
}

func (s *LinuxDoCreditService) queryProviderOrder(ctx context.Context, cfg LinuxDoCreditConfig, outTradeNo string) (*linuxDoCreditQueryResponse, string, error) {
	u, err := url.Parse(linuxDoCreditGatewayBaseURL + linuxDoCreditQueryPath)
	if err != nil {
		return nil, "", infraerrors.InternalServer("LINUXDO_CREDIT_QUERY_URL_INVALID", "failed to prepare linuxdo credit query url").WithCause(err)
	}

	q := u.Query()
	q.Set("act", "order")
	q.Set("pid", cfg.ClientID)
	q.Set("key", cfg.ClientSecret)
	q.Set("out_trade_no", outTradeNo)
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, "", infraerrors.InternalServer("LINUXDO_CREDIT_QUERY_BUILD_FAILED", "failed to build linuxdo credit query request").WithCause(err)
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, "", infraerrors.ServiceUnavailable("LINUXDO_CREDIT_QUERY_FAILED", "failed to query linuxdo credit order").WithCause(err)
	}
	defer resp.Body.Close()

	var result linuxDoCreditQueryResponse
	bodyBytes, readErr := readHTTPBody(resp.Body)
	if readErr != nil {
		return nil, "", infraerrors.ServiceUnavailable("LINUXDO_CREDIT_QUERY_READ_FAILED", "failed to read linuxdo credit query response").WithCause(readErr)
	}

	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return nil, string(bodyBytes), infraerrors.ServiceUnavailable("LINUXDO_CREDIT_QUERY_PARSE_FAILED", "failed to parse linuxdo credit query response").WithCause(err)
	}

	if resp.StatusCode == http.StatusNotFound {
		return &result, string(bodyBytes), nil
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, string(bodyBytes), infraerrors.ServiceUnavailable("LINUXDO_CREDIT_QUERY_STATUS_FAILED", "linuxdo credit query request failed")
	}

	return &result, string(bodyBytes), nil
}

func (s *LinuxDoCreditService) getLinuxDoCreditConfig(ctx context.Context, requireEnabled bool) (LinuxDoCreditConfig, error) {
	settings, err := s.settingService.GetAllSettings(ctx)
	if err != nil {
		return LinuxDoCreditConfig{}, err
	}

	cfg := linuxDoCreditConfigFromSettings(settings)
	if requireEnabled && !cfg.Enabled {
		return LinuxDoCreditConfig{}, ErrLinuxDoCreditDisabled
	}
	if cfg.ClientID == "" || cfg.ClientSecret == "" || cfg.ExchangeRate <= 0 {
		return LinuxDoCreditConfig{}, ErrLinuxDoCreditBadConfig
	}
	return cfg, nil
}

func linuxDoCreditConfigFromSettings(settings *SystemSettings) LinuxDoCreditConfig {
	if settings == nil {
		return LinuxDoCreditConfig{}
	}
	exchangeRate := settings.LinuxDoCreditExchangeRate
	if exchangeRate <= 0 {
		exchangeRate = 1
	}
	return LinuxDoCreditConfig{
		Enabled:      settings.LinuxDoCreditEnabled,
		ClientID:     strings.TrimSpace(settings.LinuxDoCreditClientID),
		ClientSecret: strings.TrimSpace(settings.LinuxDoCreditClientSecret),
		ExchangeRate: exchangeRate,
	}
}

func buildLinuxDoCreditOrderTitle(siteName string) string {
	siteName = strings.TrimSpace(siteName)
	if siteName == "" {
		siteName = "Sub2API"
	}
	title := siteName + " Balance Top-up"
	if len([]rune(title)) > 64 {
		return string([]rune(title)[:64])
	}
	return title
}

func generateLinuxDoCreditOrderNo(userID int64) (string, error) {
	buf := make([]byte, 6)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return fmt.Sprintf("ldc_%d_%d_%s", userID, time.Now().Unix(), hex.EncodeToString(buf)), nil
}

func linuxDoCreditHistoryCode(outTradeNo string) string {
	sum := md5.Sum([]byte(outTradeNo))
	return strings.ToUpper(hex.EncodeToString(sum[:]))
}

func signLinuxDoCreditFields(fields map[string]string, secret string) string {
	keys := make([]string, 0, len(fields))
	for key, value := range fields {
		if key == "sign" || key == "sign_type" {
			continue
		}
		if strings.TrimSpace(value) == "" {
			continue
		}
		keys = append(keys, key)
	}
	sort.Strings(keys)

	var builder strings.Builder
	for i, key := range keys {
		if i > 0 {
			builder.WriteByte('&')
		}
		builder.WriteString(key)
		builder.WriteByte('=')
		builder.WriteString(fields[key])
	}
	builder.WriteString(secret)
	sum := md5.Sum([]byte(builder.String()))
	return hex.EncodeToString(sum[:])
}

func normalizeLinuxDoCreditAmount(amount float64) float64 {
	if amount <= 0 {
		return 0
	}
	return math.Round(amount*100) / 100
}

func normalizeLinuxDoCreditCredit(amount float64) float64 {
	if amount <= 0 {
		return 0
	}
	return math.Round(amount*1e8) / 1e8
}

func formatLinuxDoCreditAmount(amount float64) string {
	return strconv.FormatFloat(normalizeLinuxDoCreditAmount(amount), 'f', 2, 64)
}

func parseLinuxDoCreditAmount(raw string) (float64, error) {
	value, err := strconv.ParseFloat(strings.TrimSpace(raw), 64)
	if err != nil {
		return 0, err
	}
	return normalizeLinuxDoCreditAmount(value), nil
}

func sameLinuxDoCreditAmount(left, right float64) bool {
	return math.Abs(normalizeLinuxDoCreditAmount(left)-normalizeLinuxDoCreditAmount(right)) < 0.000001
}

func flattenLinuxDoCreditValues(values url.Values) map[string]string {
	flat := make(map[string]string, len(values))
	for key, items := range values {
		if len(items) == 0 {
			continue
		}
		flat[key] = strings.TrimSpace(items[0])
	}
	return flat
}

func marshalLinuxDoCreditValues(values map[string]string) (string, error) {
	body, err := json.Marshal(values)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func readHTTPBody(body io.Reader) ([]byte, error) {
	return io.ReadAll(body)
}
