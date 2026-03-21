package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/B022MC/b022hub/internal/pkg/response"
	middleware2 "github.com/B022MC/b022hub/internal/server/middleware"
	"github.com/B022MC/b022hub/internal/service"

	"github.com/gin-gonic/gin"
)

type LinuxDoCreditHandler struct {
	linuxDoCreditService *service.LinuxDoCreditService
}

func NewLinuxDoCreditHandler(linuxDoCreditService *service.LinuxDoCreditService) *LinuxDoCreditHandler {
	return &LinuxDoCreditHandler{linuxDoCreditService: linuxDoCreditService}
}

type createLinuxDoCreditOrderRequest struct {
	Amount float64 `json:"amount" binding:"required,gt=0"`
}

func (h *LinuxDoCreditHandler) CreateOrder(c *gin.Context) {
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	var req createLinuxDoCreditOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	checkout, err := h.linuxDoCreditService.CreateCheckout(c.Request.Context(), subject.UserID, req.Amount)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	response.Success(c, gin.H{
		"order": gin.H{
			"id":              checkout.Order.ID,
			"provider":        checkout.Order.Provider,
			"out_trade_no":    checkout.Order.OutTradeNo,
			"title":           checkout.Order.Title,
			"amount":          checkout.Order.Amount,
			"credited_amount": checkout.Order.CreditedAmount,
			"status":          checkout.Order.Status,
			"created_at":      checkout.Order.CreatedAt,
		},
		"checkout_url": checkout.CheckoutURL,
		"fields":       checkout.Fields,
	})
}

func (h *LinuxDoCreditHandler) ListOrders(c *gin.Context) {
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	limit := 10
	if raw := strings.TrimSpace(c.Query("limit")); raw != "" {
		if parsed, err := strconv.Atoi(raw); err == nil {
			limit = parsed
		}
	}

	orders, err := h.linuxDoCreditService.ListOrders(c.Request.Context(), subject.UserID, limit)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	response.Success(c, gin.H{"items": orders})
}

func (h *LinuxDoCreditHandler) GetOrder(c *gin.Context) {
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	syncRemote := c.Query("sync") == "1" || strings.EqualFold(c.Query("sync"), "true")
	order, err := h.linuxDoCreditService.GetOrder(c.Request.Context(), subject.UserID, c.Param("out_trade_no"), syncRemote)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	response.Success(c, order)
}

func (h *LinuxDoCreditHandler) Notify(c *gin.Context) {
	if err := h.linuxDoCreditService.HandleNotify(c.Request.Context(), c.Request.URL.Query()); err != nil {
		c.String(http.StatusBadRequest, "fail")
		return
	}
	c.String(http.StatusOK, "success")
}
