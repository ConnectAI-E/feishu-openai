
<p align='center'>
    <img src='./docs/img3.png' alt='' width='800'/>
</p>

<br>

<p align='center'>
    在飞书上与ChatGPT随时对话
<br>
<br>
    Feishu ChatGpt
</p>

## 👻 机器人功能

- 🗣 畅所欲言，直接语音交流 🚧
- 💬 私人和群聊支持多话题同时对话，让讨论更加连贯和高效
- 🔄 持续对话联系上下文，回复对话框即可继续同一话题讨论
- ⏰ 超时自动结束对话，同时支持清除讨论历史
- 📝 支持富文本卡片回复，让信息更加丰富多彩
- 👍 交互式反馈提醒，及时获取机器人处理结果
- 🎭 支持角色扮演的场景模式，为讨论增添乐趣和创意
- 🏞 内置丰富的场景预设，让用户更方便地管理场景  🚧
- 🔙 轻松恢复历史话题，继续讨论  🚧
- 🖼 支持根据文本生成图片，增强信息呈现效果  🚧
- 🔒 内置管理员模式，使用更加安全可靠 🚧


<p align='center'>
    <img src='./docs/help.png' alt='' width='600'/>
</p>

## 🌟 项目特点
- 🍏 基于 OpenAI-[gpt-3.5-turbo](https://platform.openai.com/account/api-keys) 接口
- 🍎 通过 lark，将 ChatGPT 接入[飞书](https://open.feishu.cn/app)
- 🥒 支持[Serverless云函数](https://github.com/serverless-devs/serverless-devs)、[本地环境](https://dashboard.cpolar.com/login)、[Docker](https://www.docker.com/)、[二进制安装包](https://github.com/Leizhenpeng/feishu-chatgpt/releases/)多种渠道部署
- 🍋 基于[goCache](https://github.com/patrickmn/go-cache)内存键值对缓存


## 项目部署


######  有关飞书的配置文件说明，**[➡︎ 点击查看](#详细配置步骤)**


``` bash
git clone git@github.com:Leizhenpeng/feishu-chatgpt.git
cd feishu-chatgpt/code

```
<details>
    <summary>本地部署</summary>
<br>

如果你的服务器没有公网 IP，可以使用反向代理的方式

飞书的服务器在国内对ngrok的访问速度很慢，所以推荐使用一些国内的反向代理服务商
- [cpolar](https://dashboard.cpolar.com/)
- [natapp](https://natapp.cn/)


```bash
# 配置config.yaml
mv config.example.yaml config.yaml

//测试部署
go run main.go
cpolar http 9000

//正式部署
nohup cpolar http 9000 -log=stdout &

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
    <summary>serverless云函数(阿里云等)部署</summary>
<br>

安装[severless](https://docs.serverless-devs.com/serverless-devs/quick_start)工具
```bash
# 配置config.yaml
mv config.example.yaml config.yaml
# 安装severless cli
npm install @serverless-devs/s -g
```
一键部署
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

``` bash
docker build -t feishu-chatgpt:latest .
docker run -d --name feishu-chatgpt -p 9000:9000 \
--env APP_ID=xxx \
--env APP_SECRET=xxx \
--env APP_ENCRYPT_KEY=xxx \
--env APP_VERIFICATION_TOKEN=xxx \
--env BOT_NAME=chatGpt \
--env OPENAI_KEY=sk-xxx \
feishu-chatgpt:latest
```
------------

#### 如果需要指定OpenAi请求域名  --env API_URL=新的请求域名 \
``` bash
docker run -d --name feishu-chatgpt -p 9000:9000 \
--env APP_ID=xxx \
--env APP_SECRET=xxx \
--env APP_ENCRYPT_KEY=xxx \
--env APP_VERIFICATION_TOKEN=xxx \
--env BOT_NAME=chatGpt \
--env OPENAI_KEY=sk-xxx \
--env API_URL=https://api.openai.com \
feishu-chatgpt:latest
```
------------
小白简易化docker部署

- docker地址: https://hub.docker.com/r/leizhenpeng/feishu-chatgpt
``` bash
docker run -d --restart=always --name feishu-chatgpt2 -p 9000:9000 -v /etc/localtime:/etc/localtim:ro  \
--env APP_ID=xxx \
--env APP_SECRET=xxx \
--env APP_ENCRYPT_KEY=xxx \
--env APP_VERIFICATION_TOKEN=xxx \
--env BOT_NAME=chatGpt \
--env OPENAI_KEY=sk-xxx \
dockerproxy.com/leizhenpeng/feishu-chatgpt:latest
```

事件回调地址: http://IP:9000/webhook/event
卡片回调地址: http://IP:9000/webhook/card

把它填入飞书后台
<br>

</details>

<details>
    <summary>二进制安装包部署</summary>
<br>

1. 进入[release页面](https://github.com/Leizhenpeng/feishu-chatgpt/releases/) 下载对应的安装包
2. 解压安装包,修改config.example.yml中配置信息,另存为config.yml
3. 运行程序入口文件 `feishu-chatgpt`

</details>

## 详细配置步骤

-  获取 [OpenAI](https://platform.openai.com/account/api-keys) 的 KEY
-  创建 [飞书](https://open.feishu.cn/) 机器人
    1. 前往[开发者平台](https://open.feishu.cn/app?lang=zh-CN)创建应用,并获取到 APPID 和 Secret
    2. 前往`应用功能-机器人`, 创建机器人
    3. 从cpolar或者serverless获得公网地址,在飞书机器人后台的 `事件订阅` 板块填写。例如，
        - `http://xxxx.r6.cpolar.top`为cpolar暴露的公网地址
        - `/webhook/event`为统一的应用路由
        - 最终的回调地址为 `http://xxxx.r6.cpolar.top/webhook/event`
    4. 在飞书机器人后台的 `机器人` 板块，填写消息卡片请求网址。例如，
        - `http://xxxx.r6.cpolar.top`为cpolar暴露的公网地址
        - `/webhook/card`为统一的应用路由
        - 最终的消息卡片请求网址为 `http://xxxx.r6.cpolar.top/webhook/card`
    5. 给订阅添加下列回调事件
        - im:message
        - im:message.group_at_msg(获取群组中所有消息)
        - im:message.group_at_msg:readonly(接收群聊中@机器人消息事件)
        - im:message.p2p_msg(获取用户发给机器人的单聊消息)
        - im:message.p2p_msg:readonly(读取用户发给机器人的单聊消息)
        - im:message:send_as_bot(获取用户在群组中@机器人的消息)
        - im:chat:readonly(获取群组信息)
        - im:chat(获取与更新群组信息)
5. 发布版本，等待企业管理员审核通过

更多介绍，参考[飞书上的小计算器: Go机器人来啦](https://www.bilibili.com/video/BV12M41187rV/)


## 更多交流

企业如需定制部署，可联系WeChat: `laolei_forkway`，支持发票~

遇到其他问题，可以加入飞书群沟通~

<img src='./docs/talk.png' alt='' width='300'/>

## 赞助项目

如果你觉得这个项目对你有帮助，可以请作者买本书~

😚 谢谢你啦 😚

<details>
    <summary>赞赏二维码</summary>
    <img width="400" src="https://user-images.githubusercontent.com/50035229/222702169-cec777f9-cb0c-4f6a-90e5-418959e668e5.png">
</details>
