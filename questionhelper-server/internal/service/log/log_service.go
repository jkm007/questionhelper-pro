package log

import (
	"bytes"
	"fmt"
	"time"

	"questionhelper-server/internal/dto"
	"questionhelper-server/internal/model"
	logRepo "questionhelper-server/internal/repository/log"
	"questionhelper-server/pkg/logger"
)

// CreateOperationLog 创建操作日志
func CreateOperationLog(log *model.OperationLog) error {
	return logRepo.CreateOperationLog(log)
}

// GetOperationLog 获取操作日志详情
func GetOperationLog(id uint) (*dto.OperationLogInfo, error) {
	log, err := logRepo.FindOperationLogByID(id)
	if err != nil {
		return nil, fmt.Errorf("查询操作日志失败: %w", err)
	}
	return toOperationLogInfo(log), nil
}

// ListOperationLogs 操作日志列表
func ListOperationLogs(req *dto.OperationLogListRequest) ([]dto.OperationLogInfo, int64, error) {
	logs, total, err := logRepo.ListOperationLogs(req)
	if err != nil {
		return nil, 0, fmt.Errorf("查询操作日志列表失败: %w", err)
	}

	list := make([]dto.OperationLogInfo, 0, len(logs))
	for _, log := range logs {
		list = append(list, *toOperationLogInfo(&log))
	}

	return list, total, nil
}

// ExportOperationLogs 导出操作日志
func ExportOperationLogs(req *dto.OperationLogListRequest) (*bytes.Buffer, error) {
	// 不分页，导出全部
	req.Page = 1
	req.PageSize = 10000

	logs, _, err := logRepo.ListOperationLogs(req)
	if err != nil {
		return nil, fmt.Errorf("查询操作日志失败: %w", err)
	}

	buf := &bytes.Buffer{}

	// 写入BOM
	buf.WriteString("\xEF\xBB\xBF")

	// 写入表头
	buf.WriteString("ID,用户ID,用户名,模块,操作,资源,描述,IP,状态,时间\n")

	// 写入数据
	for _, log := range logs {
		status := "失败"
		if log.Status == 1 {
			status = "成功"
		}

		line := fmt.Sprintf("%d,%d,%s,%s,%s,%s,%s,%s,%s,%s\n",
			log.ID,
			log.UserID,
			log.Username,
			log.Module,
			log.Action,
			log.Resource,
			log.Description,
			log.IP,
			status,
			log.CreatedAt.Format("2006-01-02 15:04:05"),
		)
		buf.WriteString(line)
	}

	return buf, nil
}

// CleanOperationLogs 清理过期操作日志
func CleanOperationLogs(days int) (int64, error) {
	before := time.Now().AddDate(0, 0, -days)
	count, err := logRepo.DeleteOperationLogsBefore(before)
	if err != nil {
		return 0, fmt.Errorf("清理操作日志失败: %w", err)
	}
	logger.Infof("清理 %d 天前的操作日志，共 %d 条", days, count)
	return count, nil
}

// ListLoginLogs 登录日志列表
func ListLoginLogs(req *dto.LoginLogListRequest) ([]dto.LoginLogInfo, int64, error) {
	logs, total, err := logRepo.ListLoginLogs(req)
	if err != nil {
		return nil, 0, fmt.Errorf("查询登录日志列表失败: %w", err)
	}

	list := make([]dto.LoginLogInfo, 0, len(logs))
	for _, log := range logs {
		list = append(list, *toLoginLogInfo(&log))
	}

	return list, total, nil
}

// ExportLoginLogs 导出登录日志
func ExportLoginLogs(req *dto.LoginLogListRequest) (*bytes.Buffer, error) {
	// 不分页，导出全部
	req.Page = 1
	req.PageSize = 10000

	logs, _, err := logRepo.ListLoginLogs(req)
	if err != nil {
		return nil, fmt.Errorf("查询登录日志失败: %w", err)
	}

	buf := &bytes.Buffer{}

	// 写入BOM
	buf.WriteString("\xEF\xBB\xBF")

	// 写入表头
	buf.WriteString("ID,用户ID,用户名,IP,位置,浏览器,操作系统,状态,消息,时间\n")

	// 写入数据
	for _, log := range logs {
		status := "失败"
		if log.Status == 1 {
			status = "成功"
		}

		line := fmt.Sprintf("%d,%d,%s,%s,%s,%s,%s,%s,%s,%s\n",
			log.ID,
			log.UserID,
			log.Username,
			log.IP,
			log.Location,
			log.Browser,
			log.OS,
			status,
			log.Msg,
			log.CreatedAt.Format("2006-01-02 15:04:05"),
		)
		buf.WriteString(line)
	}

	return buf, nil
}

// CleanLoginLogs 清理过期登录日志
func CleanLoginLogs(days int) (int64, error) {
	before := time.Now().AddDate(0, 0, -days)
	count, err := logRepo.DeleteLoginLogsBefore(before)
	if err != nil {
		return 0, fmt.Errorf("清理登录日志失败: %w", err)
	}
	logger.Infof("清理 %d 天前的登录日志，共 %d 条", days, count)
	return count, nil
}

// toOperationLogInfo 转换为 OperationLogInfo DTO
func toOperationLogInfo(log *model.OperationLog) *dto.OperationLogInfo {
	return &dto.OperationLogInfo{
		ID:          log.ID,
		UserID:      log.UserID,
		Username:    log.Username,
		Module:      log.Module,
		Action:      log.Action,
		Resource:    log.Resource,
		ResourceID:  log.ResourceID,
		Description: log.Description,
		IP:          log.IP,
		UserAgent:   log.UserAgent,
		Status:      log.Status,
		ErrorMsg:    log.ErrorMsg,
		CreatedAt:   log.CreatedAt,
	}
}

// toLoginLogInfo 转换为 LoginLogInfo DTO
func toLoginLogInfo(log *model.LoginLog) *dto.LoginLogInfo {
	var userID uint
	if log.UserID != nil {
		userID = *log.UserID
	}
	return &dto.LoginLogInfo{
		ID:        log.ID,
		UserID:    userID,
		Username:  log.Username,
		IP:        log.IP,
		Location:  log.Location,
		Browser:   log.Browser,
		OS:        log.OS,
		Status:    log.Status,
		Msg:       log.Msg,
		CreatedAt: log.CreatedAt,
	}
}
