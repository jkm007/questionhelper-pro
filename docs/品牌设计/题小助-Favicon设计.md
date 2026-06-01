# 题小助 Favicon 设计规范

**项目名称**：题小助 (QuestionHelper)  
**版本**: V1.0  
**日期**: 2026-05-29  
**作者**: UI设计师

---

## 1. Favicon 设计

### 1.1 设计理念

```
Favicon 设计理念：

┌─────────────────────────────────────┐
│                                     │
│   核心元素：QH                       │
│                                     │
│   Q = Question (题目)               │
│   H = Helper (助手)                 │
│                                     │
│   设计特点：                         │
│   ✅ 简洁：2个字母，容易识别         │
│   ✅ 圆润：友好、亲切                │
│   ✅ 蓝色：专业、可信赖              │
│   ✅ 可缩放：各种尺寸都清晰          │
│                                     │
└─────────────────────────────────────┘
```

### 1.2 Favicon 形态

```
标准 Favicon (32x32)：

    ┌─────────────┐
    │ ╭─────────╮ │
    │ │         │ │
    │ │   QH    │ │   蓝色背景
    │ │         │ │   白色文字
    │ ╰─────────╯ │   圆角矩形
    └─────────────┘

简化 Favicon (16x16)：

    ┌───────┐
    │ ╭───╮ │
    │ │QH │ │   更简洁
    │ ╰───╯ │   保持可读性
    └───────┘
```

---

## 2. Favicon 尺寸规范

### 2.1 标准尺寸

| 尺寸 | 用途 | 格式 | 说明 |
|------|------|------|------|
| **16x16** | 浏览器标签 | ICO/PNG | 最小尺寸 |
| **32x32** | 浏览器标签 | ICO/PNG | 标准尺寸 |
| **48x48** | Windows快捷方式 | PNG | Windows |
| **64x64** | 高DPI屏幕 | PNG | Retina |
| **128x128** | Chrome Web Store | PNG | 应用商店 |
| **180x180** | Apple Touch Icon | PNG | iOS |
| **192x192** | Android Chrome | PNG | Android |
| **512x512** | PWA Splash | PNG | 启动画面 |

### 2.2 文件命名

```
favicon 文件结构：

assets/
├── favicon/
│   ├── favicon.ico              # 16x16 + 32x32 (多尺寸ICO)
│   ├── favicon-16x16.png        # 16x16 PNG
│   ├── favicon-32x32.png        # 32x32 PNG
│   ├── favicon-48x48.png        # 48x48 PNG
│   ├── favicon-64x64.png        # 64x64 PNG
│   ├── favicon-128x128.png      # 128x128 PNG
│   ├── apple-touch-icon.png     # 180x180 Apple
│   ├── android-chrome-192x192.png # 192x192 Android
│   ├── android-chrome-512x512.png # 512x512 Android
│   └── mstile-150x150.png       # 150x150 Windows
```

---

## 3. Favicon 设计稿

### 3.1 标准版 (推荐)

```
┌─────────────────────────────────────┐
│                                     │
│   ████████████████████████████████  │
│   █                              █  │
│   █    ╭─────────────────────╮   █  │
│   █    │                     │   █  │
│   █    │         QH          │   █  │
│   █    │                     │   █  │
│   █    ╰─────────────────────╯   █  │
│   █                              █  │
│   ████████████████████████████████  │
│                                     │
│   背景：#4A90D9 (助手蓝)            │
│   文字：白色                        │
│   圆角：20%                         │
│   字体：SF Pro Display Bold         │
│                                     │
└─────────────────────────────────────┘
```

### 3.2 简化版 (小尺寸)

```
┌─────────────────────────────────────┐
│                                     │
│   16x16 和 32x32 使用简化版：       │
│                                     │
│   ┌───────┐    ┌─────────────┐      │
│   │ ╭───╮ │    │ ╭─────────╮ │      │
│   │ │QH │ │    │ │   QH    │ │      │
│   │ ╰───╯ │    │ ╰─────────╯ │      │
│   └───────┘    └─────────────┘      │
│   16x16         32x32               │
│                                     │
│   特点：                            │
│   - 更粗的笔画                      │
│   - 更大的内边距                    │
│   - 简化细节                        │
│                                     │
└─────────────────────────────────────┘
```

