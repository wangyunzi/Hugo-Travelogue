---
categories:
  - 博客
tags:
  - jekyll
date: "2024-08-07 20:00:00"
slug: "131"
title: 使用GitHub Actions将Obsidian笔记同步到jekyll博客
---

## 介绍

在本教程中，我们将介绍如何使用 GitHub Actions 自动将 Obsidian 笔记同步到 GitHub 上的博客仓库中的特定文件夹。我们将使用 `rsync` 工具来同步文件，并利用 GitHub Actions 进行自动化处理。

## 前提条件

- GitHub 账户：你需要一个 GitHub 账户，并拥有两个仓库：一个用于 Obsidian 笔记（`wangyunzi/myob`），一个用于博客（`wangyunzi/wangyunzi.github.io`）。
- GitHub Personal Access Token (PAT)：你需要生成一个具有 `repo` 和 `workflow` 权限的 PAT。
- GitHub Actions：你需要在两个仓库中设置 GitHub Actions 工作流。

## 步骤

### （1）生成 Personal Access Token (PAT)

1. 登录到 GitHub，进入 **Settings**。
2. 选择 **Developer settings** > **Personal access tokens**。
3. 点击 **Generate new token**。
4. 输入名称，勾选 `repo` 和 `workflow` 权限。
5. 生成 token 并复制备用。

### （2）在博客仓库中设置 GitHub Actions

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
        token: ${{ secrets.PAT_TOKEN }}

    - name: Checkout blog repository
      uses: actions/checkout@v3
      with:
        repository: wangyunzi/wangyunzi.github.io
        path: wangyunzi.github.io
        token: ${{ secrets.PAT_TOKEN }}

    - name: List directories for debugging
      run: |
        echo "Listing myob"
        ls -R myob
        echo "Listing wangyunzi.github.io"
        ls -R wangyunzi.github.io

    - name: Sync _pages folder
      run: |
        rsync -av --delete myob/日记本/_pages/ wangyunzi.github.io/_pages/
        cd wangyunzi.github.io
        git config user.name "GitHub Actions"
        git config user.email "actions@github.com"
        git add _pages
        git commit -m "Sync _pages from Obsidian repository" || echo "No changes in _pages"
        cd -

    - name: Sync _album folder
      run: |
        rsync -av --delete myob/日记本/_album/ wangyunzi.github.io/_album/
        cd wangyunzi.github.io
        git config user.name "GitHub Actions"
        git config user.email "actions@github.com"
        git add _album
        git commit -m "Sync _album from Obsidian repository" || echo "No changes in _album"
        cd -

    - name: Sync _posts folder
      run: |
        rsync -av --delete myob/日记本/_posts/ wangyunzi.github.io/_posts/
        cd wangyunzi.github.io
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

4. 在博客仓库的 **Settings** > **Secrets** 中，添加 `PAT_TOKEN`，并将其设置为你生成的 Personal Access Token。

### （3）. 在 Obsidian 仓库中设置 GitHub Actions

1. 进入你的 Obsidian 仓库（`wangyunzi/myob`）。
2. 创建一个新的文件夹 `.github/workflows`。
3. 在该文件夹中创建一个名为 `Trigger_Sync.yml` 的文件，内容如下：

```yaml
name: Trigger Sync

on:
  push:
    branches:
      - master

jobs:
  trigger_sync:
    runs-on: ubuntu-latest

    steps:
    - name: Trigger Sync_Obsidian Workflow
      uses: actions/github-script@v6
      with:
        github-token: ${{ secrets.PAT_TOKEN }}
        script: |
          await github.rest.actions.createWorkflowDispatch({
            owner: 'wangyunzi',
            repo: 'wangyunzi.github.io',
            workflow_id: 'Sync_Obsidian.yml',
            ref: 'master' // 指定触发的分支
          });
```

4. 在 Obsidian 仓库的 **Settings** > **Secrets** 中，添加 `PAT_TOKEN`，并将其设置为你生成的 Personal Access Token。

### 4. 测试和验证

1. 向 Obsidian 仓库中的 `日记本/_pages`、`日记本/_album` 或 `日记本/_posts` 文件夹添加或修改文件。
2. 提交更改到 `master` 分支。
3. 检查 GitHub Actions 执行状态，确保同步流程没有错误。
4. 确认博客仓库中的 `_pages`、`_album` 和 `_posts` 文件夹已经更新。

## 总结

通过以上步骤，你可以将 Obsidian 笔记自动同步到你的博客仓库中。这样，你只需在 Obsidian 中编辑笔记，GitHub Actions 会自动将这些更改同步到博客中，保持内容的最新。

