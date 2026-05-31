#!/bin/bash

# QuestionHelper API 测试脚本
# 根据设计文档测试所有接口

BASE_URL="http://localhost:8080"
TOKEN=""
REFRESH_TOKEN=""

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 测试计数
TOTAL=0
PASSED=0
FAILED=0

# 打印测试结果
print_result() {
    local test_name="$1"
    local expected="$2"
    local actual="$3"

    TOTAL=$((TOTAL + 1))

    if echo "$actual" | grep -q "$expected"; then
        echo -e "${GREEN}✓ PASS${NC}: $test_name"
        PASSED=$((PASSED + 1))
    else
        echo -e "${RED}✗ FAIL${NC}: $test_name"
        echo -e "  Expected: $expected"
        echo -e "  Actual: $actual"
        FAILED=$((FAILED + 1))
    fi
}

# 打印测试分组
print_group() {
    echo ""
    echo -e "${YELLOW}========================================${NC}"
    echo -e "${YELLOW}$1${NC}"
    echo -e "${YELLOW}========================================${NC}"
}

# 延时函数（避免限流）
delay() {
    sleep 0.5
}

# 获取 Token
get_token() {
    local response=$(curl -s -X POST "$BASE_URL/api/v1/auth/login" \
        -H "Content-Type: application/json" \
        -d '{"username":"admin","password":"admin123"}')

    TOKEN=$(echo "$response" | grep -o '"accessToken":"[^"]*"' | cut -d'"' -f4)
    REFRESH_TOKEN=$(echo "$response" | grep -o '"refreshToken":"[^"]*"' | cut -d'"' -f4)

    if [ -z "$TOKEN" ]; then
        echo -e "${RED}Failed to get token${NC}"
        exit 1
    fi
    echo -e "${GREEN}Token obtained successfully${NC}"
    delay
}

# ==================== 认证系统测试 ====================
test_auth_system() {
    print_group "1. 用户认证系统测试"

    # 1.1 获取验证码
    local response=$(curl -s "$BASE_URL/api/v1/auth/captcha")
    print_result "1.1 获取验证码" "captchaId" "$response"
    delay

    # 1.2 用户注册
    response=$(curl -s -X POST "$BASE_URL/api/v1/auth/register" \
        -H "Content-Type: application/json" \
        -d '{"username":"testuser2","password":"Test123456"}')
    print_result "1.2 用户注册" "注册成功" "$response"
    delay

    # 1.3 用户登录
    response=$(curl -s -X POST "$BASE_URL/api/v1/auth/login" \
        -H "Content-Type: application/json" \
        -d '{"username":"admin","password":"admin123"}')
    print_result "1.3 用户登录" "accessToken" "$response"
    delay

    # 1.4 刷新Token
    response=$(curl -s -X POST "$BASE_URL/api/v1/auth/refresh" \
        -H "Content-Type: application/json" \
        -d "{\"refreshToken\":\"$REFRESH_TOKEN\"}")
    print_result "1.4 刷新Token" "accessToken" "$response"
    delay

    # 1.5 退出登录
    response=$(curl -s -X POST "$BASE_URL/api/v1/auth/logout" \
        -H "Authorization: Bearer $TOKEN")
    print_result "1.5 退出登录" "退出成功" "$response"
    delay

    # 重新登录获取Token
    get_token
}

# ==================== 用户管理测试 ====================
test_user_management() {
    print_group "2. 用户管理系统测试"

    # 2.1 获取用户列表
    local response=$(curl -s -H "Authorization: Bearer $TOKEN" "$BASE_URL/api/v1/admin/users")
    print_result "2.1 获取用户列表" "list" "$response"
    delay

    # 2.2 获取用户详情
    response=$(curl -s -H "Authorization: Bearer $TOKEN" "$BASE_URL/api/v1/admin/users/1")
    print_result "2.2 获取用户详情" "username" "$response"
    delay

    # 2.3 获取当前用户信息
    response=$(curl -s -H "Authorization: Bearer $TOKEN" "$BASE_URL/api/v1/users/me")
    print_result "2.3 获取当前用户信息" "username" "$response"
    delay

    # 2.4 获取用户角色列表
    response=$(curl -s -H "Authorization: Bearer $TOKEN" "$BASE_URL/api/v1/admin/roles")
    print_result "2.4 获取角色列表" "list" "$response"
    delay
}

