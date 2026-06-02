package user

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"questionhelper-server/internal/dto"
	"questionhelper-server/internal/service/user"
	"questionhelper-server/pkg/response"
)

type ApplyController struct{}

func NewApplyController() *ApplyController {
	return &ApplyController{}
}

// CreateApplication 创建角色申请
// @Summary      创建角色申请
// @Description  提交角色申请
// @Tags         角色申请
// @Accept       json
// @Produce      json
// @Param        req  body      dto.CreateApplicationRequest  true  "申请信息"
// @Success      200  {object}  response.Response  "申请提交成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /applications [post]
// @Security     BearerAuth
func (ctrl *ApplyController) CreateApplication(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req dto.CreateApplicationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := user.CreateApplication(userID, &req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "申请提交成功", nil)
}

// GetApplication 获取申请详情
// @Summary      获取申请详情
// @Description  根据ID获取角色申请详情
// @Tags         角色申请
// @Accept       json
// @Produce      json
// @Param        id   path      uint  true  "申请ID"
// @Success      200  {object}  response.Response{data=dto.ApplicationInfo}  "成功"
// @Failure      400  {object}  response.Response  "无效的申请ID"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /applications/{id} [get]
// @Security     BearerAuth
func (ctrl *ApplyController) GetApplication(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的申请ID")
		return
	}

	info, err := user.GetApplication(uint(id))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, info)
}

// ListApplications 申请列表(管理员)
// @Summary      申请列表
// @Description  分页查询角色申请列表（管理员）
// @Tags         角色申请
// @Accept       json
// @Produce      json
// @Param        page       query     int     false  "页码"       default(1)
// @Param        page_size  query     int     false  "每页数量"   default(10)
// @Param        status     query     int     false  "审核状态"
// @Success      200  {object}  response.Response{data=dto.PageResponse{list=[]dto.ApplicationInfo}}  "成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/applications [get]
// @Security     BearerAuth
func (ctrl *ApplyController) ListApplications(c *gin.Context) {
	var req dto.ApplicationListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	list, total, err := user.ListApplications(&req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Page(c, list, total, req.Page, req.PageSize)
}

// ReviewApplication 审核角色申请(管理员)
// @Summary      审核角色申请
// @Description  审核角色申请（管理员）
// @Tags         角色申请
// @Accept       json
// @Produce      json
// @Param        id   path      uint                            true  "申请ID"
// @Param        req  body      dto.ReviewApplicationRequest    true  "审核信息"
// @Success      200  {object}  response.Response  "审核完成"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/applications/{id}/review [put]
// @Security     BearerAuth
func (ctrl *ApplyController) ReviewApplication(c *gin.Context) {
	reviewerID := c.GetUint("user_id")

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的申请ID")
		return
	}

	var req dto.ReviewApplicationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := user.ReviewApplication(uint(id), reviewerID, &req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "审核完成", nil)
}

// GetMyApplications 获取我的申请列表
// @Summary      获取我的申请列表
// @Description  获取当前用户的角色申请列表
// @Tags         角色申请
// @Accept       json
// @Produce      json
// @Param        page       query     int  false  "页码"       default(1)
// @Param        page_size  query     int  false  "每页数量"   default(10)
// @Success      200  {object}  response.Response{data=dto.PageResponse{list=[]dto.ApplicationInfo}}  "成功"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /applications [get]
// @Security     BearerAuth
func (ctrl *ApplyController) GetMyApplications(c *gin.Context) {
	userID := c.GetUint("user_id")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	list, total, err := user.GetUserApplications(userID, page, pageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Page(c, list, total, page, pageSize)
}
