# 移动端图片资源

> 来源：`questionhelper-app/static/` 及子目录

## 目录结构

```
mobile/
├── tabbar/          # 底部导航栏图标 (10个)
├── icons/           # 功能图标 (28个)
├── empty/           # 空状态图 (7个，含IP形象)
├── error/           # 错误页 (3个，含IP形象场景插画)
├── status/          # 状态图标 (4个)
├── defaults/        # 默认占位图 (3个)
├── banners/         # Banner 图 (3个，含IP形象)
├── logo.png         # 移动端 Logo
└── splash-icon.png  # 启动页图标
```

## 重新设计说明

本次对移动端图片进行了**品牌化重新设计**：

| 类别 | 改进前 | 改进后 |
|------|--------|--------|
| **功能图标** | 低清像素化色块 | 矢量SVG渲染，品牌色渐变 + 圆角背景 |
| **错误页** | 单图标居中，无场景 | IP形象场景插画（404挠头/401挡门/网络断连） |
| **空状态** | 圆圈套通用图标 | IP形象表情 + 场景元素 + 中文文案 |
| **状态图标** | 蓝圈套色块 | 圆角方块背景 + 纯色渐变图标 |
| **Banner** | 纯色块+白星星 | 品牌渐变 + IP形象 + 文字标题 |
| **默认占位** | 灰色剪影/白底小图标 | 品牌色浅渐变 + 品牌点缀 |
| **Logo** | 像素化蓝底白字 | 矢量渲染 QH + 金色星星 |
| **Tab栏** | ✅ 已优化 | 品牌色 + 金色点缀 |

## 各分类说明

### tabbar/ (10个)
底部 5 个 Tab 的非/激活态：home, question, exam, class, my

### icons/ (28个)
arrow-right, like, liked, class, exam, profile, settings, wechat, auth, camera, creator, edit, delete, reply, favorite, about, statistics, clock, members, score, upload, creation, empty, eye, use, pending, rejected, success

### empty/ (7个)
class(班级为空), exam(考试为空), homework(作业为空), notice(通知为空), wrong(错题为空), empty(通用), no-data.svg(无数据) — 全部包含题小助IP形象表情

### error/ (3个)
404(页面未找到-挠头疑惑), 401(暂无权限-伸手挡门), network-error(网络连接失败-手持断连手机)

### status/ (4个)
icon-success(成功), icon-pending(进行中), icon-rejected(已拒绝), icon-check-success(完成勾选)

### defaults/ (3个)
default-avatar(默认头像), default-class(默认班级封面), default-class-cover(默认班级封面大图)

### banners/ (3个)
banner-1(题小助品牌), banner-2(刷题打卡), banner-3(班级挑战)
