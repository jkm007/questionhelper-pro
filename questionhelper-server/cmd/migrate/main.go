package main

import (
	"fmt"
	"log"

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

	err = db.AutoMigrate(
		// User & Auth
		&model.User{},
		&model.UserPrivacy{},
		&model.OAuthUser{},
		&model.LoginDevice{},
		&model.SecurityLog{},
		&model.PasswordHistory{},
		&model.Role{},
		&model.Menu{},
		&model.Permission{},
		&model.Tag{},
		&model.UserTag{},
		&model.UserRealName{},
		&model.RoleApplication{},
		&model.DataPermission{},
		&model.Dept{},
		// Question
		&model.Question{},
		&model.Option{},
		&model.Category{},
		&model.Knowledge{},
		&model.QuestionVersion{},
		&model.FavoriteFolder{},
		&model.QuestionFavorite{},
		&model.QuestionShare{},
		&model.QuestionAttachment{},
		&model.QuestionReview{},
		&model.SensitiveWord{},
		// Exam
		&model.Exam{},
		&model.Paper{},
		&model.PaperQuestion{},
		&model.ExamRecord{},
		&model.AnswerRecord{},
		&model.ExamWarning{},
		// Class
		&model.Class{},
		&model.ClassMember{},
		&model.Homework{},
		&model.ClassNotice{},
		// Comment
		&model.Comment{},
		&model.CommentLike{},
		&model.CommentReport{},
		// Notification
		&model.Notification{},
		// Practice
		&model.PracticeSession{},
		&model.PracticeRecord{},
		// Wrong
		&model.WrongQuestion{},
		// File
		&model.File{},
		// System
		&model.SystemSetting{},
		&model.OperationLog{},
		&model.LoginLog{},
	)

	if err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}

	fmt.Println("数据库迁移成功!")
}
