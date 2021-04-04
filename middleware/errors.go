package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/sno6/gosane/internal/http"
	"github.com/sno6/gosane/internal/sentry"
)

func Errors(sentryClient *sentry.Sentry) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Run all pre-handler middleware & handlers before handling errors.
		c.Next()

		err := c.Errors.Last()
		if err == nil {
			return
		}

		// Extract the application error context.
		appErr, ok := err.Meta.(error)

		httpErr := http.ErrorFromString(err.Err.Error())
		if httpErr != http.Internal {
			if appErr != nil && ok {
				// Only attach appErr context to the return object on non internal errors.
				httpErr.Meta = appErr.Error()
			}
		}

		uID, _ := UserUUIDFromContext(c)

		// Log the HTTP error with Sentry.
		sentryClient.CaptureError(c, httpErr, &sentry.Meta{
			UserID:  uID,
			Request: c.Request,
			IsPanic: httpErr == http.Internal,
		})

		c.JSON(httpErr.Status, httpErr)
	}
}
