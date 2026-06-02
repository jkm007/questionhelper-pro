# 移动端图片资源

> 来源：`questionhelper-app/static/` 及子目录

## 目录结构

```
mobile/
├── tabbar/          # 底部导航栏图标 (10个)
├── icons/           # 功能图标 (25个)
├── empty/           # 空状态图 (7个)
├── error/           # 错误页图 (3个)
├── status/          # 状态图标 (4个)
├── defaults/        # 默认占位图 (3个)
├── banners/         # Banner 图 (3个)
├── logo.png         # 移动端 Logo
└── splash-icon.png  # 启动页图标
```

## 各分类说明

### tabbar/ (10个)
底部 5 个 Tab 的未选中/选中态图标：home, question, exam, class, my
- 非激活态：灰色 `#B0B0B0` + 品牌色点缀
- 激活态：品牌色 `#4A90D9` + 金色点缀

### icons/ (25个)
功能操作图标：about, arrow-right, auth, camera, class, clock, creation, creator, delete, edit, empty, exam, eye, favorite, like, liked, members, pending, profile, rejected, reply, score, settings, statistics, success, upload, use, wechat

### empty/ (7个)
空状态图片：class(班级为空), exam(考试为空), homework(作业为空), notice(通知为空), wrong(错题为空), empty(通用), no-data.svg(无数据)

### error/ (3个)
错误页面：401(无权限), 404(页面未找到), network-error(网络错误)

### status/ (4个)
状态图标：check-success(成功), pending(进行中), rejected(已拒绝), success(成功标记)

### defaults/ (3个)
默认占位：default-avatar(默认头像), default-class(默认班级封面), default-class-cover(默认班级封面大图)

### banners/ (3个)
首页轮播图：banner-1, banner-2, banner-3

## 缺失 / 待补充

- [ ] **登录页背景图** — 移动端登录页需要品牌化背景
- [ ] **IP 表情包** — 品牌文档规划的 8 个表情：加油、学习、恭喜、思考、难过、加载、成就、通知
- [ ] **场景插画** — 引导页(3张)、首页装饰图等
- [ ] **加载动画** — SVG/CSS 动画（IP 转圈或翻书）
- [ ] **装饰元素** — 背景纹理、分隔线装饰、角标气泡
