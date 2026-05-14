package handler

import (
	"net/http"
	"net/url"

	"github.com/Wei-Shaw/sub2api/internal/pkg/response"
	middleware2 "github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"

	"github.com/gin-gonic/gin"
)

type WechatNotificationHandler struct {
	wechatService *service.WechatOfficialNotificationService
}

func NewWechatNotificationHandler(wechatService *service.WechatOfficialNotificationService) *WechatNotificationHandler {
	return &WechatNotificationHandler{wechatService: wechatService}
}

func (h *WechatNotificationHandler) Status(c *gin.Context) {
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
	}
	if h.wechatService == nil {
		response.Success(c, service.WechatBindingStatus{})
		return
	}
	status, err := h.wechatService.Status(c.Request.Context(), subject.UserID)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, status)
}

func (h *WechatNotificationHandler) BindURL(c *gin.Context) {
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
	}
	if h.wechatService == nil {
		response.ErrorFrom(c, service.ErrWechatOfficialNotReady)
		return
	}
	bindURL, err := h.wechatService.BuildBindURL(c.Request.Context(), subject.UserID, c.Query("return_to"))
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, service.WechatBindURL{URL: bindURL})
}

func (h *WechatNotificationHandler) Callback(c *gin.Context) {
	if h.wechatService == nil {
		redirectWechatBindResult(c, "/profile", "error", "not_ready")
		return
	}
	code := c.Query("code")
	state := c.Query("state")
	returnTo, err := h.wechatService.CompleteBind(c.Request.Context(), code, state)
	if err != nil {
		redirectWechatBindResult(c, "/profile", "error", err.Error())
		return
	}
	redirectWechatBindResult(c, returnTo, "success", "")
}

func (h *WechatNotificationHandler) Unbind(c *gin.Context) {
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
	}
	if h.wechatService == nil {
		response.ErrorFrom(c, service.ErrWechatOfficialUnbound)
		return
	}
	if err := h.wechatService.Unbind(c.Request.Context(), subject.UserID); err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, gin.H{"message": "unbound"})
}

func redirectWechatBindResult(c *gin.Context, target, result, reason string) {
	if target == "" {
		target = "/profile"
	}
	u, err := url.Parse(target)
	if err != nil {
		target = "/profile"
		u, _ = url.Parse(target)
	}
	q := u.Query()
	q.Set("wechat_bind", result)
	if reason != "" {
		q.Set("reason", reason)
	}
	u.RawQuery = q.Encode()
	c.Redirect(http.StatusFound, u.String())
}
