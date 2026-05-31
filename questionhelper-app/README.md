# 题小助 V2.0 移动端

基于 uni-app 的跨平台移动应用，支持 H5、微信小程序、App。

## 技术栈

- **框架**: uni-app (Vue 3 + TypeScript)
- **状态管理**: Pinia
- **UI**: uni-app 内置组件
- **图表**: ucharts
- **日期处理**: dayjs

## 项目结构

```
questionhelper-app/
├── pages/              # 页面
├── components/         # 公共组件
├── store/              # 状态管理
├── api/                # API 接口
├── utils/              # 工具函数
├── hooks/              # 自定义 Hooks
├── static/             # 静态资源
├── pages.json          # 页面配置
├── manifest.json       # 应用配置
├── uni.scss            # 全局样式变量
├── App.vue             # 根组件
└── main.ts             # 入口文件
```

## 快速开始

### 安装依赖

```bash
npm install
```

### 开发

```bash
# H5
npm run dev:h5

# 微信小程序
npm run dev:mp-weixin

# App
npm run dev:app
```

### 构建

```bash
# H5
npm run build:h5

# 微信小程序
npm run build:mp-weixin

# App
npm run build:app
```

## 品牌色

| 颜色 | 色值 | 用途 |
|------|------|------|
| 主色 | #4A90D9 | 按钮、链接、强调 |
| 成功 | #67C23A | 成功状态 |
| 警告 | #E6A23C | 警告状态 |
| 错误 | #F56C6C | 错误状态 |

## 页面说明

| 页面 | 路径 | 说明 |
|------|------|------|
| 首页 | pages/index/index | 数据概览、快捷入口 |
| 登录 | pages/login/index | 用户登录 |
| 题库 | pages/question/list | 题目列表 |
| 练习 | pages/practice/index | 练习入口 |
| 考试 | pages/exam/list | 考试列表 |
| 班级 | pages/class/list | 班级列表 |
| 我的 | pages/profile/index | 个人中心 |

## License

MIT
