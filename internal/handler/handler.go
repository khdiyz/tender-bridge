package handler

import (
	"tender-bridge/config"
	"tender-bridge/docs"
	"tender-bridge/internal/service"
	"tender-bridge/internal/ws"
	"tender-bridge/pkg/logger"
	"time"

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

	// Setup Swagger documentation
	h.setupSwagger(router)

	// Apply middleware
	router.Use(corsMiddleware())

	// Public routes
	h.setupPublicRoutes(router)

	// Protected API routes
	api := router.Group("/api", h.userIdentity)
	h.setupClientRoutes(api)
	h.setupContractorRoutes(api)

	// WebSocket route
	router.GET("/ws", func(c *gin.Context) {
		ws.HandleWebSocket(c.Writer, c.Request)
	})

	ws.StartWebSocketHub()

	return router
}

func (h *Handler) setupSwagger(router *gin.Engine) {
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler), func(ctx *gin.Context) {
		docs.SwaggerInfo.Host = ctx.Request.Host
		if ctx.Request.TLS != nil {
			docs.SwaggerInfo.Schemes = []string{"https"}
		}
	})
}

func (h *Handler) setupPublicRoutes(router *gin.Engine) {
	router.POST("/register", h.register)
	router.POST("/login", h.login)
}

func (h *Handler) setupClientRoutes(api *gin.RouterGroup) {
	clientTenders := api.Group("/client/tenders")
	{
		clientTenders.POST("", h.createTender)
		clientTenders.GET("", h.getTenders)
		clientTenders.GET("/:id", h.getTender)
		clientTenders.PUT("/:id", h.updateTenderStatus)
		clientTenders.DELETE("/:id", h.deleteTender)
		clientTenders.GET("/:id/bids", h.getClientTenderBids)
		clientTenders.POST("/:id/award/:bidId", h.awardBid)
	}

	users := api.Group("/users")
	{
		users.GET("/:id/tenders", h.getUserTenders)
		users.GET("/:id/bids", h.getUserBids)
	}
}

func (h *Handler) setupContractorRoutes(api *gin.RouterGroup) {
	contractorBids := api.Group("/contractor/tenders/:id/bid")
	{
		contractorBids.POST("", rateLimitMiddleware(5, time.Minute), h.submitBid)
	}

	api.GET("/contractor/bids", h.getContractorBids)
	api.DELETE("/contractor/bids/:id", h.deleteContractorBid)
}
