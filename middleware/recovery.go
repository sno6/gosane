package middleware

import (
	"log"

	"github.com/sno6/gosane/internal/http"
	"github.com/sno6/gosane/internal/sentry"

	"github.com/gin-gonic/gin"
)

func Recovery(sentryClient *sentry.Sentry, logger *log.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			uID, _ := UserUUIDFromContext(c)

			if err := recover(); err != nil {
				log.Println(err)

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
