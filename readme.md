<p align='center'>
    <img src='./docs/img3.png' alt='' width='800'/>
</p>

<details align='center'>
    <summary> 📷 点击展开完整功能截图</summary>
    <br>
    <p align='center'>
    <img src='https://user-images.githubusercontent.com/50035229/223590381-ed38db74-39f3-4e77-bd3a-aaa54d679286.png' alt='' width='800'/>
    <img src='https://user-images.githubusercontent.com/50035229/223590817-37a56eac-ab6e-4293-862a-de0988ac50b7.png' alt='' width='800'/>
    <img src='https://user-images.githubusercontent.com/50035229/223591272-88e4218c-5cb6-4925-8ede-18359bdf9097.png' alt='' width='800'/>
    </p>

</details>

<br>

<p align='center'>
   飞书 ×（GPT-3.5 + DALL·E + Whisper）
<br>
<br>
    🚀 Feishu OpenAI 🚀
</p>


## 👻 机器人功能

🗣 语音交流：直接与机器人畅所欲言 🚧

💬 多话题对话：支持私人和群聊多话题讨论，高效连贯

🖼 文本成图：支持文本成图和以图搜图

🎭 角色扮演：支持场景模式，增添讨论乐趣和创意

🔄 上下文保留：回复对话框即可继续同一话题讨论

⏰ 自动结束：超时自动结束对话，支持清除讨论历史

📝 富文本卡片：支持富文本卡片回复，信息更丰富多彩

👍 交互式反馈：即时获取机器人处理结果

🏞 场景预设：内置丰富场景预设，方便用户管理场景 🚧

🔙 历史回档：轻松回档历史对话，继续话题讨论 🚧

🔒 管理员模式：内置管理员模式，使用更安全可靠 🚧

🌐 多token负载均衡：优化生产级别的高频调用场景

↩️ 支持反向代理：方便国内用户访问

📚 与飞书文档互动：成为企业员工的超级助手 🚧


## 🌟 项目特点

