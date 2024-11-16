package handler

import (
	"net/http"
	"tender-bridge/internal/models"
	"tender-bridge/pkg/validator"

	"github.com/gin-gonic/gin"
)

type authResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// @Description Register User
// @Summary Register User
// @Tags Auth
// @Accept json
// @Produce json
// @Param signup body models.Register true "Register"
// @Success 200 {object} authResponse
// @Failure 400,404,500 {object} ErrorResponse
// @Router /api/register [post]
func (h *Handler) register(c *gin.Context) {
	var body models.Register

	if err := c.ShouldBindJSON(&body); err != nil {
		errorResponse(c, http.StatusBadRequest, err)
		return
	}

	if err := validator.ValidatePayloads(body); err != nil {
		errorResponse(c, http.StatusBadRequest, err)
		return
	}

	accessToken, refreshToken, err := h.service.Authorization.Register(body)
	if err != nil {
		fromError(c, err)
		return
	}

	c.JSON(http.StatusOK, authResponse{
		AccessToken:  accessToken.Token,
		RefreshToken: refreshToken.Token,
	})
}

// Login
// @Description Login User
// @Summary Login User
// @Tags Auth
// @Accept json
// @Produce json
// @Param login body models.Login true "Login"
// @Success 200 {object} authResponse
// @Failure 400,404,500 {object} ErrorResponse
// @Router /api/login [post]
func (h *Handler) login(c *gin.Context) {
	var body models.Login

	if err := c.ShouldBindJSON(&body); err != nil {
		errorResponse(c, http.StatusBadRequest, err)
		return
	}

	if err := validator.ValidatePayloads(body); err != nil {
		errorResponse(c, http.StatusBadRequest, err)
		return
	}

	accessToken, refreshToken, err := h.service.Authorization.Login(body)
	if err != nil {
		fromError(c, err)
		return
	}

	c.JSON(http.StatusOK, authResponse{
		AccessToken:  accessToken.Token,
		RefreshToken: refreshToken.Token,
	})
}
