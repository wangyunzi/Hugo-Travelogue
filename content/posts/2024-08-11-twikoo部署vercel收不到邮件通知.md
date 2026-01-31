---
categories:
  - 博客
tags:
  - twikoo
date: "2024-08-11 09:27:02"
slug: "137"
title: twikoo部署vercel收不到邮件通知
thumb: 
backgrounds:
---

&emsp;&emsp;前几天就一直在选择合适的评论系统，换来换去最终归于twikoo，本来以前自己成功部署过twikoo但是因为手贱删除部署信息了，所以重新设置一遍，结果最后无论怎么样设置邮件通知功能都是只能收到测试成功的邮件，收不到其他邮件通知，甚是苦恼，想着换一个，但是又陷入纠结中。

&emsp;&emsp;最后都已经删除了twikoo的相关文件了，上网搜索了之后发现了一个最佳答案怀着激动的心有重新设置一遍，终于成功，感谢这位大佬发布的解惑。

&emsp;&emsp;解决方法：进入 Settings - Deployment Protection，设置 Vercel Authentication 为 Disabled，并 Save。

相关链接：
- [vercel部署twikoo后评论收不到通知邮件问题解决方法](https://leolin66.top/article/twikoo_email_problem)
- [配置完成，测试邮件收发正常，但是正常评论的时候即发不出去也收不到提醒 #625](https://github.com/twikoojs/twikoo/issues/625)
- [System Environment Variables Overview](https://vercel.com/docs/projects/environment-variables/system-environment-variables)

ps：这个twikoo是不是有点网络问题，我最开始链接家里的移动WiFi登录不上，而且导入的评论什么的也没有了，换成自己手机热点一下就可以了，真真是怪哉！最终发现，额，如果是美国ip可能就连接不上而且看不到评论。