### 3.3 圆形版 (备选)

```
┌─────────────────────────────────────┐
│                                     │
│   圆形版本：                        │
│                                     │
│       ╭─────────────╮               │
│       │             │               │
│       │     QH      │               │
│       │             │               │
│       ╰─────────────╯               │
│                                     │
│   适用场景：                        │
│   - App图标                         │
│   - 社交媒体头像                    │
│   - 需要圆形适配的场景              │
│                                     │
└─────────────────────────────────────┘
```

---

## 4. HTML 引用代码

### 4.1 基础引用

```html
<!-- 基础 Favicon -->
<link rel="icon" type="image/x-icon" href="/favicon.ico">
<link rel="icon" type="image/png" sizes="32x32" href="/favicon-32x32.png">
<link rel="icon" type="image/png" sizes="16x16" href="/favicon-16x16.png">
```

### 4.2 完整引用 (推荐)

```html
<!-- Favicon -->
<link rel="icon" type="image/x-icon" href="/favicon.ico">
<link rel="icon" type="image/png" sizes="32x32" href="/favicon-32x32.png">
<link rel="icon" type="image/png" sizes="16x16" href="/favicon-16x16.png">

<!-- Apple -->
<link rel="apple-touch-icon" sizes="180x180" href="/apple-touch-icon.png">
<link rel="apple-touch-icon" sizes="152x152" href="/apple-touch-icon-152x152.png">
<link rel="apple-touch-icon" sizes="144x144" href="/apple-touch-icon-144x144.png">
<link rel="apple-touch-icon" sizes="120x120" href="/apple-touch-icon-120x120.png">
<link rel="apple-touch-icon" sizes="114x114" href="/apple-touch-icon-114x114.png">
<link rel="apple-touch-icon" sizes="76x76" href="/apple-touch-icon-76x76.png">
<link rel="apple-touch-icon" sizes="72x72" href="/apple-touch-icon-72x72.png">
<link rel="apple-touch-icon" sizes="60x60" href="/apple-touch-icon-60x60.png">
<link rel="apple-touch-icon" sizes="57x57" href="/apple-touch-icon-57x57.png">

<!-- Android -->
<link rel="icon" type="image/png" sizes="192x192" href="/android-chrome-192x192.png">
<link rel="icon" type="image/png" sizes="512x512" href="/android-chrome-512x512.png">

<!-- Microsoft -->
<meta name="msapplication-TileColor" content="#4A90D9">
<meta name="msapplication-TileImage" content="/mstile-150x150.png">

<!-- PWA -->
<link rel="manifest" href="/site.webmanifest">

<!-- Safari Pinned Tab -->
<link rel="mask-icon" href="/safari-pinned-tab.svg" color="#4A90D9">
```

### 4.3 Web App Manifest

```json
// site.webmanifest
{
  "name": "题小助 - 学习好帮手",
  "short_name": "题小助",
  "description": "面向大众的社交学习平台",
  "icons": [
    {
      "src": "/android-chrome-192x192.png",
      "sizes": "192x192",
      "type": "image/png"
    },
    {
      "src": "/android-chrome-512x512.png",
      "sizes": "512x512",
      "type": "image/png"
    }
  ],
  "theme_color": "#4A90D9",
  "background_color": "#ffffff",
  "display": "standalone",
  "start_url": "/"
}
```

---

## 5. 设计稿 SVG 代码

### 5.1 Favicon SVG

```svg
<!-- favicon.svg -->
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 32 32" width="32" height="32">
  <!-- 背景 -->
  <rect width="32" height="32" rx="6" ry="6" fill="#4A90D9"/>

  <!-- Q 字母 -->
  <circle cx="11" cy="12" r="5" fill="none" stroke="white" stroke-width="2.5"/>
  <line x1="14" y1="15" x2="17" y2="18" stroke="white" stroke-width="2.5" stroke-linecap="round"/>

  <!-- H 字母 -->
  <g transform="translate(22, 22)">
    <line x1="-3" y1="-4" x2="-3" y2="4" stroke="white" stroke-width="2" stroke-linecap="round"/>
    <line x1="3" y1="-4" x2="3" y2="4" stroke="white" stroke-width="2" stroke-linecap="round"/>
    <line x1="-3" y1="0" x2="3" y2="0" stroke="white" stroke-width="2" stroke-linecap="round"/>
  </g>
</svg>
```

