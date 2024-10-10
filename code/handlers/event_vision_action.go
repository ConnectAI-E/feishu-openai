package handlers

import (
	"context"
	"fmt"
	"os"
	"start-feishubot/initialization"
	"start-feishubot/services"
	"start-feishubot/services/openai"
	"start-feishubot/utils"

	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
)

type VisionAction struct { /*图片推理*/
}

func (va *VisionAction) Execute(a *ActionInfo) bool {
	if !AzureModeCheck(a) {
		return true
	}

	if isVisionCommand(a) {
		initializeVisionMode(a)
		sendVisionInstructionCard(*a.ctx, a.info.sessionId, a.info.msgId)
		return false
	}

	mode := a.handler.sessionCache.GetMode(*a.info.sessionId)

	if a.info.msgType == "image" {
		if mode != services.ModeVision {
			sendVisionModeCheckCard(*a.ctx, a.info.sessionId, a.info.msgId)
			return false
		}

		return va.handleVisionImage(a)
	}

	if a.info.msgType == "post" && mode == services.ModeVision {
		return va.handleVisionPost(a)
	}

	return true
}

func isVisionCommand(a *ActionInfo) bool {
	_, foundPic := utils.EitherTrimEqual(a.info.qParsed, "/vision", "图片推理")
	return foundPic
}

func initializeVisionMode(a *ActionInfo) {
	a.handler.sessionCache.Clear(*a.info.sessionId)
	a.handler.sessionCache.SetMode(*a.info.sessionId, services.ModeVision)
	a.handler.sessionCache.SetVisionDetail(*a.info.sessionId, services.VisionDetailHigh)
}

func (va *VisionAction) handleVisionImage(a *ActionInfo) bool {
	detail := a.handler.sessionCache.GetVisionDetail(*a.info.sessionId)
	base64, err := downloadAndEncodeImage(a.info.imageKey, a.info.msgId)
	if err != nil {
		replyWithErrorMsg(*a.ctx, err, a.info.msgId)
		return false
	}

	return va.processImageAndReply(a, base64, detail)
}

func (va *VisionAction) handleVisionPost(a *ActionInfo) bool {
	detail := a.handler.sessionCache.GetVisionDetail(*a.info.sessionId)
	var base64s []string

	for _, imageKey := range a.info.imageKeys {
		if imageKey == "" {
			continue
		}
		base64, err := downloadAndEncodeImage(imageKey, a.info.msgId)
		if err != nil {
			replyWithErrorMsg(*a.ctx, err, a.info.msgId)
			return false
		}
		base64s = append(base64s, base64)
	}

	if len(base64s) == 0 {
		replyMsg(*a.ctx, "🤖️：请发送一张图片", a.info.msgId)
		return false
	}

	return va.processMultipleImagesAndReply(a, base64s, detail)
}

func downloadAndEncodeImage(imageKey string, msgId *string) (string, error) {
	f := fmt.Sprintf("%s.png", imageKey)
	defer os.Remove(f)

	req := larkim.NewGetMessageResourceReqBuilder().MessageId(*msgId).FileKey(imageKey).Type("image").Build()
	resp, err := initialization.GetLarkClient().Im.MessageResource.Get(context.Background(), req)
	if err != nil {
		return "", err
	}

	resp.WriteFile(f)
	return openai.GetBase64FromImage(f)
}

func replyWithErrorMsg(ctx context.Context, err error, msgId *string) {
	replyMsg(ctx, fmt.Sprintf("🤖️：图片下载失败，请稍后再试～\n 错误信息: %v", err), msgId)
}

func (va *VisionAction) processImageAndReply(a *ActionInfo, base64 string, detail string) bool {
	msg := createVisionMessages("解释这个图片", base64, detail)
	completions, err := a.handler.gpt.GetVisionInfo(msg)
	if err != nil {
		replyWithErrorMsg(*a.ctx, err, a.info.msgId)
		return false
	}
	sendVisionTopicCard(*a.ctx, a.info.sessionId, a.info.msgId, completions.Content)
	return false
}

func (va *VisionAction) processMultipleImagesAndReply(a *ActionInfo, base64s []string, detail string) bool {
	msg := createMultipleVisionMessages(a.info.qParsed, base64s, detail)
	completions, err := a.handler.gpt.GetVisionInfo(msg)
	if err != nil {
		replyWithErrorMsg(*a.ctx, err, a.info.msgId)
		return false
	}
	sendVisionTopicCard(*a.ctx, a.info.sessionId, a.info.msgId, completions.Content)
	return false
}

func createVisionMessages(query, base64Image, detail string) []openai.VisionMessages {
	return []openai.VisionMessages{
		{
			Role: "user",
			Content: []openai.ContentType{
				{Type: "text", Text: query},
				{Type: "image_url", ImageURL: &openai.ImageURL{
					URL:    "data:image/jpeg;base64," + base64Image,
					Detail: detail,
				}},
			},
		},
	}
}

func createMultipleVisionMessages(query string, base64Images []string, detail string) []openai.VisionMessages {
	content := []openai.ContentType{{Type: "text", Text: query}}
	for _, base64Image := range base64Images {
		content = append(content, openai.ContentType{
			Type: "image_url",
			ImageURL: &openai.ImageURL{
				URL:    "data:image/jpeg;base64," + base64Image,
				Detail: detail,
			},
		})
	}
	return []openai.VisionMessages{{Role: "user", Content: content}}
}
