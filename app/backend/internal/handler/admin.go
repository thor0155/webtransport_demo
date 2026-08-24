package handler

import (
	"api/internal/frameworks/obj"
	"api/internal/server/wts"
	"net/http"

	"github.com/cockroachdb/errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type adminController struct {
	logger  *zap.Logger
	manager *obj.Manager
}

func (c *adminController) Register(rg *gin.RouterGroup) {
	// health endpoints (K8s)
	rg.GET("/healthz", c.Healthz)
	rg.GET("/readyz", c.Readyz)
	rg.GET("/livez", c.Livez)

	if gin.IsDebugging() {
		rg.GET("/cert", c.getCertHash)
		rg.POST("/shutdown", c.shutdown)
	}
}

func (a *adminController) shutdown(c *gin.Context) {
	defer a.manager.Shutdown("close by api")
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (*adminController) Healthz(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (*adminController) Readyz(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ready",
	})
}

func (*adminController) Livez(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "alive"})
}

func (a *adminController) getCertHash(ctx *gin.Context) {
	servers := a.manager.FindServers(func(s obj.Server) bool {
		return s.Name() == "webtransport"
	})
	if len(servers) == 0 {
		ctx.AbortWithError(http.StatusInternalServerError, errors.New("lost webtransport server"))
		return
	}
	if wts, ok := servers[0].(*wts.Server); ok {
		ctx.String(http.StatusOK, wts.GetCertHash())
	} else {
		ctx.AbortWithError(http.StatusInternalServerError, errors.New("missing TLS provider"))
	}
}

func NewAdminController(logger *zap.Logger, m *obj.Manager) *adminController {
	return &adminController{
		logger:  logger,
		manager: m,
	}
}
