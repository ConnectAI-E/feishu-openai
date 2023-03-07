package handlers

import (
	"context"
	"fmt"
	"start-feishubot/services"
	"start-feishubot/utils"
)

type ActionInfo struct {
	p         *PersonalMessageHandler
	msgId     *string
	chatId    *string
	qParsed   string
	ctx       *context.Context
	sessionId *string
}

type Action interface {
	Execute(data *ActionInfo) bool
}

type ProcessedAction struct { //消息唯一性
}

func (*ProcessedAction) Execute(data *ActionInfo) bool {
	if data.p.msgCache.IfProcessed(*data.msgId) {
		return false
	}
	data.p.msgCache.TagProcessed(*data.msgId)
	return true
}

type EmptyAction struct { /*空消息*/
}

func (*EmptyAction) Execute(data *ActionInfo) bool {
	if len(data.qParsed) == 0 {
		sendMsg(*data.ctx, "🤖️：你想知道什么呢~", data.chatId)
		fmt.Println("msgId", *data.msgId, "message.text is empty")
		return false
	}
	return true
}

type ClearAction struct { /*清除消息*/
}

func (*ClearAction) Execute(data *ActionInfo) bool {
	if _, foundClear := utils.EitherTrimEqual(data.qParsed, "/clear", "清除"); foundClear {
		sendClearCacheCheckCard(*data.ctx, data.sessionId, data.msgId)
		return false
	}
	return true
}

type RolePlayAction struct { /*角色扮演*/
}

func (*RolePlayAction) Execute(data *ActionInfo) bool {
	if system, foundSystem := utils.EitherCutPrefix(data.qParsed, "/system ", "角色扮演 "); foundSystem {
		data.p.sessionCache.Clear(*data.sessionId)
		systemMsg := append([]services.Messages{}, services.Messages{
			Role: "system", Content: system,
		})
		data.p.sessionCache.SetMsg(*data.sessionId, systemMsg)
		sendSystemInstructionCard(*data.ctx, data.sessionId, data.msgId, system)
		return false
	}
	return true
}

type HelpAction struct { /*帮助*/
}

func (*HelpAction) Execute(data *ActionInfo) bool {
	if _, foundHelp := utils.EitherTrimEqual(data.qParsed, "/help", "帮助"); foundHelp {
		sendHelpCard(*data.ctx, data.sessionId, data.msgId)
		return false
	}
	return true
}

type PicAction struct { /*图片*/
}

func (*PicAction) Execute(data *ActionInfo) bool {
	// 开启图片创作模式
	if _, foundPic := utils.EitherTrimEqual(data.qParsed,
		"/picture", "图片创作"); foundPic {
		data.p.sessionCache.Clear(*data.sessionId)
		data.p.sessionCache.SetMode(*data.sessionId,
			services.ModePicCreate)
		sendPicCreateInstructionCard(*data.ctx, data.sessionId,
			data.msgId)
		return false
	}

	// 生成图片
	mode := data.p.sessionCache.GetMode(*data.sessionId)
	if mode == services.ModePicCreate {
		bs64, err := data.p.gpt.GenerateOneImage(data.qParsed,
			"256x256")
		if err != nil {
			replyMsg(*data.ctx, fmt.Sprintf(
				"🤖️：图片生成失败，请稍后再试～\n错误信息: %v", err), data.msgId)
			return false
		}
		replayImageByBase64(*data.ctx, bs64, data.msgId)
		return false
	}

	return true
}

type MessageAction struct { /*消息*/
}

func (*MessageAction) Execute(data *ActionInfo) bool {
	msg := data.p.sessionCache.GetMsg(*data.sessionId)
	msg = append(msg, services.Messages{
		Role: "user", Content: data.qParsed,
	})
	completions, err := data.p.gpt.Completions(msg)
	if err != nil {
		replyMsg(*data.ctx, fmt.Sprintf("🤖️：消息机器人摆烂了，请稍后再试～\n错误信息: %v", err), data.msgId)
		return false
	}
	msg = append(msg, completions)
	data.p.sessionCache.SetMsg(*data.sessionId, msg)
	//if new topic
	if len(msg) == 2 {
		fmt.Println("new topic", msg[1].Content)
		sendNewTopicCard(*data.ctx, data.sessionId, data.msgId, completions.Content)
		return false
	}
	err = replyMsg(*data.ctx, completions.Content, data.msgId)
	if err != nil {
		replyMsg(*data.ctx, fmt.Sprintf("🤖️：消息机器人摆烂了，请稍后再试～\n错误信息: %v", err), data.msgId)
		return false
	}
	return true
}
