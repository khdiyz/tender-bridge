package handler

import (
	"errors"
	"net/http"
	"strings"
	"tender-bridge/config"

	"github.com/gin-gonic/gin"
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
