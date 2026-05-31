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
func (ctrl *VersionController) RollbackVersion(c *gin.Context) {
	userID := c.GetUint("user_id")

	questionID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的题目ID")
		return
	}

	version, err := strconv.Atoi(c.Param("version"))
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
