package handler

import (
	"Auth/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) singUp(c *gin.Context) {

	var input models.SignUpInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Invalid input body")
		return
	}

	err := h.services.Authorization.CreateUser(c, input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	//StatusCode - 201 - User created
	c.Status(201)
}

func (h *Handler) singIn(c *gin.Context) {
	var input models.SignInInput

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Invalid input body")
		return
	}

	res, err := h.services.Authorization.GetUser(c.Request.Context(), models.SignInInput{
		Email:    input.Email,
		Password: input.Password,
	})

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, tokenResponse{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
	})
}

func (h *Handler) userRefresh(c *gin.Context) {
	type refreshInput struct {
		Token string `json:"token" binding:"required"`
	}
	var inp refreshInput
	if err := c.BindJSON(&inp); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid input body")

		return
	}

	res, err := h.services.RefreshTokens(c.Request.Context(), inp.Token)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.JSON(http.StatusOK, tokenResponse{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
	})
}
