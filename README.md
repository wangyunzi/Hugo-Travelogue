长街短梦（Hugo 版）
===================

这是基于原 Travelogue 布局迁移到 Hugo 的版本，保留现有样式与页面结构。

快速开始
--------

- 本地预览：`hugo server -D`
- 生成静态文件：`hugo`

目录结构
--------

- `content/posts/`：文章
- `content/album/`：相册内容（相册首页为 `content/album/_index.md`）
- `content/`：独立页面（关于/友链/链圈/瞬间/时间线）
- `layouts/`：Hugo 模板与 partials
- `static/assets/`：原主题的 CSS/JS/图片资源
- `data/`：社交配置与 RSS 数据
- `rss/`：RSS 源列表（仅供抓取工具使用）
- `logs/`：抓取日志

说明
----

- RSS 输出文件名设为 `feed.xml`，用于站内搜索脚本读取。
- 菜单来源于页面 front matter 的 `menu.main.weight`。
- 文章默认开启阅读时间、标签/分类、分享按钮等配置，集中在 `content/posts/_index.md` 的 `cascade` 中。
- 原 Jekyll 相关文件已整理到 `legacy_jekyll/` 供备查。
- 主题来源：Travelogue（MIT License）。

维护脚本
--------

重命名文章文件（按 front matter 日期/标题）：

```
python scripts/rename_posts.py --posts-dir content/posts --dry-run
python scripts/rename_posts.py --posts-dir content/posts
```

评论数据转换（Twikoo 导出）：

```
python scripts/convert_comments.py --input /path/to/twikoo.json --output new_comments.json
```

RSS 抓取
--------

- 本地运行：`go run GrabLatestRSS/main.go`
- 源列表：`rss/rss_feeds.txt`
- 头像配置：`data/avatar_data.json`
- 输出：`data/rss_data.json`
- 日志：`logs/error.log`

目录与文件详解
-------------

### 根目录
- `.github/workflows/`：自动化工作流
  - `ScheduledRssRawler.yml`：定时抓取 RSS，更新 `data/rss_data.json`
  - `Sync_Obsidian.yml`：从 Obsidian 仓库同步内容到本仓库
  - `Notion_Sync.yml`：Notion → Markdown → 图床同步
  - `sync_trigger.yml`：仓库事件触发 Obsidian 同步
- `hugo.yaml`：Hugo 主配置（站点地址、分页、渲染、分类法等）
- `README.md`：项目说明与维护指南
- `CNAME`：GitHub Pages 自定义域名
- `.gitignore` / `.editorconfig`：Git 忽略与编辑器格式规范
- `.hugo_build.lock`：Hugo 构建锁文件（自动生成）
- `LICENSE`：主题/代码许可
- `new_comments.json`：评论数据（转换后的静态 JSON）

### 内容目录（content/）
- `content/posts/`：博客文章（Markdown）
- `content/album/`：相册内容（`_index.md` 为相册首页配置）
- `content/about.md` / `friend.md` / `links.md` / `archive.md` / `memos.md`：独立页面

### 模板目录（layouts/）
- `layouts/_default/`：通用布局
  - `baseof.html`：全站基础框架
  - `list.html`：列表页模板
  - `single.html`：通用单页模板
  - `_markup/render-image.html`：Markdown 图片渲染钩子
- `layouts/index.html`：首页布局
- `layouts/posts/single.html`：文章详情页
- `layouts/album/list.html` / `layouts/album/single.html`：相册列表/详情页
- `layouts/page/archive.html` / `layouts/page/links.html`：归档/链圈页
- `layouts/partials/`：组件/片段
  - `head.html`：页面 head（CSS 引用）
  - `js.html`：全站脚本
  - `sidebar.html`：侧边栏
  - `pagination.html` / `next.html`：分页/下一篇
  - `share.html`：分享按钮
  - `comments.html` / `twikoo.html` / `disqus.html`：评论相关
  - `search-lunr.html`：站内搜索
  - `loader.html`：加载动画

### 静态资源（static/）
- `static/assets/css/`：样式
  - `base.css` / `layout.css` / `sidebar.css` / `loader.css` / `syntax-highlighting.css`
  - `custom.css`：从模板内联样式收拢的自定义样式
- `static/assets/js/`：脚本
  - `functions.js`：主题交互逻辑
  - `lunr.js` / `lunr.zh.js`：搜索
  - `rss.js`：链圈相关逻辑（如无使用可考虑清理）
  - `memos.js`：瞬间/动态
  - `twikoo.all.min.js`：评论脚本
- `static/assets/js/libs/`：第三方库（jQuery/slidebars/rrssb 等）
- `static/assets/images/`：主题图片资源

### 数据与脚本
- `data/`：Hugo 数据源
  - `social.yml`：分享按钮开关
  - `rss_data.json`：RSS 抓取结果
  - `avatar_data.json`：RSS 头像配置
- `rss/`：RSS 源列表（`rss_feeds.txt`）
- `logs/`：RSS 抓取日志
- `scripts/`：维护脚本
  - `rename_posts.py`：按日期/标题重命名文章文件
  - `convert_comments.py`：Twikoo 评论数据转换

### RSS 抓取模块
- `GrabLatestRSS/`：主抓取程序（Go）
  - 读取：`rss/rss_feeds.txt`、`data/avatar_data.json`
  - 输出：`data/rss_data.json`、`logs/error.log`
- `api/_GrabLatestRSS/`：旧版/备用的抓取程序（未必仍在使用）

### 旧站备份与构建产物
- `legacy_jekyll/`：原 Jekyll 版本备份
- `public/`：Hugo 构建产物（可删，重新构建会生成）
