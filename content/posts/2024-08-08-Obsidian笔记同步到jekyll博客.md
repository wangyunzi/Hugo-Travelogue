---
categories:
  - 博客
tags:
  - jekyll
date: "2024-08-08 08:18:00"
slug: "132"
title: Obsidian笔记同步到jekyll博客
---

&emsp;&emsp;本来这篇文章应该在昨天就写好，但是鉴于自己流程写的有点乱加上最后非要请AI 帮我折腾一下就导致出现问题了，然后自己重新慢慢开始部署一番，终于才算把流程弄好。

### 前提条件
- GitHub 账户：你需要一个 GitHub 账户，并拥有两个仓库：一个用于 Obsidian 笔记（`wangyunzi/myob`），一个用于博客（`wangyunzi/wangyunzi.github.io`）。
- GitHub Personal Access Token (PAT)：你需要生成三个具有 `repo` 和 `workflow` 权限的 PAT，复制到备忘录中以备后用。
- GitHub Actions：你需要在两个仓库中设置 GitHub Actions 工作流。

### 在博客仓库中设置 GitHub Actions
1. 进入你的博客仓库（`wangyunzi/wangyunzi.github.io`）。
2. 创建一个新的文件夹 `.github/workflows`。
3. 在该文件夹中创建一个名为 `Sync_Obsidian.yml` 的文件，内容如下：

```yaml
name: Sync_Obsidian

on:
  workflow_dispatch:
    inputs:
      ref:
        description: 'Branch to deploy'
        required: true
        default: 'master'

jobs:
  sync:
    runs-on: ubuntu-latest

    steps:
    - name: Checkout Obsidian repository
      uses: actions/checkout@v3
      with:
        repository: wangyunzi/myob
        path: myob
        token: ${{ secrets.OBSIDIAN_REPO_TOKEN }}

    - name: Checkout blog repository
      uses: actions/checkout@v3
      with:
        repository: wangyunzi/wangyunzi.github.io
        path: wangyunzi.github.io
        token: ${{ secrets.BLOG_REPO_TOKEN }}

    - name: List directories for debugging
      run: |
        echo "Listing myob"
        ls -R myob
        echo "Listing wangyunzi.github.io"
        ls -R wangyunzi.github.io

    - name: Sync _pages folder
      run: |
        rsync -av --delete myob/日记本/pages/ wangyunzi.github.io/_pages/
        cd wangyunzi.github.io
        find _pages -name "*.md" -exec sed -i -E '
          # ISO 8601格式（例如 2024-08-07T16:00:00）
          s/date: ([0-9]{4}-[0-9]{2}-[0-9]{2})T([0-9]{2}:[0-9]{2}:[0-9]{2})([+-][0-9]{2}:[0-9]{2}|)/date: \1 \2 +0800/
          # 带有时区的格式（例如 2024-08-07T16:00:00+08:00）
          s/date: ([0-9]{4}-[0-9]{2}-[0-9]{2})T([0-9]{2}:[0-9]{2}:[0-9]{2})[+-][0-9]{2}:[0-9]{2}/date: \1 \2 +0800/
          # 不带时区的格式（例如 2024-08-07 16:00:00）
          s/date: ([0-9]{4}-[0-9]{2}-[0-9]{2}) ([0-9]{2}:[0-9]{2}:[0-9]{2})/date: \1 \2 +0800/
        ' {} \;
        git config user.name "GitHub Actions"
        git config user.email "actions@github.com"
        git add _pages
        git commit -m "Sync _pages from Obsidian repository" || echo "No changes in _pages"
        cd -

    - name: Sync _album folder
      run: |
        rsync -av --delete myob/日记本/album/ wangyunzi.github.io/_album/
        cd wangyunzi.github.io
        find _album -name "*.md" -exec sed -i -E '
          # ISO 8601格式（例如 2024-08-07T16:00:00）
          s/date: ([0-9]{4}-[0-9]{2}-[0-9]{2})T([0-9]{2}:[0-9]{2}:[0-9]{2})([+-][0-9]{2}:[0-9]{2}|)/date: \1 \2 +0800/
          # 带有时区的格式（例如 2024-08-07T16:00:00+08:00）
          s/date: ([0-9]{4}-[0-9]{2}-[0-9]{2})T([0-9]{2}:[0-9]{2}:[0-9]{2})[+-][0-9]{2}:[0-9]{2}/date: \1 \2 +0800/
          # 不带时区的格式（例如 2024-08-07 16:00:00）
          s/date: ([0-9]{4}-[0-9]{2}-[0-9]{2}) ([0-9]{2}:[0-9]{2}:[0-9]{2})/date: \1 \2 +0800/
        ' {} \;
        git config user.name "GitHub Actions"
        git config user.email "actions@github.com"
        git add _album
        git commit -m "Sync _album from Obsidian repository" || echo "No changes in _album"
        cd -

    - name: Sync _posts folder
      run: |
        rsync -av --delete myob/日记本/posts/ wangyunzi.github.io/_posts/
        cd wangyunzi.github.io
        find _posts -name "*.md" -exec sed -i -E '
          # ISO 8601格式（例如 2024-08-07T16:00:00）
          s/date: ([0-9]{4}-[0-9]{2}-[0-9]{2})T([0-9]{2}:[0-9]{2}:[0-9]{2})([+-][0-9]{2}:[0-9]{2}|)/date: \1 \2 +0800/
          # 带有时区的格式（例如 2024-08-07T16:00:00+08:00）
          s/date: ([0-9]{4}-[0-9]{2}-[0-9]{2})T([0-9]{2}:[0-9]{2}:[0-9]{2})[+-][0-9]{2}:[0-9]{2}/date: \1 \2 +0800/
          # 不带时区的格式（例如 2024-08-07 16:00:00）
          s/date: ([0-9]{4}-[0-9]{2}-[0-9]{2}) ([0-9]{2}:[0-9]{2}:[0-9]{2})/date: \1 \2 +0800/
        ' {} \;
        git config user.name "GitHub Actions"
        git config user.email "actions@github.com"
        git add _posts
        git commit -m "Sync _posts from Obsidian repository" || echo "No changes in _posts"
        cd -

    - name: Push changes
      run: |
        cd wangyunzi.github.io
        git push
```
文件说明：这个文件的作用就是把你想要从obsidian里面同步的特定文件同步至博客仓库，`myob/日记本/posts/` 这样的代表obsidian里面的文件夹，`wangyunzi.github.io/_posts/`代表博客里面的文件夹，可以同步多个文件夹，需保证文件存在且路径正确，否则会报错。鉴于昨天测试的时候发现jekyll出现文章时间问题，也就是obsidian同步过去的文章时间不是在东八区，而是在东七区，我也不知道怎么回事，导致jekyll的文章时间在未来时刻，也就是不会生成文章，建议将格式写成 `date: "2024-08-08 08:18:00"` 或者 `date: 2024-08-08 08:18:00 +0800 +0800` , 这样需要把`Sync_Obsidian.yml`文件里面转换时间的格式去掉，否则就会在同步的时候检测时间格式，并且转换成加上时区的格式。`OBSIDIAN_REPO_TOKEN`和 `BLOG_REPO_TOKEN`需要勾选正确的权限，以及设置在`wangyunzi.github.io`仓库中。

