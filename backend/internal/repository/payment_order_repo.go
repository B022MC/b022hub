package repository

import (
	"context"
	"strings"
	"time"

	dbent "github.com/B022MC/b022hub/ent"
	"github.com/B022MC/b022hub/ent/paymentorder"
	"github.com/B022MC/b022hub/internal/service"
)

type paymentOrderRepository struct {
	client *dbent.Client
}

func NewPaymentOrderRepository(client *dbent.Client) service.PaymentOrderRepository {
	return &paymentOrderRepository{client: client}
}

func (r *paymentOrderRepository) Create(ctx context.Context, order *service.PaymentOrder) error {
	created, err := clientFromContext(ctx, r.client).PaymentOrder.Create().
		SetProvider(order.Provider).
		SetOutTradeNo(order.OutTradeNo).
		SetUserID(order.UserID).
		SetTitle(order.Title).
		SetAmount(order.Amount).
		SetCreditedAmount(order.CreditedAmount).
		SetStatus(order.Status).
		SetNillableProviderTradeNo(nullableString(order.ProviderTradeNo)).
		SetNillableRawProviderPayload(nullableString(order.RawProviderPayload)).
		SetNillablePaidAt(order.PaidAt).
		Save(ctx)
	if err != nil {
		return err
	}
	applyPaymentOrderEntity(order, created)
	return nil
}

func (r *paymentOrderRepository) GetByOutTradeNo(ctx context.Context, outTradeNo string) (*service.PaymentOrder, error) {
	entity, err := clientFromContext(ctx, r.client).PaymentOrder.Query().
		Where(paymentorder.OutTradeNoEQ(outTradeNo)).
		Only(ctx)
	if err != nil {
		return nil, translatePersistenceError(err, service.ErrPaymentOrderNotFound, nil)
	}
	return paymentOrderEntityToService(entity), nil
}

func (r *paymentOrderRepository) GetByUserAndOutTradeNo(ctx context.Context, userID int64, outTradeNo string) (*service.PaymentOrder, error) {
	entity, err := clientFromContext(ctx, r.client).PaymentOrder.Query().
		Where(paymentorder.UserIDEQ(userID), paymentorder.OutTradeNoEQ(outTradeNo)).
		Only(ctx)
	if err != nil {
		return nil, translatePersistenceError(err, service.ErrPaymentOrderNotFound, nil)
	}
	return paymentOrderEntityToService(entity), nil
}

func (r *paymentOrderRepository) ListByUser(ctx context.Context, userID int64, limit int) ([]service.PaymentOrder, error) {
	if limit <= 0 {
		limit = 10
	}

	entities, err := clientFromContext(ctx, r.client).PaymentOrder.Query().
		Where(paymentorder.UserIDEQ(userID)).
		Order(dbent.Desc(paymentorder.FieldCreatedAt)).
		Limit(limit).
		All(ctx)
	if err != nil {
		return nil, err
	}

	orders := make([]service.PaymentOrder, 0, len(entities))
	for _, entity := range entities {
		orders = append(orders, *paymentOrderEntityToService(entity))
	}
	return orders, nil
}

func (r *paymentOrderRepository) MarkPaid(ctx context.Context, outTradeNo string, providerTradeNo string, rawProviderPayload string, paidAt time.Time) (*service.PaymentOrder, bool, error) {
	client := clientFromContext(ctx, r.client)

	affected, err := client.PaymentOrder.Update().
		Where(
			paymentorder.OutTradeNoEQ(outTradeNo),
			paymentorder.StatusEQ(service.PaymentOrderStatusPending),
		).
		SetStatus(service.PaymentOrderStatusPaid).
		SetNillableProviderTradeNo(nullableString(providerTradeNo)).
		SetNillableRawProviderPayload(nullableString(rawProviderPayload)).
		SetPaidAt(paidAt).
		Save(ctx)
	if err != nil {
		return nil, false, err
	}

	order, err := r.GetByOutTradeNo(ctx, outTradeNo)
	if err != nil {
		return nil, false, err
	}

	return order, affected > 0, nil
}

func paymentOrderEntityToService(entity *dbent.PaymentOrder) *service.PaymentOrder {
	if entity == nil {
		return nil
	}
	order := &service.PaymentOrder{}
	applyPaymentOrderEntity(order, entity)
	return order
}

func applyPaymentOrderEntity(dst *service.PaymentOrder, src *dbent.PaymentOrder) {
	dst.ID = src.ID
	dst.Provider = src.Provider
	dst.OutTradeNo = src.OutTradeNo
	dst.ProviderTradeNo = valueOrEmpty(src.ProviderTradeNo)
	dst.UserID = src.UserID
	dst.Title = src.Title
	dst.Amount = src.Amount
	dst.CreditedAmount = src.CreditedAmount
	dst.Status = src.Status
	dst.RawProviderPayload = valueOrEmpty(src.RawProviderPayload)
	dst.PaidAt = src.PaidAt
	dst.CreatedAt = src.CreatedAt
	dst.UpdatedAt = src.UpdatedAt
}

func nullableString(value string) *string {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}
	return &value
}

func valueOrEmpty(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}
