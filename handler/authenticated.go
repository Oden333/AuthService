package handler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (h *Handler) getUserIdentity(c *gin.Context) {
	id, err := h.parseAuthHeader(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
	}
	c.Set("userId", id)
}

func (h *Handler) parseAuthHeader(c *gin.Context) (string, error) {
	header := c.GetHeader("Authorization")
	if header == "" {
		return "", errors.New("empty auth header")
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return "", errors.New("invalid auth header")
	}

	if len(headerParts[1]) == 0 {
		return "", errors.New("token is empty")
	}

	return h.services.ParseToken(headerParts[1])
}

func (h *Handler) getUserAccount(c *gin.Context) {
	idFromCtx, ok := c.Get("userId")
	if !ok {
		newErrorResponse(c, http.StatusInternalServerError, "User context not found")
	}

	idStr, ok := idFromCtx.(string)
	if !ok {
		newErrorResponse(c, http.StatusInternalServerError, "Ctx is of invalid type")
	}

	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "Error finding user by Id")
	}
	user, err := h.services.GetUserById(c.Request.Context(), id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]string{
		"Name":  user.Name,
		"Email": user.Email,
	})
}
