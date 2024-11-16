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

type submitBidResponse struct {
	Id    uuid.UUID `json:"id"`
	Price int64     `json:"price"`
}

// @Description Create Bid
// @Summary Create Bid
// @Tags Bid
// @Accept json
// @Produce json
// @Param id path string true "tender id"
// @Param create body models.CreateBid true "Create bid"
// @Success 200 {object} createResponse
// @Failure 400,401,404,500 {object} ErrorResponse
// @Router /api/contractor/tenders/{id}/bid [post]
// @Security ApiKeyAuth
func (h *Handler) submitBid(c *gin.Context) {
	userInfo, err := getUserInfo(c)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, err)
		return
	}

	if userInfo.Role != config.RoleContractor {
		errorResponse(c, http.StatusForbidden, errors.New("permission denied"))
		return
	}

	tenderId, err := getUUIDParam(c, "id")
	if err != nil {
		errorResponse(c, http.StatusNotFound, errors.New("error: Tender not found"))
		return
	}

	var body models.CreateBid
	if err := c.ShouldBindJSON(&body); err != nil {
		errorResponse(c, http.StatusBadRequest, err)
		return
	}
	body.TenderId = tenderId
	body.ContractorId = userInfo.Id

	if err := validator.ValidatePayloads(body); err != nil {
		errorResponse(c, http.StatusBadRequest, err)
		return
	}

	bidId, err := h.service.Bid.CreateBid(body)
	if err != nil {
		fromError(c, err)
		return
	}

	c.JSON(http.StatusCreated, submitBidResponse{
		Id:    bidId,
		Price: body.Price,
	})
}

func (h *Handler) getContractorBids(c *gin.Context) {
	userInfo, err := getUserInfo(c)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, err)
		return
	}

	if userInfo.Role != config.RoleContractor {
		errorResponse(c, http.StatusForbidden, errors.New("permission denied"))
		return
	}

	pagination, err := listPagination(c)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, err)
		return
	}

	var filter models.BidFilter
	filter.Limit = pagination.Limit
	filter.Offset = pagination.Offset
	filter.ContractorId = userInfo.Id

	bids, _, err := h.service.Bid.GetBids(filter)
	if err != nil {
		fromError(c, err)
		return
	}

	c.JSON(http.StatusOK, bids)
}

func (h *Handler) getClientTenderBids(c *gin.Context) {
	userInfo, err := getUserInfo(c)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, err)
		return
	}

	if userInfo.Role != config.RoleClient {
		errorResponse(c, http.StatusForbidden, errors.New("permission denied"))
		return
	}

	tenderId, err := getUUIDParam(c, "id")
	if err != nil {
		errorResponse(c, http.StatusNotFound, errors.New("error: Tender not found or access denied"))
		return
	}

	pagination, err := listPagination(c)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, err)
		return
	}

	var filter models.BidFilter
	filter.Limit = pagination.Limit
	filter.Offset = pagination.Offset
	filter.TenderId = tenderId

	bids, _, err := h.service.Bid.GetBids(filter)
	if err != nil {
		fromError(c, err)
		return
	}

	c.JSON(http.StatusOK, bids)
}

func (h *Handler) awardBid(c *gin.Context) {
	userInfo, err := getUserInfo(c)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, err)
		return
	}

	if userInfo.Role != config.RoleClient {
		errorResponse(c, http.StatusForbidden, errors.New("permission denied"))
		return
	}

	tenderId, err := getUUIDParam(c, "id")
	if err != nil {
		errorResponse(c, http.StatusNotFound, errors.New("error: Tender not found or access denied"))
		return
	}

	bidId, err := getUUIDParam(c, "bidId")
	if err != nil {
		errorResponse(c, http.StatusNotFound, errors.New("error: Bid not found"))
		return
	}

	err = h.service.Bid.AwardBid(userInfo.Id, tenderId, bidId)
	if err != nil {
		fromError(c, err)
		return
	}

	c.JSON(http.StatusOK, ErrorResponse{
		ErrorMessage: "Bid awarded successfully",
	})
}

func (h *Handler) deleteContractorBid(c *gin.Context) {
	userInfo, err := getUserInfo(c)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, err)
		return
	}

	if userInfo.Role != config.RoleContractor {
		errorResponse(c, http.StatusForbidden, errors.New("permission denied"))
		return
	}

	bidId, err := getUUIDParam(c, "id")
	if err != nil {
		errorResponse(c, http.StatusNotFound, errors.New("error: Bid not found or access denied"))
		return
	}

	if err = h.service.Bid.DeleteContractorBid(userInfo.Id, bidId); err != nil {
		fromError(c, err)
		return
	}

	c.JSON(http.StatusOK, BaseResponse{
		Message: "Bid deleted successfully",
	})
}
