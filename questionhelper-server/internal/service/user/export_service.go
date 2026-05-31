package user

import (
	"bytes"
	"fmt"
	"strconv"

	"questionhelper-server/internal/dto"
	userRepo "questionhelper-server/internal/repository/user"
)

// ExportUsers 导出用户为CSV
func ExportUsers(req *dto.UserListRequest) (*bytes.Buffer, error) {
	users, err := userRepo.ExportUsers(req)
	if err != nil {
		return nil, fmt.Errorf("查询用户数据失败: %w", err)
	}

	// 创建CSV内容
	buf := &bytes.Buffer{}

	// 写入BOM(解决Excel中文乱码)
	buf.WriteString("\xEF\xBB\xBF")

	// 写入表头
	buf.WriteString("ID,用户名,昵称,手机号,邮箱,性别,状态,实名认证,注册时间\n")

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

		line := fmt.Sprintf("%d,%s,%s,%s,%s,%s,%s,%s,%s\n",
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

// ExportUsersByTag 按标签导出用户
func ExportUsersByTag(tagID uint) (*bytes.Buffer, error) {
	users, _, err := userRepo.FindByTagID(tagID, 1, 10000)
	if err != nil {
		return nil, fmt.Errorf("查询用户数据失败: %w", err)
	}

	// 创建CSV内容
	buf := &bytes.Buffer{}

	// 写入BOM
	buf.WriteString("\xEF\xBB\xBF")

	// 写入表头
	buf.WriteString("ID,用户名,昵称,手机号,邮箱,状态\n")

	// 写入数据
	for _, u := range users {
		status := "禁用"
		if u.Status == 1 {
			status = "正常"
		} else if u.Status == 2 {
			status = "注销中"
		}

		line := fmt.Sprintf("%d,%s,%s,%s,%s,%s\n",
			u.ID,
			u.Username,
			u.Nickname,
			u.Phone,
			u.Email,
			status,
		)
		buf.WriteString(line)
	}

	return buf, nil
}

// ExportUsersByID 根据ID列表导出用户
func ExportUsersByID(ids []uint) (*bytes.Buffer, error) {
	buf := &bytes.Buffer{}

	// 写入BOM
	buf.WriteString("\xEF\xBB\xBF")

	// 写入表头
	buf.WriteString("ID,用户名,昵称,手机号,邮箱,性别,状态,角色\n")

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

		line := fmt.Sprintf("%d,%s,%s,%s,%s,%s,%s,%s\n",
			u.ID,
			u.Username,
			u.Nickname,
			u.Phone,
			u.Email,
			gender,
			status,
			roles,
		)
		buf.WriteString(line)
	}

	return buf, nil
}

// importUsers 从CSV导入用户(示例实现)
func importUsers(data []byte) (int, error) {
	lines := bytes.Split(data, []byte("\n"))
	if len(lines) <= 1 {
		return 0, fmt.Errorf("数据为空")
	}

	count := 0
	// 跳过表头
	for i := 1; i < len(lines); i++ {
		line := string(bytes.TrimSpace(lines[i]))
		if line == "" {
			continue
		}

		// 解析CSV行(简单实现，生产环境应使用csv库)
		fields := bytes.Split(lines[i], []byte(","))
		if len(fields) < 3 {
			continue
		}

		username := string(fields[0])
		password := string(fields[1])
		nickname := string(fields[2])

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

// parseCSVLine 解析CSV行
func parseCSVLine(line string) []string {
	// 简单实现，生产环境应使用encoding/csv
	fields := []string{}
	field := ""
	inQuote := false

	for _, c := range line {
		switch {
		case c == '"':
			inQuote = !inQuote
		case c == ',' && !inQuote:
			fields = append(fields, field)
			field = ""
		default:
			field += string(c)
		}
	}
	fields = append(fields, field)

	return fields
}

// Atoi 安全的字符串转整数
func safeAtoi(s string) int {
	n, _ := strconv.Atoi(s)
	return n
}
