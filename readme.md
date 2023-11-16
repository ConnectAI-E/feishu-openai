<a href="https://www.connectai-e.com" target="_blank" rel="noreferrer">
  <img align="center" width="100%" height="auto" src="https://github-production-user-asset-6210df.s3.amazonaws.com/50035229/250274252-7f07a95e-b5aa-4dd8-90e7-5fb3bfb863c7.svg" >
</a>
<br>
<br>
<details align='center' style="margin-top: 10px">
    <summary> 📷 点击展开完整功能</summary>
    <br>
    <p align='center'>
        <img src='https://github.com/ConnectAI-E/Feishu-OpenAI/assets/110169811/b3543ae5-23cc-4fd7-829b-74656a300901' width='800'/>
        <img src='https://github.com/ConnectAI-E/Feishu-OpenAI/assets/110169811/948b051a-c25a-4825-8a18-5b72f6660b87' width='800'/>
        <img src='https://github.com/ConnectAI-E/Feishu-OpenAI/assets/110169811/a2139cc3-021e-4820-8e95-53a70541f136' width='800'/>
        <img src='https://github.com/ConnectAI-E/Feishu-OpenAI/assets/110169811/8e700709-ec0d-4c5e-b4af-710e42a0ba69' width='800'/>
        <img src='https://github.com/ConnectAI-E/Feishu-OpenAI/assets/110169811/001cb314-5fe2-4c93-86a7-66ea747b7855' width='800'/>
        <img src='https://github.com/ConnectAI-E/Feishu-OpenAI/assets/110169811/4dc12ca6-5305-4bc5-bf1a-40354c17ef29' width='800'/>
        <img src='https://github.com/ConnectAI-E/Feishu-OpenAI/assets/110169811/0c3178a0-e648-4013-a042-c1f83dbbe847' width='800'/>   
    </p>
</details>

<br>

<p align='center'>
   飞书 ×（GPT-4 + DALL·E + Whisper）
<br>
<br>
    🚀 Feishu OpenAI 🚀
</p>

  
<p align='center'>
   www.connectai-e.com
</p>


## 商业支持

如果开源版无法满足您公司的需求，推荐您体验商业共创版

- 内置开箱即用的 Azure Openai: 无需部署到海外，即可获得数十倍的性能提升
- 掌控全局的 Admin Panel: AI资源管理、对话日志查询、风险词规避和对话权限管理
- 专人技术支持: 配备专业部署交付人员与后期一对一维护服务
- 同时提供在线 Saas 版/企业级私有部署

查看更多内容: https://www.connectai-e.com

企业客户咨询: 15623677720

<img width="300" src="https://github.com/DDMeaqua/Feishu-OpenAI/assets/110169811/0a07f7e2-183f-4ee4-a40a-e33ee5c440e4">
<img width="300" src="https://github.com/ConnectAI-E/Feishu-OpenAI/assets/110169811/ef313ce0-158e-4314-a932-e6b1180cc750">



## 👻 机器人功能

🗣 语音交流：私人直接与机器人畅所欲言

💬 多话题对话：支持私人和群聊多话题讨论，高效连贯

🖼 文本成图：支持文本成图和以图搜图

🛖 场景预设：内置丰富场景列表，一键切换AI角色

🎭 角色扮演：支持场景模式，增添讨论乐趣和创意

🤖 AI模式：内置4种AI模式，感受AI的智慧与创意

🔄 上下文保留：回复对话框即可继续同一话题讨论

⏰ 自动结束：超时自动结束对话，支持清除讨论历史

📝 富文本卡片：支持富文本卡片回复，信息更丰富多彩

👍 交互式反馈：即时获取机器人处理结果

🎰 余额查询：即时获取token消耗情况

🔙 历史回档：轻松回档历史对话，继续话题讨论 🚧

🔒 管理员模式：内置管理员模式，使用更安全可靠 🚧

🌐 多token负载均衡：优化生产级别的高频调用场景

↩️ 支持反向代理：为不同地区的用户提供更快、更稳定的访问体验

📚 与飞书文档互动：成为企业员工的超级助手 🚧

