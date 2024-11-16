package handler

import (
	"errors"
	"net/http"
	"tender-bridge/internal/models"
	"tender-bridge/pkg/validator"

	"github.com/gin-gonic/gin"
)

type authResponse struct {
	Token string `json:"token"`
}

// @Description Register User
// @Summary Register User
// @Tags Auth
// @Accept json
// @Produce json
// @Param signup body models.Register true "Register"
// @Success 200 {object} authResponse
// @Failure 400,404,500 {object} ErrorResponse
// @Router /register [post]
func (h *Handler) register(c *gin.Context) {
	var body models.Register

	if err := c.ShouldBindJSON(&body); err != nil {
		errorResponse(c, http.StatusBadRequest, err)
		return
	}

	if body.Username == "" && body.Email == "" {
		errorResponse(c, http.StatusBadRequest, errors.New("username or email cannot be empty"))
		return
	}

	if err := validator.ValidatePayloads(body); err != nil {
		errorResponse(c, http.StatusBadRequest, err)
		return
	}

	accessToken, _, err := h.service.Authorization.Register(body)
	if err != nil {
		fromError(c, err)
		return
	}

	c.JSON(http.StatusCreated, authResponse{
		Token: accessToken.Token,
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
// @Router /login [post]
func (h *Handler) login(c *gin.Context) {
	var body models.Login

	if err := c.ShouldBindJSON(&body); err != nil {
		errorResponse(c, http.StatusBadRequest, err)
		return
	}

	if body.Username == "" || body.Password == "" {
		errorResponse(c, http.StatusBadRequest, errors.New("error: Username and password are required"))
		return
	}

	if err := validator.ValidatePayloads(body); err != nil {
		errorResponse(c, http.StatusBadRequest, err)
		return
	}

	accessToken, _, err := h.service.Authorization.Login(body)
	if err != nil {
		fromError(c, err)
		return
	}

	c.JSON(http.StatusOK, authResponse{
		Token: accessToken.Token,
	})
}
