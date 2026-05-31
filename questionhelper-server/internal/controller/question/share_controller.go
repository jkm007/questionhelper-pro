package question

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"questionhelper-server/internal/dto"
	"questionhelper-server/internal/service/question"
	"questionhelper-server/pkg/response"
)

type ShareController struct{}

func NewShareController() *ShareController {
	return &ShareController{}
}

// CreateShare 创建分享
func (ctrl *ShareController) CreateShare(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req dto.CreateShareRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	info, err := question.CreateShare(userID, &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, info)
}

// GetShare 获取分享
func (ctrl *ShareController) GetShare(c *gin.Context) {
	code := c.Param("code")

	info, err := question.GetShare(code)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, info)
}

// RevokeShare 撤销分享
func (ctrl *ShareController) RevokeShare(c *gin.Context) {
	userID := c.GetUint("user_id")

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的分享ID")
		return
	}

	if err := question.RevokeShare(uint(id), userID); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "撤销成功", nil)
}

// ListMyShares 我的分享列表
func (ctrl *ShareController) ListMyShares(c *gin.Context) {
	userID := c.GetUint("user_id")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	list, total, err := question.ListMyShares(userID, page, pageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Page(c, list, total, page, pageSize)
}
