

## 项目特点

- gpt3
- 飞书机器人
- 责任链的设计模式
- zap日志记录
- goCache缓存

## 项目介绍

聊天机器人，当然得在聊天软件上使用呀！

## 功能解释

### 责任链-设计模式

划重点@bro

千万不要用if else，这样的代码，不仅可读性差，而且，如果要增加一个处理器，就需要修改代码，违反了开闭原则

用户发送的文本消息，根据消息内容，匹配到对应的处理器，处理器处理消息，返回结果给用户

这种匹配，可以使用责任链模式，将匹配的逻辑抽象成一个个的处理器，然后将这些处理器串联起来，形成一个链条。

用户发送的消息，从链条的头部开始，依次匹配，匹配到后，就不再继续匹配，直接返回结果给用户


！！！切记！！！

责任链模式[参考代码](https://refactoringguru.cn/design-patterns/chain-of-responsibility)



### 日志记录

- 按照文件大小切割


### 相关阅读

- [在Go语言项目中使用Zap日志库](https://www.liwenzhou.com/posts/Go/zap/)
- 
- [飞书 User_ID、Open_ID 与 Union_ID 区别](https://www.feishu.cn/hc/zh-CN/articles/794300086214)
- 
