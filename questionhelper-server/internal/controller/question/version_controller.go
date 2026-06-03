package question

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"questionhelper-server/internal/service/question"
	"questionhelper-server/pkg/response"
)

type VersionController struct{}

func NewVersionController() *VersionController {
	return &VersionController{}
}

// ListVersions 获取版本列表
// @Summary      获取题目版本列表
// @Description  获取指定题目的所有版本记录
// @Tags         题目版本
// @Accept       json
// @Produce      json
// @Param        id   path      uint  true  "题目ID"
// @Success      200  {object}  response.Response{data=[]dto.VersionInfo}  "成功"
// @Failure      400  {object}  response.Response  "无效的题目ID"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /questions/{id}/versions [get]
// @Security     BearerAuth
func (ctrl *VersionController) ListVersions(c *gin.Context) {
	questionID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的题目ID")
		return
	}

	list, err := question.ListVersions(uint(questionID))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, list)
}

// GetVersionDetail 获取版本详情
// @Summary      获取版本详情
// @Description  管理员获取指定版本的详细信息
// @Tags         题目版本
// @Accept       json
// @Produce      json
// @Param        id          path      uint  true  "题目ID"
// @Param        versionId   path      uint  true  "版本ID"
// @Success      200  {object}  response.Response{data=dto.VersionDetail}  "成功"
// @Failure      400  {object}  response.Response  "无效的版本ID"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/questions/{id}/versions/{versionId} [get]
// @Security     BearerAuth
func (ctrl *VersionController) GetVersionDetail(c *gin.Context) {
	versionID, err := strconv.ParseUint(c.Param("versionId"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的版本ID")
		return
	}

	detail, err := question.GetVersionDetail(uint(versionID))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, detail)
}

// RollbackVersion 回滚版本
// @Summary      回滚版本
// @Description  管理员将题目回滚到指定版本
// @Tags         题目版本
// @Accept       json
// @Produce      json
// @Param        id        path      uint  true  "题目ID"
// @Param        version   path      int   true  "版本号"
// @Success      200  {object}  response.Response  "回滚成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/questions/{id}/versions/{version}/rollback [post]
// @Security     BearerAuth
func (ctrl *VersionController) RollbackVersion(c *gin.Context) {
	userID := c.GetUint("user_id")

	questionID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的题目ID")
		return
	}

	version, err := strconv.Atoi(c.Param("versionId"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的版本号")
		return
	}

	if err := question.RollbackVersion(int(questionID), version, userID); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "回滚成功", nil)
}
