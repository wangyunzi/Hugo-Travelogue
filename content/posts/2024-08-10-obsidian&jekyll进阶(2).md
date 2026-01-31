---
categories:
  - 博客
tags:
  - jekyll
date: "2024-08-10 05:42:40"
slug: "134"
title: obsidian&jekyll进阶(2)
thumb: 
backgrounds:
---

&emsp;&emsp;鉴于一篇只写了电脑端的如何新建文章发布博客，在移动端同样测试之后出现一点问题，于是小小修改一下代码，我的脚本文件夹`myob/assets/new_blog.js`：

```js
module.exports = async (params) => {
    const { app, quickAddApi } = params;

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

    // 如果用户取消输入标题，则终止脚本执行
    if (!userTitle) {
        new Notice('操作已取消，未创建新笔记。');
        return;
    }

    const currentDate = new Date().toISOString().split('T')[0];

    // 设置 YAML 头部内容，确保 date 字段始终用双引号包裹
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

    // 组合文件路径，使用字符串拼接代替 path 模块
    const filePath = `${obsidianFolder}/${fileName}`;

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

修改点总结：
1. 移除 path 模块：直接使用字符串拼接来生成 filePath，使其适应 Obsidian 环境。
2. 确保 YAML 内容正确格式化：代码生成的 YAML 头部确保在 Obsidian 中能够被正确识别。

演示视频如下：
<iframe src="//player.bilibili.com/player.html?isOutside=true&aid=112936148337403&bvid=BV1ndYVeUESs&cid=500001644542130&p=1" scrolling="no" border="0" frameborder="no" framespacing="0" allowfullscreen="true"></iframe>

&emsp;&emsp;此版本应能在 Obsidian 的 QuickAdd 插件中顺利运行，并创建带有 YAML 头部的新笔记文件。

&emsp;&emsp;如果单纯把文章放在文件夹里不见天日也不好，有时候隔太长时间不知道自己写了什么，所以利用`dataviewjs`获取最新的博客文章状态展示：
![](https://blog.wangyunzi.com/2024/08/525d9eac0660e8a38c3a249ee0cc76d5.png)

代码如下，需要利用dataview插件：
````
```dataviewjs
const folderPath = "日记本/posts"; // 替换为你的文件夹路径

let pages = dv.pages().filter(p => p.file.path.startsWith(folderPath)) // 获取指定文件夹下的文件
    .filter(p => p.file.name.match(/^\d{4}-\d{2}-\d{2}/)) // 过滤出符合日期格式的文件
    .sort(p => p.file.name.substring(0, 10), 'desc') // 按文件名中的日期排序
    .limit(3);

for (let page of pages) {
    let filePath = page.file.path.replace(/ /g, '%20'); // 替换空格为 %20
    dv.paragraph(`- [${page.file.name}](${filePath})`);
}
```
````

一个小bug就是dataviewjs代码块不能包含在callout代码块，要不然就会出现插件冲突问题，实在有点遗憾。