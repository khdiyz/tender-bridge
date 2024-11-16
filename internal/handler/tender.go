package handler

import (
	"errors"
	"net/http"
	"tender-bridge/config"
	"tender-bridge/internal/models"
	"tender-bridge/pkg/validator"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	idQuery     = "id"
	searchQuery = "search"
)

type createTenderResponse struct {
	Id    uuid.UUID `json:"id"`
	Title string    `json:"title"`
}

// @Description Create Tender
// @Summary Create Tender
// @Tags Tender
// @Accept json
// @Produce json
// @Param create body models.CreateTender true "Create tender"
// @Success 201 {object} createTenderResponse
// @Failure 400,404,500 {object} ErrorResponse
// @Router /api/client/tenders [post]
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
		errorResponse(c, http.StatusBadRequest, errors.New("error: Invalid input"))
		return
	}

	tenderId, err := h.service.Tender.CreateTender(body)
	if err != nil {
		fromError(c, err)
		return
	}

	c.JSON(http.StatusCreated, createTenderResponse{
		Id:    tenderId,
		Title: body.Title,
	})
}

// @Description Get Tenders
// @Summary Get Tenders
// @Tags Tender
// @Accept json
// @Produce json
// @Param limit query int64 true "limit" default(10)
// @Param page  query int64 true "page" default(1)
// @Param search  query string false "search"
// @Success 200 {object} []models.Tender
// @Failure 400,401,404,500 {object} ErrorResponse
// @Router /api/client/tenders [get]
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

	search := c.Query(searchQuery)
	if search != "" {
		filter.Search = search
	}

	tenders, _, err := h.service.Tender.GetTenders(filter)
	if err != nil {
		fromError(c, err)
		return
	}

	c.JSON(http.StatusOK, tenders)
}

// @Description Get Tender
// @Summary Get Tender
// @Tags Tender
// @Accept json
// @Produce json
// @Param id path string true "tender id"
// @Success 200 {object} models.Tender
// @Failure 400,401,404,500 {object} ErrorResponse
// @Router /api/client/tenders/{id} [get]
// @Security ApiKeyAuth
func (h *Handler) getTender(c *gin.Context) {
	id, err := getUUIDParam(c, idQuery)
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

// @Description Update Tender Status
// @Summary Update Tender Status
// @Tags Tender
// @Accept json
// @Produce json
// @Param id path string true "tender id"
// @Param update body models.UpdateTenderStatus true "update tender status"
// @Success 200 {object} BaseResponse
// @Failure 400,401,404,500 {object} ErrorResponse
// @Router /api/client/tenders/{id} [put]
// @Security ApiKeyAuth
func (h *Handler) updateTenderStatus(c *gin.Context) {
	id, err := getUUIDParam(c, "id")
	if err != nil {
		errorResponse(c, http.StatusNotFound, errors.New("error: Tender not found"))
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

	var body models.UpdateTenderStatus
	if err = c.ShouldBindJSON(&body); err != nil {
		errorResponse(c, http.StatusBadRequest, err)
		return
	}
	body.Id = id

	if err = h.service.Tender.UpdateTenderStatus(body); err != nil {
		fromError(c, err)
		return
	}

	c.JSON(http.StatusOK, BaseResponse{
		Message: "Tender status updated",
	})
}

// @Description Delete Tender
// @Summary Delete Tender
// @Tags Tender
// @Accept json
// @Produce json
// @Param id path string true "tender id"
// @Success 200 {object} BaseResponse
// @Failure 400,401,404,500 {object} ErrorResponse
// @Router /api/client/tenders/{id} [delete]
// @Security ApiKeyAuth
func (h *Handler) deleteTender(c *gin.Context) {
	id, err := getUUIDParam(c, "id")
	if err != nil {
		errorResponse(c, http.StatusNotFound, errors.New("error: Tender not found or access denied"))
		return
	}

	userInfo, err := getUserInfo(c)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, err)
		return
	}

	if userInfo.Role != config.RoleClient {
		errorResponse(c, http.StatusForbidden, errors.New("error: Tender not found or access denied"))
		return
	}

	if err = h.service.Tender.DeleteTender(id); err != nil {
		fromError(c, err)
		return
	}

	c.JSON(http.StatusOK, BaseResponse{
		Message: "Tender deleted successfully",
	})
}
