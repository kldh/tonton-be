package monitor

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ldhk/tonton-be/pkg/telemetry/logging"
	"github.com/ldhk/tonton-be/pkg/telemetry/tracing"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func Gin(r gin.IRouter, serviceName string) {
	r.Use(otelgin.Middleware(serviceName), ginLog())
}

func ginLog() gin.HandlerFunc {
	ginPath := func(c *gin.Context) string {
		if c.FullPath() != "" {
			return c.FullPath()
		}

		return c.Request.URL.Path
	}

	return func(c *gin.Context) {
		ctx, start := c.Request.Context(), time.Now()

		ctx, l := logging.WithFields(ctx, map[string]interface{}{
			"trace_id":                 tracing.TraceID(ctx),
			"http.server.request_path": ginPath(c),
		})

		l.WithFields(map[string]interface{}{
			"http.server.method":    c.Request.Method,
			"http.server.client_ip": c.ClientIP(),
		}).Infof("server: received request")

		c.Request = c.Request.WithContext(ctx)
		c.Next()

		l = l.WithFields(map[string]interface{}{
			"http.server.status":  c.Writer.Status(),
			"http.server.latency": time.Since(start).String(),
		})

		if len(c.Errors) > 0 {
			l = l.WithField("error", c.Errors)
		}

		if c.Writer.Status() >= 500 {
			l.Errorf("http.server: sending response")
		} else {
			l.Infof("http.server: sending response")
		}
	}
}