- 🍏 基于 OpenAI-[gpt-3.5-turbo](https://platform.openai.com/account/api-keys) 接口
- 🍎 通过 lark，将 ChatGPT 接入[飞书](https://open.feishu.cn/app)
- 🥒 支持[Serverless 云函数](https://github.com/serverless-devs/serverless-devs)、[本地环境](https://dashboard.cpolar.com/login)、[Docker](https://www.docker.com/)、[二进制安装包](https://github.com/Leizhenpeng/feishu-chatgpt/releases/)多种渠道部署
- 🍋 基于[goCache](https://github.com/patrickmn/go-cache)内存键值对缓存

## 项目部署

###### 有关飞书的配置文件说明，**[➡︎ 点击查看](#详细配置步骤)**

<details>
    <summary>本地部署</summary>
<br>

```bash
git clone git@github.com:Leizhenpeng/feishu-chatgpt.git
cd feishu-chatgpt/code
```

如果你的服务器没有公网 IP，可以使用反向代理的方式

飞书的服务器在国内对 ngrok 的访问速度很慢，所以推荐使用一些国内的反向代理服务商

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

更多详细介绍，参考[飞书上的小计算器: Go 机器人来啦](https://www.bilibili.com/video/BV1nW4y1378T/)

<br>

</details>

<details>
    <summary>serverless云函数(阿里云等)部署</summary>
<br>

```bash
git clone git@github.com:Leizhenpeng/feishu-chatgpt.git
cd feishu-chatgpt/code
```

安装[severless](https://docs.serverless-devs.com/serverless-devs/quick_start)工具

```bash
# 配置config.yaml
mv config.example.yaml config.yaml
# 安装severless cli
npm install @serverless-devs/s -g
```

安装完成后，请根据您本地环境，根据下面教程部署`severless`

- 本地 `linux`/`mac os` 环境

1. 修改`s.yaml`中的部署地区和部署秘钥

```
edition: 1.0.0
name: feishuBot-chatGpt
access: "aliyun" #  修改自定义的秘钥别称

vars: # 全局变量
region: "cn-hongkong" # 修改云函数想要部署地区

```

2. 一键部署

```bash
cd ..
s deploy
```

- 本地`windows`

1. 首先打开本地`cmd`命令提示符工具，运行`go env`检查你电脑上 go 环境变量设置, 确认以下变量和值

```cmd
set GO111MODULE=on
set GOARCH=amd64
set GOOS=linux
set CGO_ENABLED=0
```

如果值不正确，比如您电脑上为`set GOOS=windows`, 请运行以下命令设置`GOOS`变量值

```cmd
go env -w GOOS=linux
```

2. 修改`s.yaml`中的部署地区和部署秘钥

```
edition: 1.0.0
name: feishuBot-chatGpt
access: "aliyun" #  修改自定义的秘钥别称

vars: # 全局变量
  region: "cn-hongkong" #  修改云函数想要部署地区

```

3. 修改`s.yaml`中的`pre-deploy`, 去除第二步`run`前面的环变量改置部分

```
  pre-deploy:
        - run: go mod tidy
          path: ./code
        - run: go build -o
            target/main main.go  # 删除GO111MODULE=on GOOS=linux GOARCH=amd64 CGO_ENABLED=0
          path: ./code

```

4. 一键部署

```bash
cd ..
s deploy
```

更多详细介绍，参考[仅需 1min，用 Serverless 部署基于 gin 的飞书机器人](https://www.bilibili.com/video/BV1nW4y1378T/)
<br>

</details>

<details>
    <summary>docker部署</summary>
<br>

```bash
docker build -t feishu-chatgpt:latest .
docker run -d --name feishu-chatgpt -p 9000:9000 \
--env APP_ID=xxx \
--env APP_SECRET=xxx \
--env APP_ENCRYPT_KEY=xxx \
--env APP_VERIFICATION_TOKEN=xxx \
--env BOT_NAME=chatGpt \
--env OPENAI_KEY="sk-xxx1,sk-xxx2,sk-xxx3" \
--env API_URL="https://api.openai.com" \
--env HTTP_PROXY="" \
feishu-chatgpt:latest
```
注意:
- `BOT_NAME` 为飞书机器人名称，例如 `chatGpt`
- `OPENAI_KEY` 为openai key，多个key用逗号分隔，例如 `sk-xxx1,sk-xxx2,sk-xxx3`
- `HTTP_PROXY` 为宿主机的proxy地址，例如 `http://host.docker.internal:7890`
- `API_URL` 为openai api 接口地址，例如 `https://api.openai.com`, 没有反向代理的话，可以不用设置
---

小白简易化 docker 部署

- docker 地址: https://hub.docker.com/r/leizhenpeng/feishu-chatgpt

```bash
docker run -d --restart=always --name feishu-chatgpt2 -p 9000:9000 -v /etc/localtime:/etc/localtim:ro  \
--env APP_ID=xxx \
--env APP_SECRET=xxx \
--env APP_ENCRYPT_KEY=xxx \
--env APP_VERIFICATION_TOKEN=xxx \
--env BOT_NAME=chatGpt \
--env OPENAI_KEY="sk-xxx1,sk-xxx2,sk-xxx3" \
--env API_URL=https://api.openai.com \
--env HTTP_PROXY="" \
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

1. 进入[release 页面](https://github.com/Leizhenpeng/feishu-chatgpt/releases/) 下载对应的安装包
2. 解压安装包,修改 config.example.yml 中配置信息,另存为 config.yml
3. 运行程序入口文件 `feishu-chatgpt`

事件回调地址: http://IP:9000/webhook/event
卡片回调地址: http://IP:9000/webhook/card

</details>

## 详细配置步骤

<details align='left'>
    <summary> 📸 点击展开飞书机器人配置的分步截图指导</summary>
    <br>
    <p align='center'>
    <img src='https://user-images.githubusercontent.com/50035229/223943381-39e0466f-2a5e-472a-9863-94eafb5f17b0.png' alt='' width='800'/>
    <img src='https://user-images.githubusercontent.com/50035229/223943448-228de5cb-0929-4d80-8087-8d8624dd6ddf.png' alt='' width='800'/>
    <img src='https://user-images.githubusercontent.com/50035229/223943485-ef331784-7940-4657-b128-70c98391e72f.png' alt='' width='800'/>
    <img src='https://user-images.githubusercontent.com/50035229/223943527-60e6653a-eb6e-4062-a076-b6c9da934352.png' alt='' width='800'/>
    <img src='https://user-images.githubusercontent.com/50035229/223943972-f49adf9f-af5f-463a-8c7a-c1f0cac0e8c3.png' alt='' width='800'/>
      <img src='https://user-images.githubusercontent.com/50035229/223944060-7ef630a4-4248-4509-852b-cad8bfffeefc.png' alt='' width='800'/>
      <img src='https://user-images.githubusercontent.com/50035229/223944230-aff586be-31cc-40de-9b1a-7d4e259d54dd.png' alt='' width='800'/>
      <img src='https://user-images.githubusercontent.com/50035229/223944350-917d115c-6c82-4d8b-9ec8-b5c82331a2dc.png' alt='' width='800'/>
      <img src='https://user-images.githubusercontent.com/50035229/223944381-97396156-f5e2-467f-aaf6-b1f6e1c446b2.png' alt='' width='800'/>
      <img src='https://user-images.githubusercontent.com/50035229/223945122-f7ab3d9a-6742-43d2-970e-ddb0f284c7fa.png' alt='' width='800'/>
      <img src='https://user-images.githubusercontent.com/50035229/223944507-8d1a08d7-8b5b-4f32-a90d-fd338164ec82.png' alt='' width='800'/>
      <img src='https://user-images.githubusercontent.com/50035229/223944515-fb505e84-c840-484a-8df5-612f60bf27ea.png' alt='' width='800'/>
      <img src='https://user-images.githubusercontent.com/50035229/223944590-ad61320f-c14a-4542-80ad-dee2e6469b67.png' alt='' width='800'/>
    </p>
</details>


- 获取 [OpenAI](https://platform.openai.com/account/api-keys) 的 KEY
- 创建 [飞书](https://open.feishu.cn/) 机器人
  1. 前往[开发者平台](https://open.feishu.cn/app?lang=zh-CN)创建应用,并获取到 APPID 和 Secret
  2. 前往`应用功能-机器人`, 创建机器人
  3. 从 cpolar 或者 serverless 获得公网地址,在飞书机器人后台的 `事件订阅` 板块填写。例如，
     - `http://xxxx.r6.cpolar.top`为 cpolar 暴露的公网地址
     - `/webhook/event`为统一的应用路由
     - 最终的回调地址为 `http://xxxx.r6.cpolar.top/webhook/event`
  4. 在飞书机器人后台的 `机器人` 板块，填写消息卡片请求网址。例如，
     - `http://xxxx.r6.cpolar.top`为 cpolar 暴露的公网地址
     - `/webhook/card`为统一的应用路由
     - 最终的消息卡片请求网址为 `http://xxxx.r6.cpolar.top/webhook/card`
  5. 在事件订阅板块，搜索三个词`机器人进群`、 `接收消息`、 `消息已读`, 把他们后面所有的权限全部勾选。
  进入权限管理界面，搜索`图片`, 勾选`获取与上传图片或文件资源`。
  最终会添加下列回调事件
     - im:resource(获取与上传图片或文件资源)
     - im:message
     - im:message.group_at_msg(获取群组中所有消息)
     - im:message.group_at_msg:readonly(接收群聊中@机器人消息事件)
     - im:message.p2p_msg(获取用户发给机器人的单聊消息)
     - im:message.p2p_msg:readonly(读取用户发给机器人的单聊消息)
     - im:message:send_as_bot(获取用户在群组中@机器人的消息)
     - im:chat:readonly(获取群组信息)
     - im:chat(获取与更新群组信息)
 

5. 发布版本，等待企业管理员审核通过

更多介绍，参考[飞书上的小计算器: Go 机器人来啦](https://www.bilibili.com/video/BV12M41187rV/)

## 更多交流

如需协助部署，或者其他定制服务，可联系下面的WeChat，支持发票~

遇到问题，可以加入飞书群沟通~

<img src='./docs/talk.png' alt='' width='200'/>

## 交朋友 或者 鼓励一下

如果你觉得这个项目对你有帮助，可以请作者买本书~

<img width="400" src="https://user-images.githubusercontent.com/50035229/224462896-28f7b2d5-f443-4cc2-9790-7b72e5c53f15.png">

😚 谢谢你啦 😚

