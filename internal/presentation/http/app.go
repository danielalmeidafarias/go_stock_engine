package http

import (
	"log"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type GinApp struct {
	gin *gin.Engine
}

func (g GinApp) Run() {
	if err := g.gin.Run(":8080"); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}

func NewGinApp(handler *ProductStockHandler) GinApp {
	r := gin.Default()

	stock := r.Group("/stock")
	{
		stock.POST("", handler.Create)
		stock.GET("", handler.GetAll)
		stock.GET("/:id", handler.GetOne)
		stock.PUT("/:id", handler.Update)
		stock.DELETE("/:id", handler.Delete)
		stock.GET("/category/:category", handler.GetByCategory)
	}

	r.GET("/restock/priorities", handler.GetRestockPriorities)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return GinApp{
		gin: r,
	}
}
