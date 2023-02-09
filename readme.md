
<p align='center'>
  <img src='https://p6-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/91d1c7af087646aea2c550665c01796b~tplv-k3u1fbpfcp-watermark.image?' alt='' width='900'/>
</p>

<br>

<p align='center'>
    在飞书与ChatGPT随时对话，智慧随身。
    <br>
    Feishu ChatGpt
</p>

## 项目特点
- 🍏 openai-[gpt3](https://platform.openai.com/account/api-keys)
- 🍎 [飞书](https://open.feishu.cn/app)机器人
- 🥒 支持[Serverless](https://github.com/serverless-devs/serverless-devs)、[本地环境](https://dashboard.cpolar.com/login)、[Docker](https://www.docker.com/) 多渠道部署
- 🍐 基于[责任链](https://refactoringguru.cn/design-patterns/chain-of-responsibility/go/example)的消息处理器，轻松自定义扩展命令

[//]: # (- 🍊 [zap]&#40;https://github.com/uber-go/zap&#41;日志记录)

[//]: # (- )
- 🍋 基于[goCache](https://github.com/patrickmn/go-cache)内存键值对缓存


## 项目部署


######  有关飞书具体的配置文件说明，**[➡︎ 点击查看](#详细配置步骤)**


``` bash
git clone git@github.com:Leizhenpeng/feishu-chatGpt.git
cd feishu-chatGpt/code

# 配置config.yaml
mv config.example.yaml config.yaml
```
<details>
    <summary>本地部署</summary>
    <br>

如果你的服务器没有公网 IP，可以使用反向代理的方式

飞书的服务器在国内对ngrok的访问速度很慢，所以推荐使用一些国内的反向代理服务商
- [cpolar](https://dashboard.cpolar.com/)
- [natapp](https://natapp.cn/)


```bash
//测试部署
go run main.go
cpolar http 9000

//正式部署
nohup cpolar http 8080 -log=stdout &

//查看服务器状态
https://dashboard.cpolar.com/status

// 下线服务
ps -ef | grep cpolar
kill -9 PID
```

更多详细介绍，参考[飞书上的小计算器: Go机器人来啦](https://www.bilibili.com/video/BV1nW4y1378T/)

    <br>

</details>


<details>
    <summary>serverless部署</summary>
<br>

``` bash
cd ..
s deploy
```

更多详细介绍，参考[仅需1min，用Serverless部署基于 gin 的飞书机器人](https://www.bilibili.com/video/BV1nW4y1378T/)
    <br>

</details>


<details>
    <summary>docker部署</summary>
    <br>

待补充
    <br>

</details>


## 功能解释

### 责任链-设计模式

划重点@bro

千万不要用if else，这样的代码，不仅可读性差，而且，如果要增加一个处理器，就需要修改代码，违反了开闭原则

用户发送的文本消息，根据消息内容，匹配到对应的处理器，处理器处理消息，返回结果给用户

这种匹配，可以使用责任链模式，将匹配的逻辑抽象成一个个的处理器，然后将这些处理器串联起来，形成一个链条。

用户发送的消息，从链条的头部开始，依次匹配，匹配到后，就不再继续匹配，直接返回结果给用户


！！！切记！！！

责任链模式[参考代码](https://refactoringguru.cn/design-patterns/chain-of-responsibility)



## 详细配置步骤

-  获取 [OpenAI](https://platform.openai.com/account/api-keys) 的 KEY
-  创建 [飞书](https://open.feishu.cn/) 机器人
    1. 前往[开发者平台](https://open.feishu.cn/app?lang=zh-CN)创建应用,并获取到 APPID 和 Secret
    2. 打开机器人能力
    3. 从cpolar或者serverless获得公网地址,例如`http://xxxx.r6.cpolar.top/webhook/event` ,在飞书机器人的 `事件订阅` 板块填写回调地址。
    4. 给订阅添加下列回调事件
        - im:message
        - im:message.group_at_msg
        - im:message.group_at_msg:readonly
        - im:message.p2p_msg
        - im:message.p2p_msg:readonly
        - im:message:send_as_bot
    5. 发布版本，等待企业管理员审核通过

更多介绍，参考[飞书上的小计算器: Go机器人来啦](https://www.bilibili.com/video/BV12M41187rV/)



### 相关阅读

- [go-cache](https://github.com/patrickmn/go-cache)

- [在Go语言项目中使用Zap日志库](https://www.liwenzhou.com/posts/Go/zap/)

- [飞书 User_ID、Open_ID 与 Union_ID 区别](https://www.feishu.cn/hc/zh-CN/articles/794300086214)

- [飞书重复接受到消息](https://open.feishu.cn/document/uAjLw4CM/ukTMukTMukTM/reference/im-v1/message/events/receive)
