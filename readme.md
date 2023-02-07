
<p align='center'>
  <img src='https://p6-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/da7632aafca8419d93ace47fa6479feb~tplv-k3u1fbpfcp-watermark.image?' alt='' width='900'/>
</p>


<p align='center'>
在飞书与ChatGPT随时对话，智慧随身。
<br>
 Feishu ChatGpt
</p>

## 项目特点

- 🍏 openai [gpt3](https://platform.openai.com/account/api-keys)
- 🥒 [serverless一键部署](https://github.com/serverless-devs/serverless-devs)
- 🍎 [飞书](https://open.feishu.cn/app)机器人
- 🍐 [责任链](https://refactoringguru.cn/design-patterns/chain-of-responsibility/go/example)的设计模式
- 🍊 [zap](https://github.com/uber-go/zap)日志记录
- 🍋 [goCache](https://github.com/patrickmn/go-cache)内存键值对缓存


## 部署
``` bash
git clone git@github.com:Leizhenpeng/feishu-chatGpt.git
cd feishu-chatGpt/code

# 配置config.yaml
mv config.example.yaml config.yaml

# serverless部署
cd ..
s deploy
```

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

- [go-cache](https://github.com/patrickmn/go-cache)

- [在Go语言项目中使用Zap日志库](https://www.liwenzhou.com/posts/Go/zap/)

- [飞书 User_ID、Open_ID 与 Union_ID 区别](https://www.feishu.cn/hc/zh-CN/articles/794300086214)

- [飞书重复接受到消息](https://open.feishu.cn/document/uAjLw4CM/ukTMukTMukTM/reference/im-v1/message/events/receive)
