package http

import "github.com/gin-gonic/gin"

func SetupRouter(handler *OrderHandler) *gin.Engine {
	r := gin.Default()

	r.GET("/orders", handler.ListOrders)
	r.POST("/orders", handler.CreateOrder)             // [cite: 75]
	r.PATCH("/orders/:id/cancel", handler.CancelOrder) // [cite: 85]

	return r
}