4. 继续在该文件夹下面创建一个名为 `Sync_Trigger.yml` 的文件，内容如下：

```yaml
name: Trigger Sync Obsidian

on:
  repository_dispatch:
    types: [sync_from_obsidian]

jobs:
  trigger:
    runs-on: ubuntu-latest

    steps:
    - name: Trigger Sync Obsidian workflow
      uses: actions/github-script@v6
      with:
        github-token: ${{ secrets.Trigger_TOKEN }}
        script: |
          await github.rest.actions.createWorkflowDispatch({
            owner: 'wangyunzi',
            repo: 'wangyunzi.github.io',
            workflow_id: 'Sync_Obsidian.yml', // 确保这是你的工作流文件名
            ref: 'master' // 确保这是你的目标分支
          });
```
文件说明：这个文件夹的作用是为了触发上面一个文件夹。

### 在 Obsidian 仓库中设置 GitHub Actions

1. 进入你的 Obsidian 仓库（`wangyunzi/myob`）。
2. 创建一个新的文件夹 `.github/workflows`。
3. 在该文件夹中创建一个名为 `Trigger_Sync.yml` 的文件，内容如下：

```yaml
name: Notify Sync_Obsidian

on:
  push:
    branches:
      - master  # 或者其他你想监听的分支

jobs:
  notify:
    runs-on: ubuntu-latest

    steps:
    - name: Trigger Sync_Obsidian Workflow
      uses: actions/github-script@v6
      with:
        github-token: ${{ secrets.Trigger_TOKEN }}
        script: |
          await github.rest.repos.createDispatchEvent({
            owner: 'wangyunzi',
            repo: 'wangyunzi.github.io',
            event_type: 'sync_from_obsidian'
          });

```

文件说明：当obsidian出现push动作的时候就会触发这个文件，继而触发上面一个action动作。

最后，自动化obsidian写作就完成了，在obsidian中写好文章或者其他页面，然后git push到github仓库，就会触发后面的实现发布博客。

ps：最后发布完了之后发现怎么jekyll连代码都显示不正确啊，第一个文件直接给我少了东西，文件里面还是好好的。