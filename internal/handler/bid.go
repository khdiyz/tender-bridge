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

// @Description Submit Bid
// @Summary Submit Bid
// @Tags Bid
// @Accept json
// @Produce json
// @Param id path string true "tender id"
// @Param create body models.CreateBid true "Submit bid"
// @Success 201 {object} submitBidResponse
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

	bidId, err := h.service.Bid.SubmitBid(body)
	if err != nil {
		fromError(c, err)
		return
	}

	c.JSON(http.StatusCreated, submitBidResponse{
		Id:    bidId,
		Price: body.Price,
	})
}

// @Description Get Contractor Bids
// @Summary Get Contractor Bids
// @Tags Bid
// @Accept json
// @Produce json
// @Success 200 {object} []models.Bid
// @Failure 400,401,404,500 {object} ErrorResponse
// @Router /api/contractor/bids [get]
// @Security ApiKeyAuth
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

// @Description Get Client Tender Bids
// @Summary Get Client Tender Bids
// @Tags Bid
// @Accept json
// @Produce json
// @Param id path string true "tender id"
// @Success 200 {object} []models.Bid
// @Failure 400,401,404,500 {object} ErrorResponse
// @Router /api/client/tenders/{id}/bids [get]
// @Security ApiKeyAuth
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

// @Description Award Bid
// @Summary Award Bid
// @Tags Tender
// @Accept json
// @Produce json
// @Param id path string true "tender id"
// @Param bidId path string true "tender id"
// @Success 200 {object} BaseResponse
// @Failure 400,401,404,500 {object} ErrorResponse
// @Router /api/client/tenders/{id}/award/{bidId} [post]
// @Security ApiKeyAuth
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

	c.JSON(http.StatusOK, BaseResponse{
		Message: "Bid awarded successfully",
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

// @Description Get User Bids
// @Summary Get User Bids
// @Tags Bid
// @Accept json
// @Produce json
// @Param id path string true "user id"
// @Param page query string true "page" Default(1)
// @Param limit query string true "limit" Default(10)
// @Success 200 {object} []models.Bid
// @Failure 400,401,404,500 {object} ErrorResponse
// @Router /api/users/{id}/bids [get]
// @Security ApiKeyAuth
func (h *Handler) getUserBids(c *gin.Context) {
	pagination, err := listPagination(c)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, err)
		return
	}

	userId, err := getUUIDParam(c, "id")
	if err != nil {
		errorResponse(c, http.StatusNotFound, errors.New("error: User not found"))
		return
	}

	userInfo, err := getUserInfo(c)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, err)
		return
	}
	if userInfo.Id != userId {
		errorResponse(c, http.StatusForbidden, errors.New("access denied"))
		return
	}

	var bidFilter models.BidFilter
	bidFilter.Limit = pagination.Limit
	bidFilter.Offset = pagination.Offset
	bidFilter.ContractorId = userId

	bids, _, err := h.service.Bid.GetBids(bidFilter)
	if err != nil {
		fromError(c, err)
		return
	}

	c.JSON(http.StatusOK, bids)
}
