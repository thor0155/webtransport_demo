package https

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/requestid"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type RouterRegistry interface {
	Register(router *gin.RouterGroup)
}

type RouterRegistries = []RouterRegistry

type routerRegistryCollector []RouterRegistry

func ConvGinEngineToHandler(g *gin.Engine) http.Handler {
	return g
}

func NewGinRouter(log *zap.Logger, cfg *HttpServerConfig, routerRegistries RouterRegistries) *gin.Engine {

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(requestid.New())
	r.Use(ginzap.GinzapWithConfig(log, &ginzap.Config{
		TimeFormat: time.RFC3339,
		UTC:        false,
		SkipPaths: []string{
			"/healthz",
			"/readyz",
			"/livez",
		},
		Context: func(ctx *gin.Context) []zapcore.Field {
			var fields []zapcore.Field

			if len(ctx.Errors) > 0 {
				errors := make([]error, 0, len(ctx.Errors))
				for _, e := range ctx.Errors {
					errors = append(errors, e.Err)
				}
				fields = append(fields, zap.Errors("errors", errors))
			}
			return append(fields, zap.String("request_id", requestid.Get(ctx)))
		},
	}))
	r.Use(ginzap.RecoveryWithZap(log, true))

	if cfg.Cors.Enabled {

		corscfg := cors.DefaultConfig()
		corscfg.AllowCredentials = cfg.Cors.AllowCredentials
		corscfg.AllowOrigins = cfg.Cors.AllowOrigins
		// Wildcard(*) will not work with Access-Control-Allow-Originand need to have a list of origin/domains as mentioned above.
		corscfg.AllowAllOrigins = len(corscfg.AllowOrigins) == 0

		corscfg.AddAllowHeaders("Authorization")
		if len(cfg.Cors.ExposeHeaders) > 0 {
			corscfg.AddExposeHeaders(cfg.Cors.ExposeHeaders...)
		}

		r.Use(cors.New(corscfg))
	}

	root := r.Group("/")
	for _, router := range routerRegistries {
		router.Register(root)
	}

	return r
}
