---
categories: 博客
tags: ["学习"]
pstatus: 已发布
ptype: posts
date: 2024-06-07 23:34:00 +0800
slug: "125"
id: 81e6f110-1ea2-407c-9799-c78f18f1e3b9
title: Netlify部署hugo博客失败怎么办？
created_time: 2024-06-07T15:34:00.000Z
cover: 
icon: 
last_edited_time: 2024-06-07T15:56:00.000Z
archived: false
created_by_object: user
created_by_id: 82c1f8d9-63cb-4b06-aa84-69c69af9ea5d
last_edited_by_object: user
last_edited_by_id: 82c1f8d9-63cb-4b06-aa84-69c69af9ea5d
---

起初是因为自己修改页面的时候发现GitHub+Netlify部署的时候出现一个部署成功一个失败，道理来讲一个博客每次修改应该只需要部署一次，但是不知道为什么netlify部署了两次，且名称都显示不一样，错误如下图所示：

![](https://blog.wangyunzi.com/2024/06/OmdQZN.jpg)



经过查询得知应该是根目录的`netlify.toml`文件出现问题，于是询问AI后得知这个错误表明在构建过程中无法找到指定的基本目录（base directory）`/opt/build/repo/blog` ，于是修改方案如下，将原来的netlify.toml代码：

```go
[build]
publish = "public"
command = "hugo --gc --minify"
```

修改成下面：

```go
[build]
base = ""
publish = "public"
command = "hugo --gc --minify"
```

最后显示部署成功。


