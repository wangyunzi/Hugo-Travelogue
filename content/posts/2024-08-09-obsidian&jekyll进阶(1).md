---
categories:
  - 博客
tags:
  - jekyll
date: "2024-08-09 11:29:07"
slug: "133"
title: obsidian&jekyll进阶(1)
thumb: 
backgrounds:
---

&emsp;&emsp;在上一篇自动化利用obsidian写jekyll博客之后，日常写博客又会出现一个问题，就是左侧文件列表栏的文件太多，每次都需要特意打开posts文件夹然后滑到最下面新建一个博客文章，虽然有了`Templater`插件，但是这个插件有点过于死板了，于是想到了`QuickAdd`插件，通过运行一个`JavaScript`脚本新建文章并且自动填充好`yaml`，包括date、slug以及title等日常必须，如果`slug`还是需要自己填写的话就未必有点得不偿失了，有可能自己半个月不写博客导致忘记上一篇`slug`数据，因此每次都要自己重新点开上一篇的数据信息查看，脚本最好能后遍历自己的所有文章slug然后选取最大的数字最后只需要在新建文章＋1即可。

&emsp;&emsp;具体操作：需要安装`obsidian-quickadd`插件，添加宏步骤中的脚本动作即可。脚本信息如下：

```
module.exports = async (params) => {
  const { app, quickAddApi } = params;
  const path = require('path');

  // 设置你的 Obsidian 文件夹路径
  const obsidianFolder = '日记本/posts'; // 替换为实际路径

  // 输出调试信息
  console.log('Obsidian folder path:', obsidianFolder);

  let mdFiles;
  try {
      mdFiles = await app.vault.adapter.list(obsidianFolder);
      console.log('Files in folder:', mdFiles.files);
  } catch (error) {
      console.error('Error reading directory:', error);
      new Notice('Failed to read Obsidian folder. Please check the folder path.');
      return;
  }

  const markdownFiles = mdFiles.files.filter(file => file.endsWith('.md'));

  // 初始化最大 slug
  let maxSlug = 0;

  // 正则表达式匹配 slug
  const slugPattern = /slug:\s*"(\d+)"/;

  // 遍历文件以找到最大 slug
  for (const file of markdownFiles) {
      const content = await app.vault.adapter.read(file);
      const match = content.match(slugPattern);
      if (match) {
          maxSlug = Math.max(maxSlug, parseInt(match[1], 10));
      }
  }

  // 新文章的 slug 为最大 slug + 1
  const newSlug = maxSlug + 1;

  // 让用户输入标题
  const userTitle = await quickAddApi.inputPrompt("请输入标题");
  const currentDate = new Date().toISOString().split('T')[0];

  // 设置 YAML 头部内容，date 字段添加双引号，title 仅包括用户输入的标题
  const yamlContent = `---
categories:
- 
tags:
- 
date: "${new Date().toISOString().replace('T', ' ').split('.')[0]}"
slug: "${newSlug}"
title: "${userTitle}"
thumb: 
backgrounds: 
---
`;

  // 获取新笔记的文件名，以“yyyy-mm-dd-标题”的格式命名
  const fileName = `${currentDate}-${userTitle}.md`;

  // 创建新笔记
  const filePath = path.join(obsidianFolder, fileName);

  try {
      await app.vault.create(filePath, yamlContent + '\n# Your Note Content\n');
      new Notice(`New note created at: ${filePath}`);

      // 在新标签中打开新创建的笔记
      const createdFile = app.vault.getAbstractFileByPath(filePath);
      app.workspace.getLeaf(true).openFile(createdFile);
  } catch (error) {
      console.error('Error creating note:', error);
      new Notice('Failed to create note. Please check the folder path and permissions.');
  }
};

```

&emsp;&emsp;演示视频如下：

<iframe src="//player.bilibili.com/player.html?isOutside=true&aid=112931903767041&bvid=BV18gYgeXEdQ&cid=500001643783998&p=1" scrolling="no" border="0" frameborder="no" framespacing="0" allowfullscreen="true"></iframe>

