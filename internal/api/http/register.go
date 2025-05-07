package http

import "github.com/gin-gonic/gin"

func registerHandler(router *gin.Engine, handler *Handler, middlewares ...gin.HandlerFunc) {
	endpoints := router.Group("/api", middlewares...)

	endpoints.POST("/move", handler.Move)
}
