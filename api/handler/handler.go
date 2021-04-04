package handler

import (
	"github.com/gin-gonic/gin"
)

// A Handler defines anything that can handle an HTTP(S) request.
type Handler interface {
	Path() string
	Method() string
	HandleFunc(c *gin.Context)
}

// A HandlerGroup defines a group of handlers that are logically grouped.
type HandlerGroup interface {
	RelativePath() string
	Handlers() []Handler
}

// MiddlewareChainer can be implemented by single route handlers or handler groups.
// If a handler group implements this interface, each handler will perform the actions
// of the middleware chain in array order.
type MiddlewareChainer interface {
	MiddlewareChain() gin.HandlersChain
}
