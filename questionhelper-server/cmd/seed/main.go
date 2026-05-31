package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"questionhelper-server/internal/model"
	"questionhelper-server/pkg/config"
	"questionhelper-server/pkg/database"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	database.InitMySQL(cfg.MySQL)
	db := database.DB

	// 先执行自动迁移
	err = db.AutoMigrate(&model.Menu{})
	if err != nil {
		log.Fatalf("迁移失败: %v", err)
	}

	// 检查是否已有菜单数据
	var count int64
	db.Model(&model.Menu{}).Count(&count)
	if count > 0 {
		fmt.Printf("菜单表已有 %d 条数据，是否清空后重新初始化？(y/n): ", count)
		var input string
		fmt.Scanln(&input)
		if input != "y" && input != "Y" {
			fmt.Println("取消操作")
			os.Exit(0)
		}
		// 清空菜单表（先删除子记录）
		db.Exec("DELETE FROM role_menus")
		db.Exec("SET FOREIGN_KEY_CHECKS = 0")
		db.Exec("DELETE FROM menus")
		db.Exec("SET FOREIGN_KEY_CHECKS = 1")
		db.Exec("ALTER TABLE menus AUTO_INCREMENT = 1")
		fmt.Println("已清空菜单表")
	}

	// 读取SQL文件
	sqlFile := "scripts/menu_seed.sql"
	file, err := os.Open(sqlFile)
	if err != nil {
		log.Fatalf("打开SQL文件失败: %v", err)
	}
	defer file.Close()

	// 逐行解析并执行SQL
	scanner := bufio.NewScanner(file)
	var sqlBuilder strings.Builder
	lineNum := 0
	executedCount := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		// 跳过空行和注释
		if line == "" || strings.HasPrefix(line, "--") || strings.HasPrefix(line, "#") {
			continue
		}

		sqlBuilder.WriteString(line)
		sqlBuilder.WriteString(" ")

		// 如果行以分号结尾，执行SQL
		if strings.HasSuffix(line, ";") {
			sql := strings.TrimSpace(sqlBuilder.String())
			if sql != "" && !strings.HasPrefix(sql, "SELECT") {
				result := db.Exec(sql)
				if result.Error != nil {
					log.Printf("执行SQL失败 (行 %d): %v\nSQL: %s", lineNum, result.Error, sql[:min(100, len(sql))])
				} else {
					executedCount++
				}
			}
			sqlBuilder.Reset()
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("读取SQL文件失败: %v", err)
	}

	fmt.Printf("菜单种子数据初始化成功！共执行 %d 条SQL语句\n", executedCount)

	// 验证结果
	var menuCount int64
	db.Model(&model.Menu{}).Count(&menuCount)
	fmt.Printf("当前菜单总数: %d\n", menuCount)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
