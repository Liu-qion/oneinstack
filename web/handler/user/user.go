package user

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"oneinstack/core"
	"oneinstack/internal/services/user"
	"oneinstack/utils"
	"oneinstack/web/input"
)

func LoginHandler(c *gin.Context) {
	var req input.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		core.HandleError(c, http.StatusUnauthorized, core.ErrBadRequest, err)
		return
	}
	user, ok := user.CheckUserPassword(req.Username, req.Password)
	if !ok {
		core.HandleError(c, http.StatusUnauthorized, core.ErrUnauthorizedAP, nil)
		return
	}
	token, err := utils.GenerateJWT(user.Username)
	if err != nil {
		core.HandleError(c, http.StatusInternalServerError, core.ErrInternalServerError, gin.H{"error": "could not generate token"})
		return
	}
	core.HandleSuccess(c, gin.H{"token": token})
}