### 5.2 简化 SVG (16x16)

```svg
<!-- favicon-16x16.svg -->
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 16 16" width="16" height="16">
  <!-- 背景 -->
  <rect width="16" height="16" rx="3" ry="3" fill="#4A90D9"/>

  <!-- Q 字母 -->
  <circle cx="5" cy="6" r="2.5" fill="none" stroke="white" stroke-width="1.5"/>
  <line x1="7" y1="8" x2="9" y2="10" stroke="white" stroke-width="1.5" stroke-linecap="round"/>

  <!-- H 字母 -->
  <g transform="translate(11, 11)">
    <line x1="-1.5" y1="-2" x2="-1.5" y2="2" stroke="white" stroke-width="1.2" stroke-linecap="round"/>
    <line x1="1.5" y1="-2" x2="1.5" y2="2" stroke="white" stroke-width="1.2" stroke-linecap="round"/>
    <line x1="-1.5" y1="0" x2="1.5" y2="0" stroke="white" stroke-width="1.2" stroke-linecap="round"/>
  </g>
</svg>
```

---

## 6. 设计工具生成指南

### 6.1 使用 Figma 生成

```
Figma 操作步骤：

1. 创建 512x512 画布
2. 绘制圆角矩形 (96px圆角)
3. 填充渐变色 (#4A90D9 → #357ABD)
4. 添加 QH 文字 (白色, Bold)
5. 导出为各种尺寸 PNG

导出设置：
- 格式：PNG
- 缩放：1x, 2x, 3x
- 命名：favicon-{size}x{size}.png
```

### 6.2 使用在线工具

```
推荐工具：

1. Favicon.io
   - 网址：https://favicon.io/
   - 功能：从文字/图片生成Favicon
   - 输出：ICO + PNG + HTML代码

2. RealFaviconGenerator
   - 网址：https://realfavicongenerator.net/
   - 功能：全平台Favicon生成
   - 输出：所有尺寸 + 配置代码

3. Favicon Generator
   - 网址：https://www.favicon-generator.org/
   - 功能：批量生成Favicon
   - 输出：多尺寸PNG + ICO
```

---

## 7. 使用场景

### 7.1 浏览器

```
浏览器标签页效果：

┌─────────────────────────────────────┐
│ [QH] 题小助 - 学习好帮手           │
├─────────────────────────────────────┤
│                                     │
│         页面内容...                 │
│                                     │
└─────────────────────────────────────┘

16x16 或 32x32 Favicon
```

### 7.2 移动设备

```
手机桌面图标：

    ┌─────────────┐
    │             │
    │    ╭───╮    │
    │    │QH │    │   180x180 或 192x192
    │    ╰───╯    │   圆角处理
    │             │
    │   题小助    │
    └─────────────┘

iOS: apple-touch-icon.png
Android: android-chrome-192x192.png
```

### 7.3 社交媒体

```
社交媒体头像：

    ╭─────────────╮
    │             │
    │     QH      │   400x400 或更大
    │             │   圆形裁剪
    ╰─────────────╯

微信、微博、抖音等
```

---

## 附录

### A. 文件清单

| 文件名 | 尺寸 | 格式 | 用途 |
|--------|------|------|------|
| favicon.ico | 16x16, 32x32 | ICO | 浏览器标签 |
| favicon-16x16.png | 16x16 | PNG | 浏览器标签 |
| favicon-32x32.png | 32x32 | PNG | 浏览器标签 |
| favicon-48x48.png | 48x48 | PNG | Windows |
| favicon-64x64.png | 64x64 | PNG | Retina |
| favicon-128x128.png | 128x128 | PNG | Chrome |
| apple-touch-icon.png | 180x180 | PNG | iOS |
| android-chrome-192x192.png | 192x192 | PNG | Android |
| android-chrome-512x512.png | 512x512 | PNG | PWA |
| mstile-150x150.png | 150x150 | PNG | Windows |

### B. 色值参考

| 颜色 | 色值 | 用途 |
|------|------|------|
| 主色 | #4A90D9 | 背景 |
| 主色-深 | #357ABD | 渐变 |
| 白色 | #FFFFFF | 文字 |
| 金色 | #FFD700 | 装饰 |

---

**文档结束**