# ==================== 菜单权限测试 ====================
test_menu_permission() {
    print_group "3. 菜单权限系统测试"

    # 3.1 获取用户菜单
    local response=$(curl -s -H "Authorization: Bearer $TOKEN" "$BASE_URL/api/v1/menus/routes")
    print_result "3.1 获取用户路由菜单" "name" "$response"
    delay

    # 3.2 获取菜单列表
    response=$(curl -s -H "Authorization: Bearer $TOKEN" "$BASE_URL/api/v1/admin/menus")
    print_result "3.2 获取菜单列表" "name" "$response"
    delay

    # 3.3 获取菜单树
    response=$(curl -s -H "Authorization: Bearer $TOKEN" "$BASE_URL/api/v1/admin/menus/tree")
    print_result "3.3 获取菜单树" "children" "$response"
    delay
}

# ==================== 题库管理测试 ====================
test_question_management() {
    print_group "4. 题库管理系统测试"

    # 4.1 获取题目列表
    local response=$(curl -s -H "Authorization: Bearer $TOKEN" "$BASE_URL/api/v1/questions")
    print_result "4.1 获取题目列表" "code" "$response"
    delay

    # 4.2 获取分类列表
    response=$(curl -s -H "Authorization: Bearer $TOKEN" "$BASE_URL/api/v1/categories")
    print_result "4.2 获取分类列表" "code" "$response"
    delay

    # 4.3 获取分类树
    response=$(curl -s -H "Authorization: Bearer $TOKEN" "$BASE_URL/api/v1/categories/tree")
    print_result "4.3 获取分类树" "code" "$response"
    delay

    # 4.4 获取知识点列表
    response=$(curl -s -H "Authorization: Bearer $TOKEN" "$BASE_URL/api/v1/knowledge-points")
    print_result "4.4 获取知识点列表" "code" "$response"
    delay
}

# ==================== 考试管理测试 ====================
test_exam_management() {
    print_group "5. 考试管理系统测试"

    # 5.1 获取试卷列表
    local response=$(curl -s -H "Authorization: Bearer $TOKEN" "$BASE_URL/api/v1/admin/papers")
    print_result "5.1 获取试卷列表" "code" "$response"
    delay

    # 5.2 获取考试列表
    response=$(curl -s -H "Authorization: Bearer $TOKEN" "$BASE_URL/api/v1/admin/exams")
    print_result "5.2 获取考试列表" "code" "$response"
    delay

    # 5.3 获取成绩列表
    response=$(curl -s -H "Authorization: Bearer $TOKEN" "$BASE_URL/api/v1/admin/scores")
    print_result "5.3 获取成绩列表" "code" "$response"
    delay
}

# ==================== 班级管理测试 ====================
test_class_management() {
    print_group "6. 班级管理系统测试"

    # 6.1 获取班级列表
    local response=$(curl -s -H "Authorization: Bearer $TOKEN" "$BASE_URL/api/v1/classes")
    print_result "6.1 获取班级列表" "code" "$response"
    delay

    # 6.2 获取管理员班级列表
    response=$(curl -s -H "Authorization: Bearer $TOKEN" "$BASE_URL/api/v1/admin/classes")
    print_result "6.2 获取管理员班级列表" "code" "$response"
    delay
}

# ==================== 练习管理测试 ====================
test_practice_management() {
    print_group "7. 练习管理系统测试"

    # 7.1 获取练习历史
    local response=$(curl -s -H "Authorization: Bearer $TOKEN" "$BASE_URL/api/v1/practices")
    print_result "7.1 获取练习历史" "code" "$response"
    delay

    # 7.2 获取练习统计
    response=$(curl -s -H "Authorization: Bearer $TOKEN" "$BASE_URL/api/v1/practices/stats")
    print_result "7.2 获取练习统计" "code" "$response"
    delay
}

# ==================== 错题本测试 ====================
test_wrong_questions() {
    print_group "8. 错题本系统测试"

    # 8.1 获取错题列表
    local response=$(curl -s -H "Authorization: Bearer $TOKEN" "$BASE_URL/api/v1/wrong-questions")
    print_result "8.1 获取错题列表" "code" "$response"
    delay

    # 8.2 获取错题分析
    response=$(curl -s -H "Authorization: Bearer $TOKEN" "$BASE_URL/api/v1/wrong-questions/analysis")
    print_result "8.2 获取错题分析" "code" "$response"
    delay
}

