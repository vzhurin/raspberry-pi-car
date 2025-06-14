package http

import "github.com/gin-gonic/gin"

func registerHandler(router *gin.Engine, handler *Handler, middlewares ...gin.HandlerFunc) {
	endpoints := router.Group("", middlewares...)

	api := endpoints.Group("/api")
	api.POST("/move", handler.Move)

	static := endpoints.Group("/static")
	static.Static("/control", "./static/control")
}
