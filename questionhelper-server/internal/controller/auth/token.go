package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"questionhelper-server/internal/dto"
	"questionhelper-server/internal/service/auth"
	"questionhelper-server/pkg/response"
)

// RefreshToken 刷新令牌
func (ac *AuthController) RefreshToken(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	result, err := auth.RefreshToken(&req, ac.jwtCfg)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, err.Error())
		return
	}

	response.Success(c, result)
}
