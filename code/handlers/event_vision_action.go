package handlers

import (
	"context"
	"fmt"
	"os"
	"start-feishubot/initialization"
	"start-feishubot/logger"
	"start-feishubot/services"
	"start-feishubot/services/openai"
	"start-feishubot/utils"

	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
)

type VisionAction struct { /*图片推理*/
}

func (*VisionAction) Execute(a *ActionInfo) bool {
	check := AzureModeCheck(a)
	if !check {
		return true
	}
	// 开启图片创作模式
	if _, foundPic := utils.EitherTrimEqual(a.info.qParsed,
		"/vision", "图片推理"); foundPic {
		a.handler.sessionCache.Clear(*a.info.sessionId)
		a.handler.sessionCache.SetMode(*a.info.sessionId,
			services.ModeVision)
		a.handler.sessionCache.SetVisionDetail(*a.info.sessionId,
			services.VisionDetailHigh)
		sendVisionInstructionCard(*a.ctx, a.info.sessionId,
			a.info.msgId)
		return false
	}

	mode := a.handler.sessionCache.GetMode(*a.info.sessionId)
	fmt.Println("a.info.msgType: ", a.info.msgType)

	logger.Debug("MODE:", mode)

	// 收到一张图片,且不在图片推理模式下, 提醒是否切换到图片推理模式
	if a.info.msgType == "image" && mode != services.ModeVision {
		sendVisionModeCheckCard(*a.ctx, a.info.sessionId, a.info.msgId)
		return false
	}

	// todo
	//return false

	if a.info.msgType == "image" && mode == services.ModeVision {
		//保存图片
		imageKey := a.info.imageKey
		//fmt.Printf("fileKey: %s \n", imageKey)
		msgId := a.info.msgId
		//fmt.Println("msgId: ", *msgId)
		req := larkim.NewGetMessageResourceReqBuilder().MessageId(
			*msgId).FileKey(imageKey).Type("image").Build()
		resp, err := initialization.GetLarkClient().Im.MessageResource.Get(context.Background(), req)
		fmt.Println(resp, err)
		if err != nil {
			//fmt.Println(err)
			replyMsg(*a.ctx, fmt.Sprintf("🤖️：图片下载失败，请稍后再试～\n 错误信息: %v", err),
				a.info.msgId)
			return false
		}

		f := fmt.Sprintf("%s.png", imageKey)
		fmt.Println(f)
		resp.WriteFile(f)
		defer os.Remove(f)

		base64, err := openai.GetBase64FromImage(f)
		if err != nil {
			replyMsg(*a.ctx, fmt.Sprintf("🤖️：图片下载失败，请稍后再试～\n 错误信息: %v", err),
				a.info.msgId)
			return false
		}
		//
		var msg []openai.VisionMessages
		detail := a.handler.sessionCache.GetVisionDetail(*a.info.sessionId)
		// 如果没有提示词，默认模拟ChatGPT

		content2 := []openai.ContentType{
			{Type: "text", Text: "图片里面有什么", ImageURL: nil},
			{Type: "image_url", ImageURL: &openai.ImageURL{
				URL:    "data:image/jpeg;base64," + base64,
				Detail: detail,
			}},
		}

		msg = append(msg, openai.VisionMessages{
			Role: "user", Content: content2,
		})

		// get ai mode as temperature
		fmt.Println("msg: ", msg)
		completions, err := a.handler.gpt.GetVisionInfo(msg)
		if err != nil {
			replyMsg(*a.ctx, fmt.Sprintf(
				"🤖️：消息机器人摆烂了，请稍后再试～\n错误信息: %v", err), a.info.msgId)
			return false
		}
		sendOldTopicCard(*a.ctx, a.info.sessionId, a.info.msgId,
			completions.Content)
		return false
		//a.handler.sessionCache.SetMsg(*a.info.sessionId, msg)

	}

	if a.info.msgType == "post" && mode == services.ModeVision {
		fmt.Println(a.info.imageKeys)
		fmt.Println(a.info.qParsed)
		imagesKeys := a.info.imageKeys
		var base64s []string
		if len(imagesKeys) == 0 {
			replyMsg(*a.ctx, "🤖️：请发送一张图片", a.info.msgId)
			return false
		}
		//保存图片
		for i := 0; i < len(imagesKeys); i++ {
			if imagesKeys[i] == "" {
				continue
			}
			imageKey := imagesKeys[i]
			msgId := a.info.msgId
			//fmt.Println("msgId: ", *msgId)
			req := larkim.NewGetMessageResourceReqBuilder().MessageId(
				*msgId).FileKey(imageKey).Type("image").Build()
			resp, err := initialization.GetLarkClient().Im.MessageResource.Get(context.Background(), req)
			if err != nil {
				//fmt.Println(err)
				replyMsg(*a.ctx, fmt.Sprintf("🤖️：图片下载失败，请稍后再试～\n 错误信息: %v", err),
					a.info.msgId)
				return false
			}

			f := fmt.Sprintf("%s.png", imageKey)
			fmt.Println(f)
			resp.WriteFile(f)
			defer os.Remove(f)

			base64, err := openai.GetBase64FromImage(f)
			base64s = append(base64s, base64)
			if err != nil {
				replyMsg(*a.ctx, fmt.Sprintf("🤖️：图片下载失败，请稍后再试～\n 错误信息: %v", err),
					a.info.msgId)
				return false
			}
		}

		var msg []openai.VisionMessages
		detail := a.handler.sessionCache.GetVisionDetail(*a.info.sessionId)
		// 如果没有提示词，默认模拟ChatGPT

		content0 := []openai.ContentType{
			{Type: "text", Text: a.info.qParsed, ImageURL: nil},
		}
		// 循环数组
		for i := 0; i < len(base64s); i++ {
			content1 := []openai.ContentType{
				{Type: "image_url", ImageURL: &openai.ImageURL{
					URL:    "data:image/jpeg;base64," + base64s[i],
					Detail: detail,
				}},
			}
			content0 = append(content0, content1...)
		}

		msg = append(msg, openai.VisionMessages{
			Role: "user", Content: content0,
		})

		// get ai mode as temperature
		fmt.Println("msg: ", msg)
		completions, err := a.handler.gpt.GetVisionInfo(msg)
		if err != nil {
			replyMsg(*a.ctx, fmt.Sprintf(
				"🤖️：消息机器人摆烂了，请稍后再试～\n错误信息: %v", err), a.info.msgId)
			return false
		}
		sendOldTopicCard(*a.ctx, a.info.sessionId, a.info.msgId,
			completions.Content)
		return false
		//a.handler.sessionCache.SetMsg(*a.info.sessionId, msg)

		return false

	}

	return true
}
