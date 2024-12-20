package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"tender-bridge/config"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

const (
	AuthorizationHeader = "Authorization"
	UserCtx             = "user_id"
	RoleCtx             = "role"
)

func (h *Handler) userIdentity(c *gin.Context) {
	header := c.GetHeader(AuthorizationHeader)
	if header == "" {
		errorResponse(c, http.StatusUnauthorized, errors.New("error: Missing token"))
		c.Abort()
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		errorResponse(c, http.StatusUnauthorized, errors.New("invalid auth header"))
		c.Abort()
		return
	}

	if len(headerParts[1]) == 0 {
		errorResponse(c, http.StatusUnauthorized, errors.New("token is empty"))
		c.Abort()
		return
	}

	claims, err := h.service.Authorization.ParseToken(headerParts[1])
	if err != nil {
		errorResponse(c, http.StatusUnauthorized, err)
		c.Abort()
		return
	}

	if claims.Type != config.TokenTypeAccess {
		errorResponse(c, http.StatusUnauthorized, errors.New("invalid token type"))
		c.Abort()
		return
	}

	c.Set(UserCtx, claims.UserId)
	c.Set(RoleCtx, claims.Role)
	c.Next()
}

func corsMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		ctx.Writer.Header().Set("Content-Type", "application/json")
		ctx.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		ctx.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With,Access-Control-Request-Method, Access-Control-Request-Headers")
		ctx.Header("Access-Control-Max-Age", "3600")
		ctx.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH, HEAD")
		if ctx.Request.Method == "OPTIONS" {
			ctx.AbortWithStatus(204)
			return
		}
		ctx.Next()
	}
}

var redisClient *redis.Client

func init() {
	cfg := config.GetConfig()

	redisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.RedisHost, cfg.RedisPort),
		Password: cfg.RedisPassword,
		DB:       0,
	})
}

// RateLimitMiddleware enforces rate limits for a specific endpoint
func rateLimitMiddleware(limit int, duration time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		userInfo, err := getUserInfo(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		if userInfo.Id.String() == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		ctx := context.Background()
		key := "rate_limit:" + userInfo.Id.String()

		// Increment the count and set expiration if key doesn't exist
		count, err := redisClient.Incr(ctx, key).Result()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to enforce rate limit"})
			c.Abort()
			return
		}

		// Set expiration on the first request
		if count == 1 {
			redisClient.Expire(ctx, key, duration)
		}

		// Check if the limit is exceeded
		if int(count) > limit {
			ttl, _ := redisClient.TTL(ctx, key).Result()
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":       "Rate limit exceeded",
				"retry_after": ttl.Seconds(),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