# ==================== 评论系统测试 ====================
test_comment_system() {
    print_group "9. 评论系统测试"

    # 9.1 获取评论列表
    local response=$(curl -s -H "Authorization: Bearer $TOKEN" "$BASE_URL/api/v1/comments")
    print_result "9.1 获取评论列表" "code" "$response"
    delay
}

# ==================== 通知系统测试 ====================
test_notification_system() {
    print_group "10. 通知系统测试"

    # 10.1 获取通知列表
    local response=$(curl -s -H "Authorization: Bearer $TOKEN" "$BASE_URL/api/v1/notifications")
    print_result "10.1 获取通知列表" "code" "$response"
    delay

    # 10.2 获取未读通知数
    response=$(curl -s -H "Authorization: Bearer $TOKEN" "$BASE_URL/api/v1/notifications/unread-count")
    print_result "10.2 获取未读通知数" "code" "$response"
    delay
}

# ==================== 统计分析测试 ====================
test_statistics() {
    print_group "11. 统计分析系统测试"

    # 11.1 获取统计概览
    local response=$(curl -s -H "Authorization: Bearer $TOKEN" "$BASE_URL/api/v1/statistics/overview")
    print_result "11.1 获取统计概览" "code" "$response"
    delay

    # 11.2 获取用户统计
    response=$(curl -s -H "Authorization: Bearer $TOKEN" "$BASE_URL/api/v1/statistics/users")
    print_result "11.2 获取用户统计" "code" "$response"
    delay

    # 11.3 获取练习统计
    response=$(curl -s -H "Authorization: Bearer $TOKEN" "$BASE_URL/api/v1/statistics/questions")
    print_result "11.3 获取练习统计" "code" "$response"
    delay

    # 11.4 获取考试统计
    response=$(curl -s -H "Authorization: Bearer $TOKEN" "$BASE_URL/api/v1/statistics/exams")
    print_result "11.4 获取考试统计" "code" "$response"
    delay

    # 11.5 获取班级统计
    response=$(curl -s -H "Authorization: Bearer $TOKEN" "$BASE_URL/api/v1/statistics/classes")
    print_result "11.5 获取班级统计" "code" "$response"
    delay
}

# ==================== 系统管理测试 ====================
test_system_management() {
    print_group "12. 系统管理测试"

    # 12.1 获取系统设置
    local response=$(curl -s -H "Authorization: Bearer $TOKEN" "$BASE_URL/api/v1/admin/configs")
    print_result "12.1 获取系统设置" "code" "$response"
    delay

    # 12.2 获取操作日志
    response=$(curl -s -H "Authorization: Bearer $TOKEN" "$BASE_URL/api/v1/admin/oper-logs")
    print_result "12.2 获取操作日志" "code" "$response"
    delay

    # 12.3 获取登录日志
    response=$(curl -s -H "Authorization: Bearer $TOKEN" "$BASE_URL/api/v1/admin/login-logs")
    print_result "12.3 获取登录日志" "code" "$response"
    delay
}

# ==================== 文件管理测试 ====================
test_file_management() {
    print_group "13. 文件管理测试"

    # 13.1 测试文件上传（创建测试文件）
    echo "test content" > /tmp/test_upload.txt
    local response=$(curl -s -X POST "$BASE_URL/api/v1/files" \
        -H "Authorization: Bearer $TOKEN" \
        -F "file=@/tmp/test_upload.txt")
    print_result "13.1 文件上传" "code" "$response"
    rm -f /tmp/test_upload.txt
    delay
}

# ==================== 主测试流程 ====================
main() {
    echo "=========================================="
    echo "QuestionHelper API 测试"
    echo "=========================================="
    echo ""

    # 获取初始Token
    get_token

    # 运行所有测试
    test_auth_system
    test_user_management
    test_menu_permission
    test_question_management
    test_exam_management
    test_class_management
    test_practice_management
    test_wrong_questions
    test_comment_system
    test_notification_system
    test_statistics
    test_system_management
    test_file_management

    # 打印测试总结
    echo ""
    echo "=========================================="
    echo "测试总结"
    echo "=========================================="
    echo -e "总计: $TOTAL"
    echo -e "${GREEN}通过: $PASSED${NC}"
    echo -e "${RED}失败: $FAILED${NC}"

    if [ $FAILED -eq 0 ]; then
        echo -e "\n${GREEN}所有测试通过！${NC}"
        exit 0
    else
        echo -e "\n${RED}有测试失败，请检查！${NC}"
        exit 1
    fi
}

# 运行测试
main
