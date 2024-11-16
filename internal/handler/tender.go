package handler

import (
	"errors"
	"math"
	"net/http"
	"tender-bridge/config"
	"tender-bridge/internal/models"
	"tender-bridge/pkg/validator"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type createResponse struct {
	Id uuid.UUID `json:"id"`
}

// @Description Create Tender
// @Summary Create Tender
// @Tags Tender
// @Accept json
// @Produce json
// @Param create body models.CreateTender true "Create tender"
// @Success 200 {object} createResponse
// @Failure 400,404,500 {object} ErrorResponse
// @Router /api/tenders [post]
// @Security ApiKeyAuth
func (h *Handler) createTender(c *gin.Context) {
	var body models.CreateTender

	userInfo, err := getUserInfo(c)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, err)
		return
	}

	if userInfo.Role != config.RoleClient {
		errorResponse(c, http.StatusForbidden, errors.New("permission denied"))
		return
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		errorResponse(c, http.StatusBadRequest, err)
		return
	}
	body.ClientId = userInfo.Id

	if err := validator.ValidatePayloads(body); err != nil {
		errorResponse(c, http.StatusBadRequest, err)
		return
	}

	tenderId, err := h.service.Tender.CreateTender(body)
	if err != nil {
		fromError(c, err)
		return
	}

	c.JSON(http.StatusOK, createResponse{
		Id: tenderId,
	})
}

type getTendersResponse struct {
	Tenders    []models.Tender   `json:"data"`
	Pagination models.Pagination `json:"pagination"`
}

// @Description Get Tenders
// @Summary Get Tenders
// @Tags Tender
// @Accept json
// @Produce json
// @Param limit query int64 true "limit" default(10)
// @Param page  query int64 true "page" default(1)
// @Param search  query string false "search"
// @Success 200 {object} getTendersResponse
// @Failure 400,401,404,500 {object} ErrorResponse
// @Router /api/tenders [get]
// @Security ApiKeyAuth
func (h *Handler) getTenders(c *gin.Context) {
	pagination, err := listPagination(c)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, err)
		return
	}

	var filter models.TenderFilter
	filter.Limit = pagination.Limit
	filter.Offset = pagination.Offset

	search := c.Query("search")
	if search != "" {
		filter.Search = search
	}

	tenders, total, err := h.service.Tender.GetTenders(filter)
	if err != nil {
		fromError(c, err)
		return
	}

	pagination.TotalCount = total
	pageCount := math.Ceil(float64(pagination.TotalCount) / float64(pagination.Limit))
	pagination.PageCount = int(pageCount)

	c.JSON(http.StatusOK, getTendersResponse{
		Tenders:    tenders,
		Pagination: pagination,
	})
}

// @Description Get Tender
// @Summary Get Tender
// @Tags Tender
// @Accept json
// @Produce json
// @Param id path string true "tender id"
// @Success 200 {object} models.Tender
// @Failure 400,401,404,500 {object} ErrorResponse
// @Router /api/tenders/{id} [get]
// @Security ApiKeyAuth
func (h *Handler) getTender(c *gin.Context) {
	id, err := getUUIDParam(c, "id")
	if err != nil {
		errorResponse(c, http.StatusBadRequest, err)
		return
	}

	tender, err := h.service.Tender.GetTender(id)
	if err != nil {
		fromError(c, err)
		return
	}

	c.JSON(http.StatusOK, tender)
}

// @Description Update Tender
// @Summary Update Tender
// @Tags Tender
// @Accept json
// @Produce json
// @Param id path string true "tender id"
// @Param update body models.UpdateTender true "update tender"
// @Success 200 {object} createResponse
// @Failure 400,401,404,500 {object} ErrorResponse
// @Router /api/tenders/{id} [put]
// @Security ApiKeyAuth
func (h *Handler) updateTender(c *gin.Context) {
	id, err := getUUIDParam(c, "id")
	if err != nil {
		errorResponse(c, http.StatusBadRequest, err)
		return
	}

	userInfo, err := getUserInfo(c)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, err)
		return
	}

	if userInfo.Role != config.RoleClient {
		errorResponse(c, http.StatusForbidden, errors.New("permission denied"))
		return
	}

	var body models.UpdateTender
	if err = c.ShouldBindJSON(&body); err != nil {
		errorResponse(c, http.StatusBadRequest, err)
		return
	}
	body.Id = id
	body.ClientId = userInfo.Id

	if err := validator.ValidatePayloads(body); err != nil {
		errorResponse(c, http.StatusBadRequest, err)
		return
	}

	if err = h.service.Tender.UpdateTender(body); err != nil {
		fromError(c, err)
		return
	}

	c.JSON(http.StatusOK, createResponse{
		Id: id,
	})
}

// @Description Delete Tender
// @Summary Delete Tender
// @Tags Tender
// @Accept json
// @Produce json
// @Param id path string true "tender id"
// @Success 200 {object} createResponse
// @Failure 400,401,404,500 {object} ErrorResponse
// @Router /api/tenders/{id} [delete]
// @Security ApiKeyAuth
func (h *Handler) deleteTender(c *gin.Context) {
	id, err := getUUIDParam(c, "id")
	if err != nil {
		errorResponse(c, http.StatusBadRequest, err)
		return
	}

	userInfo, err := getUserInfo(c)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, err)
		return
	}

	if userInfo.Role != config.RoleClient {
		errorResponse(c, http.StatusForbidden, errors.New("permission denied"))
		return
	}

	if err = h.service.Tender.DeleteTender(id); err != nil {
		fromError(c, err)
		return
	}

	c.JSON(http.StatusOK, createResponse{
		Id: id,
	})
}
