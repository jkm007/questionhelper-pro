package question

import (
	"errors"
	"fmt"

	"questionhelper-server/internal/dto"
	"questionhelper-server/internal/model"
	questionRepo "questionhelper-server/internal/repository/question"
	"questionhelper-server/pkg/logger"
)

// ListCategories 分类列表
func ListCategories() ([]dto.CategoryInfo, error) {
	categories, err := questionRepo.FindAllCategories()
	if err != nil {
		return nil, fmt.Errorf("查询分类失败: %w", err)
	}

	list := make([]dto.CategoryInfo, 0, len(categories))
	for _, c := range categories {
		list = append(list, toCategoryInfo(&c))
	}
	return list, nil
}

// GetCategoryTree 分类树
func GetCategoryTree() ([]dto.CategoryInfo, error) {
	categories, err := questionRepo.FindCategoryTree()
	if err != nil {
		return nil, fmt.Errorf("查询分类树失败: %w", err)
	}

	tree := make([]dto.CategoryInfo, 0, len(categories))
	for _, c := range categories {
		tree = append(tree, toCategoryInfo(&c))
	}
	return tree, nil
}

// CreateCategory 创建分类
func CreateCategory(req *dto.CreateCategoryRequest) error {
	category := &model.Category{
		ParentID: req.ParentID,
		Name:     req.Name,
		Sort:     req.Sort,
	}

	if err := questionRepo.CreateCategory(category); err != nil {
		return fmt.Errorf("创建分类失败: %w", err)
	}

	logger.Infof("创建分类成功: %s", req.Name)
	return nil
}

// UpdateCategory 更新分类
func UpdateCategory(id uint, req *dto.UpdateCategoryRequest) error {
	category, err := questionRepo.FindCategoryByID(id)
	if err != nil {
		return errors.New("分类不存在")
	}

	if req.Name != "" {
		category.Name = req.Name
	}
	if req.ParentID != nil {
		category.ParentID = req.ParentID
	}
	category.Sort = req.Sort

	if err := questionRepo.UpdateCategory(category); err != nil {
		return fmt.Errorf("更新分类失败: %w", err)
	}

	logger.Infof("更新分类 %d 成功", id)
	return nil
}

// DeleteCategory 删除分类
func DeleteCategory(id uint) error {
	if err := questionRepo.DeleteCategoryByID(id); err != nil {
		return fmt.Errorf("删除分类失败: %w", err)
	}

	logger.Infof("删除分类 %d 成功", id)
	return nil
}

// toCategoryInfo 转换为 CategoryInfo DTO
func toCategoryInfo(c *model.Category) dto.CategoryInfo {
	info := dto.CategoryInfo{
		ID:       c.ID,
		ParentID: c.ParentID,
		Name:     c.Name,
	}
	if len(c.Children) > 0 {
		info.Children = make([]dto.CategoryInfo, 0, len(c.Children))
		for _, child := range c.Children {
			info.Children = append(info.Children, toCategoryInfo(&child))
		}
	}
	return info
}
