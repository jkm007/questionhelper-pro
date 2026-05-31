package user

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"strconv"
	"strings"

	"questionhelper-server/internal/dto"
	userRepo "questionhelper-server/internal/repository/user"
)

// ExportUsers 导出用户为CSV
func ExportUsers(req *dto.UserListRequest) (*bytes.Buffer, error) {
	users, err := userRepo.ExportUsers(req)
	if err != nil {
		return nil, fmt.Errorf("查询用户数据失败: %w", err)
	}

	buf := &bytes.Buffer{}

	// 写入BOM(解决Excel中文乱码)
	buf.WriteString("\xEF\xBB\xBF")

	w := csv.NewWriter(buf)

	// 写入表头
	if err := w.Write([]string{"ID", "用户名", "昵称", "手机号", "邮箱", "性别", "状态", "实名认证", "注册时间"}); err != nil {
		return nil, fmt.Errorf("写入CSV表头失败: %w", err)
	}

	// 写入数据
	for _, u := range users {
		gender := "未知"
		if u.Gender == 1 {
			gender = "男"
		} else if u.Gender == 2 {
			gender = "女"
		}

		status := "禁用"
		if u.Status == 1 {
			status = "正常"
		} else if u.Status == 2 {
			status = "注销中"
		}

		isReal := "否"
		if u.IsReal {
			isReal = "是"
		}

		record := []string{
			strconv.FormatUint(uint64(u.ID), 10),
			u.Username,
			u.Nickname,
			u.Phone,
			u.Email,
			gender,
			status,
			isReal,
			u.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		if err := w.Write(record); err != nil {
			return nil, fmt.Errorf("写入CSV数据失败: %w", err)
		}
	}

	w.Flush()
	if err := w.Error(); err != nil {
		return nil, fmt.Errorf("刷新CSV写入器失败: %w", err)
	}

	return buf, nil
}

// ExportUsersByTag 按标签导出用户
func ExportUsersByTag(tagID uint) (*bytes.Buffer, error) {
	users, _, err := userRepo.FindByTagID(tagID, 1, 10000)
	if err != nil {
		return nil, fmt.Errorf("查询用户数据失败: %w", err)
	}

	buf := &bytes.Buffer{}

	// 写入BOM
	buf.WriteString("\xEF\xBB\xBF")

	w := csv.NewWriter(buf)

	// 写入表头
	if err := w.Write([]string{"ID", "用户名", "昵称", "手机号", "邮箱", "状态"}); err != nil {
		return nil, fmt.Errorf("写入CSV表头失败: %w", err)
	}

	// 写入数据
	for _, u := range users {
		status := "禁用"
		if u.Status == 1 {
			status = "正常"
		} else if u.Status == 2 {
			status = "注销中"
		}

		record := []string{
			strconv.FormatUint(uint64(u.ID), 10),
			u.Username,
			u.Nickname,
			u.Phone,
			u.Email,
			status,
		}
		if err := w.Write(record); err != nil {
			return nil, fmt.Errorf("写入CSV数据失败: %w", err)
		}
	}

	w.Flush()
	if err := w.Error(); err != nil {
		return nil, fmt.Errorf("刷新CSV写入器失败: %w", err)
	}

	return buf, nil
}

// ExportUsersByID 根据ID列表导出用户
func ExportUsersByID(ids []uint) (*bytes.Buffer, error) {
	buf := &bytes.Buffer{}

	// 写入BOM
	buf.WriteString("\xEF\xBB\xBF")

	w := csv.NewWriter(buf)

	// 写入表头
	if err := w.Write([]string{"ID", "用户名", "昵称", "手机号", "邮箱", "性别", "状态", "角色"}); err != nil {
		return nil, fmt.Errorf("写入CSV表头失败: %w", err)
	}

	for _, id := range ids {
		u, err := userRepo.FindByID(id)
		if err != nil {
			continue
		}

		gender := "未知"
		if u.Gender == 1 {
			gender = "男"
		} else if u.Gender == 2 {
			gender = "女"
		}

		status := "禁用"
		if u.Status == 1 {
			status = "正常"
		} else if u.Status == 2 {
			status = "注销中"
		}

		roles := ""
		for i, role := range u.Roles {
			if i > 0 {
				roles += "|"
			}
			roles += role.Name
		}

		record := []string{
			strconv.FormatUint(uint64(u.ID), 10),
			u.Username,
			u.Nickname,
			u.Phone,
			u.Email,
			gender,
			status,
			roles,
		}
		if err := w.Write(record); err != nil {
			return nil, fmt.Errorf("写入CSV数据失败: %w", err)
		}
	}

	w.Flush()
	if err := w.Error(); err != nil {
		return nil, fmt.Errorf("刷新CSV写入器失败: %w", err)
	}

	return buf, nil
}

// importUsers 从CSV导入用户
func importUsers(data []byte) (int, error) {
	// 去掉BOM前缀
	cleanData := data
	if len(cleanData) >= 3 && cleanData[0] == 0xEF && cleanData[1] == 0xBB && cleanData[2] == 0xBF {
		cleanData = cleanData[3:]
	}

	reader := csv.NewReader(bytes.NewReader(cleanData))
	reader.TrimLeadingSpace = true

	records, err := reader.ReadAll()
	if err != nil {
		return 0, fmt.Errorf("解析CSV数据失败: %w", err)
	}

	if len(records) <= 1 {
		return 0, fmt.Errorf("数据为空")
	}

	count := 0
	// 跳过表头，从第二行开始
	for _, record := range records[1:] {
		if len(record) < 3 {
			continue
		}

		username := strings.TrimSpace(record[0])
		password := strings.TrimSpace(record[1])
		nickname := strings.TrimSpace(record[2])

		if username == "" || password == "" {
			continue
		}

		// 创建用户
		req := &dto.CreateUserRequest{
			Username: username,
			Password: password,
			Nickname: nickname,
		}

		if err := CreateUser(req); err != nil {
			continue
		}
		count++
	}

	return count, nil
}

// ExportUsersExcel 导出用户为Excel格式(简化版，使用制表符分隔)
func ExportUsersExcel(req *dto.UserListRequest) (*bytes.Buffer, error) {
	users, err := userRepo.ExportUsers(req)
	if err != nil {
		return nil, fmt.Errorf("查询用户数据失败: %w", err)
	}

	buf := &bytes.Buffer{}

	// 写入表头
	buf.WriteString("ID\t用户名\t昵称\t手机号\t邮箱\t性别\t状态\t实名认证\t注册时间\n")

	// 写入数据
	for _, u := range users {
		gender := "未知"
		if u.Gender == 1 {
			gender = "男"
		} else if u.Gender == 2 {
			gender = "女"
		}

		status := "禁用"
		if u.Status == 1 {
			status = "正常"
		} else if u.Status == 2 {
			status = "注销中"
		}

		isReal := "否"
		if u.IsReal {
			isReal = "是"
		}

		line := fmt.Sprintf("%d\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
			u.ID,
			u.Username,
			u.Nickname,
			u.Phone,
			u.Email,
			gender,
			status,
			isReal,
			u.CreatedAt.Format("2006-01-02 15:04:05"),
		)
		buf.WriteString(line)
	}

	return buf, nil
}

// Atoi 安全的字符串转整数
func safeAtoi(s string) int {
	n, _ := strconv.Atoi(s)
	return n
}