🎥 话题内容秒转PPT：让你的汇报从此变得更加简单 🚧

📊 表格分析：轻松导入飞书表格，提升数据分析效率 🚧

🍊 私有数据训练：利用公司产品信息对GPT二次训练，更好地满足客户个性化需求 🚧



## 🌟 项目特点

- 🍏 对话基于 OpenAI (https://platform.openai.com/account/api-keys) 接口
- 🍎 通过 lark，将 ChatGPT 接入[飞书](https://open.feishu.cn/app)和[飞书国际版](https://www.larksuite.com/)
- 🥒
  支持[Serverless 云函数](https://github.com/serverless-devs/serverless-devs)、[本地环境](https://dashboard.cpolar.com/login)、[Docker](https://www.docker.com/)、[二进制安装包](https://github.com/Leizhenpeng/feishu-chatgpt/releases/)
  等多种渠道部署
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
    <summary>使用 Railway 平台一键部署</summary>


Railway 是一家国外的 Serverless 平台，支持多种语言，可以一键将 GitHub 上的代码仓库部署到 Railway 平台，然后在 Railway
平台上配置环境变量即可。部署本项目的流程如下：

#### 1. 生成 Railway 项目

点击下方按钮即可创建一个对应的 Railway 项目，其会自动 Fork 本项目到你的 GitHub 账号下。

[![Deploy on Railway](https://railway.app/button.svg)](https://railway.app/template/10D-TF?referralCode=oMcVS2)

#### 2. 配置环境变量

在打开的页面中，配置环境变量，每个变量的说明如下图所示：


<img src='https://user-images.githubusercontent.com/50035229/225005602-88d8678f-9d17-4dc5-8d1e-4abf64fb84fd.png' alt='Railway 环境变量' width='500px'/>

#### 3. 部署项目

填写完环境变量后，点击 Deploy 就完成了项目的部署。部署完成后还需获取对应的域名用于飞书机器人访问，如下图所示：

<img src='https://user-images.githubusercontent.com/50035229/225006236-57cb3c8a-1b7d-4bfe-9c9b-099cb9179027.png' alt='Railway 域名' width='500px'/>

如果不确定自己部署是否成功，可以通过访问上述获取到的域名 (https://xxxxxxxx.railway.app/ping) 来查看是否返回了`pong`
，如果返回了`pong`，说明部署成功。

</details>


<details>
    <summary>Repl.it部署</summary>
<br>

The fastest way to deploy the feishu-openai to `repl.it` is to click the `run on repl.it` button below.

<a href="https://replit.com/@Leizhenpeng/Feishu-OpenAI#.replit">
  <img alt="Run on Repl.it" src="https://repl.it/badge/github/alist-org/alist-replit" style="height: 40px; width: 190px;" />
</a>

Remember switch to `secrets` tab then edit `System environment variables`.You can also edit raw json:

```json
{
  "APP_ID": "",
  "APP_SECRET": "",
  "APP_ENCRYPT_KEY": "",
  "APP_VERIFICATION_TOKEN": "",
  "BOT_NAME": "ChatGPT",
  "OPENAI_KEY": "sk-",
  "OPENAI_MODEL": "gpt-3.5-turbo"
}
```

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
- `HTTP_PROXY` 为宿主机的proxy地址，例如 `http://host.docker.internal:7890`,没有代理的话，可以不用设置
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

--- 

部署azure版本

```bash
docker build -t feishu-chatgpt:latest .
docker run -d --name feishu-chatgpt -p 9000:9000 \
--env APP_ID=xxx \
--env APP_SECRET=xxx \
--env APP_ENCRYPT_KEY=xxx \
--env APP_VERIFICATION_TOKEN=xxx \
--env BOT_NAME=chatGpt \
--env AZURE_ON=true \
--env AZURE_API_VERSION=xxx \
--env AZURE_RESOURCE_NAME=xxx \
--env AZURE_DEPLOYMENT_NAME=xxx \
--env AZURE_OPENAI_TOKEN=xxx \
feishu-chatgpt:latest
```

注意:

- `BOT_NAME` 为飞书机器人名称，例如 `chatGpt`
- `AZURE_ON` 为是否使用azure ,请填写 `true`
- `AZURE_API_VERSION` 为azure api版本 例如 `2023-03-15-preview`
- `AZURE_RESOURCE_NAME` 为azure 资源名称 类似 `https://{AZURE_RESOURCE_NAME}.openai.azure.com`
- `AZURE_DEPLOYMENT_NAME` 为azure 部署名称 类似 `https://{AZURE_RESOURCE_NAME}.openai.azure.com/deployments/{AZURE_DEPLOYMENT_NAME}/chat/completions`
- `AZURE_OPENAI_TOKEN` 为azure openai token

</details>

<details>
    <summary>docker-compose 部署</summary>
<br>

编辑 docker-compose.yaml，通过 environment 配置相应环境变量（或者通过 volumes 挂载相应配置文件），然后运行下面的命令即可

```bash
# 构建镜像
docker compose build

# 启动服务
docker compose up -d

# 停止服务
docker compose down
```

事件回调地址: http://IP:9000/webhook/event
卡片回调地址: http://IP:9000/webhook/card

</details>

<details>
    <summary>二进制安装包部署</summary>
<br>

1. 进入[release 页面](https://github.com/Leizhenpeng/feishu-chatgpt/releases/) 下载对应的安装包
2. 解压安装包,修改 config.example.yml 中配置信息,另存为 config.yaml
3. 目录下添加文件 `role_list.yaml`，自定义角色，可以从这里获取：[链接](https://github.com/Leizhenpeng/feishu-chatgpt/blob/master/code/role_list.yaml)
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
      <img src='https://user-images.githubusercontent.com/50035229/230003546-36450f2f-b6e9-4292-8b40-3a4aa8a05a64.png' alt='' width='800'/>
      <img src='https://user-images.githubusercontent.com/50035229/223945122-f7ab3d9a-6742-43d2-970e-ddb0f284c7fa.png' alt='' width='800'/>
      <img src='https://user-images.githubusercontent.com/50035229/223944507-8d1a08d7-8b5b-4f32-a90d-fd338164ec82.png' alt='' width='800'/>
      <img src='https://user-images.githubusercontent.com/50035229/223944515-fb505e84-c840-484a-8df5-612f60bf27ea.png' alt='' width='800'/>
      <img src='https://user-images.githubusercontent.com/50035229/223944590-ad61320f-c14a-4542-80ad-dee2e6469b67.png' alt='' width='800'/>
    </p>
</details>


- 获取 [OpenAI](https://platform.openai.com/account/api-keys) 的 KEY( 🙉 下面有免费的 KEY 供大家测试部署 )
- 创建 [飞书](https://open.feishu.cn/) 机器人
    1. 前往[开发者平台](https://open.feishu.cn/app?lang=zh-CN)创建应用,并获取到 APPID 和 Secret
    2. 前往`应用功能-机器人`, 创建机器人
    3. 从 cpolar、serverless 或 Railway 获得公网地址，在飞书机器人后台的 `事件订阅` 板块填写。例如，
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


## 一起交流

遇到问题，可以加入飞书群沟通~

<img src='https://user-images.githubusercontent.com/13283837/232570671-1058555f-c9e5-4f64-889b-1d8efd0101ba.png' alt='' width='200'/>


## 企联AI

| <div style="width:200px">AI</div> |             <img width=120> SDK <img width=120>              |                         Application                          |
| :-------------------------------: | :----------------------------------------------------------: | :----------------------------------------------------------: |
|              🎒OpenAI              |    [Go-OpenAI](https://github.com/ConnectAI-E/Go-OpenAI)     | [🏅Feishu-OpenAI](https://github.com/ConnectAI-E/Feishu-OpenAI), [🎖Lark-OpenAI](https://github.com/ConnectAI-E/Lark-OpenAI), [Feishu-EX-ChatGPT](https://github.com/ConnectAI-E/Feishu-EX-ChatGPT), [🎖Feishu-OpenAI-Stream-Chatbot](https://github.com/ConnectAI-E/Feishu-OpenAI-Stream-Chatbot), [Feishu-TLDR](https://github.com/ConnectAI-E/Feishu-TLDR),[Feishu-OpenAI-Amazing](https://github.com/ConnectAI-E/Feishu-OpenAI-Amazing), [Feishu-Oral-Friend](https://github.com/ConnectAI-E/Feishu-Oral-Friend), [Feishu-OpenAI-Base-Helper](https://github.com/ConnectAI-E/Feishu-OpenAI-Base-Helper), [Feishu-Vector-Knowledge-Management](https://github.com/ConnectAI-E/Feishu-Vector-Knowledge-Management), [Feishu-OpenAI-PDF-Helper](https://github.com/ConnectAI-E/Feishu-OpenAI-PDF-Helper), [🏅Dingtalk-OpenAI](https://github.com/ConnectAI-E/Dingtalk-OpenAI), [Wework-OpenAI](https://github.com/ConnectAI-E/Wework-OpenAI), [WeWork-OpenAI-Node](https://github.com/ConnectAI-E/WeWork-OpenAI-Node), [llmplugin](https://github.com/ConnectAI-E/llmplugin) |
|             🤖 AutoGPT             |                            ------                            | [🏅AutoGPT-Next-Web](https://github.com/ConnectAI-E/AutoGPT-Next-Web) |
|         🎭 Stablediffusion         |                            ------                            | [🎖Feishu-Stablediffusion](https://github.com/ConnectAI-E/Feishu-Stablediffusion) |
|           🍎 Midjourney            | [Go-Midjourney](https://github.com/ConnectAI-E/Go-Midjourney) | [🏅Feishu-Midjourney](https://github.com/ConnectAI-E/Feishu-Midjourney), [🔥MidJourney-Web](https://github.com/ConnectAI-E/MidJourney-Web), [Dingtalk-Midjourney](https://github.com/ConnectAI-E/Dingtalk-Midjourney) |
|            🍍 文心一言             |    [Go-Wenxin](https://github.com/ConnectAI-E/Go-Wenxin)     | [Feishu-Wenxin](https://github.com/ConnectAI-E/Feishu-Wenxin), [Dingtalk-Wenxin](https://github.com/ConnectAI-E/Dingtalk-Wenxin), [Wework-Wenxin](https://github.com/ConnectAI-E/Wework-Wenxin) |
|             💸 Minimax             |   [Go-Minimax](https://github.com/ConnectAI-E/Go-Minimax)    | [Feishu-Minimax](https://github.com/ConnectAI-E/Feishu-Minimax), [Dingtalk-Minimax](https://github.com/ConnectAI-E/Dingtalk-Minimax), [Wework-Minimax](https://github.com/ConnectAI-E/Wework-Minimax) |
|             ⛳️ CLAUDE              |    [Go-Claude](https://github.com/ConnectAI-E/Go-Claude)     | [Feishu-Claude](https://github.com/ConnectAI-E/Feishu-Claude), [DingTalk-Claude](https://github.com/ConnectAI-E/DingTalk-Claude), [Wework-Claude](https://github.com/ConnectAI-E/Wework-Claude) |
|              🥁 PaLM               |      [Go-PaLM](https://github.com/ConnectAI-E/go-PaLM)       | [Feishu-PaLM](https://github.com/ConnectAI-E/Feishu-PaLM),[DingTalk-PaLM](https://github.com/ConnectAI-E/DingTalk-PaLM),[Wework-PaLM](https://github.com/ConnectAI-E/Wework-PaLM) |
|             🎡 Prompt              |                            ------                            | [📖 Prompt-Engineering-Tutior](https://github.com/ConnectAI-E/Prompt-Engineering-Tutior) |
|             🍋 ChatGLM             |                            ------                            | [Feishu-ChatGLM](https://github.com/ConnectAI-E/Feishu-ChatGLM) |
|            ⛓ LangChain            |                            ------                            | [📖 LangChain-Tutior](https://github.com/ConnectAI-E/LangChain-Tutior) |
|            🪄 One-click            |                            ------                            | [🎖Awesome-One-Click-Deployment](https://github.com/ConnectAI-E/Awesome-One-Click-Deployment) |






开源社区：https://github.com/ConnectAI-E

产品日志：https://www.connectai-e.com/logger



