package middleware

import (
	"github.com/sno6/gosane/internal/http"
	"github.com/sno6/gosane/internal/sentry"

	"github.com/gin-gonic/gin"
)

func Recovery(sentryClient *sentry.Sentry) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			uID, _ := UserUUIDFromContext(c)

			if err := recover(); err != nil {
				httpErr := http.Internal
				httpErr.Meta = err

				sentryClient.CaptureError(c, httpErr, &sentry.Meta{
					UserID:  uID,
					Request: c.Request,
					IsPanic: true,
				})

				c.AbortWithStatus(httpErr.Status)
			}
		}()

		c.Next()
	}
}
