package handler

import (
	"tender-bridge/config"
	"tender-bridge/docs"
	"tender-bridge/internal/service"
	"tender-bridge/pkg/logger"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handler struct {
	service *service.Service
	logger  *logger.Logger
}

func NewHandler(service *service.Service, loggers *logger.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  loggers,
	}
}

func (h *Handler) InitRoutes(cfg *config.Config) *gin.Engine {
	router := gin.Default()

	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler), func(ctx *gin.Context) {
		docs.SwaggerInfo.Host = ctx.Request.Host
		if ctx.Request.TLS != nil {
			docs.SwaggerInfo.Schemes = []string{"https"}
		}
	})

	router.Use(corsMiddleware())

	// auth handlers
	router.POST("/register", h.register)
	router.POST("/login", h.login)

	api := router.Group("/api", h.userIdentity)
	{
		tenders := api.Group("/client/tenders")
		{
			tenders.POST("", h.createTender)
			tenders.GET("", h.getTenders)
			tenders.GET("/:id", h.getTender)
			tenders.PUT("/:id", h.updateTenderStatus)
			tenders.DELETE("/:id", h.deleteTender)
			tenders.GET("/:id/bids", h.getClientTenderBids)
			tenders.POST(":id/award/:bidId", h.awardBid)
		}

		bids := api.Group("/contractor/tenders/:id/bid")
		{
			bids.POST("", h.submitBid)
		}

		api.GET("/contractor/bids", h.getContractorBids)
		api.DELETE("/contractor/bids/:id", h.deleteContractorBid)
	}

	return router
}
