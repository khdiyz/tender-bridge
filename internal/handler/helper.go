package handler

import (
	"errors"
	"fmt"
	"strconv"
	"tender-bridge/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var (
	errInvalidUserId   = errors.New("invalid user id")
	errInvalidUserRole = errors.New("invalid user role")
)

type UserInfo struct {
	Id   uuid.UUID
	Role string
}

func getUserInfo(ctx *gin.Context) (UserInfo, error) {
	var userInfo UserInfo

	userId, ok := ctx.Get(UserCtx)
	if !ok {
		return UserInfo{}, errInvalidUserId
	}

	userInfo.Id, ok = userId.(uuid.UUID)
	if !ok {
		return UserInfo{}, errInvalidUserId
	}

	role, ok := ctx.Get(RoleCtx)
	if !ok {
		return UserInfo{}, errInvalidUserRole
	}

	userInfo.Role, ok = role.(string)
	if !ok {
		return UserInfo{}, errInvalidUserRole
	}

	return userInfo, nil
}

func listPagination(c *gin.Context) (models.Pagination, error) {
	page, err := getPageQuery(c)
	if err != nil {
		return models.Pagination{}, err
	}

	limit, err := getLimitQuery(c)
	if err != nil {
		return models.Pagination{}, err
	}

	offset, limit := calculatePagination(page, limit)

	return models.Pagination{
		Page:   page,
		Limit:  limit,
		Offset: offset,
	}, nil
}

func getPageQuery(c *gin.Context) (int, error) {
	pageStr := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		return 0, errors.New("invalid page parameter")
	}
	return page, nil
}

func getLimitQuery(c *gin.Context) (int, error) {
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		return 0, errors.New("invalid limit parameter")
	}
	return limit, nil
}

func calculatePagination(page, limit int) (int, int) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit
	return offset, limit
}

func getUUIDParam(c *gin.Context, param string) (uuid.UUID, error) {
	paramValue := c.Param(param)
	if paramValue != "" {
		id, err := uuid.Parse(paramValue)
		if err != nil {
			return uuid.Nil, fmt.Errorf("invalid param %v", paramValue)
		}
		return id, nil
	}
	return uuid.Nil, errors.New("empty param value")
}
