package handler

import (
	"github.com/B022MC/b022hub/internal/service"
	"github.com/gin-gonic/gin"
)

func bindSubscriptionRequestProxy(c *gin.Context, router *service.SubscriptionRequestProxyRouter, subscription *service.UserSubscription) {
	if c == nil || c.Request == nil || router == nil || subscription == nil {
		return
	}
	c.Request = c.Request.WithContext(router.BindRandomProxy(c.Request.Context()))
}
