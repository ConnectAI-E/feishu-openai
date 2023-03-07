package handlers

import (
	"context"
	"fmt"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	"start-feishubot/services"
	"start-feishubot/utils"
)

type MsgInfo struct {
	handlerType HandlerType
	msgType     string
	msgId       *string
	chatId      *string
	qParsed     string
	sessionId   *string
	mention     []*larkim.MentionEvent
}
type ActionInfo struct {
	handler *MessageHandler
	ctx     *context.Context
	info    *MsgInfo
}

type Action interface {
	Execute(a *ActionInfo) bool
}

type ProcessedUniqueAction struct { //消息唯一性
}

func (*ProcessedUniqueAction) Execute(a *ActionInfo) bool {
	if a.handler.msgCache.IfProcessed(*a.info.msgId) {
		return false
	}
	a.handler.msgCache.TagProcessed(*a.info.msgId)
	return true
}

type ProcessMentionAction struct { //是否机器人应该处理
}

func (*ProcessMentionAction) Execute(a *ActionInfo) bool {
	// 私聊直接过
	if a.info.handlerType == UserHandler {
		return true
	}
	// 群聊判断是否提到机器人
	if a.info.handlerType == GroupHandler {
		if a.handler.judgeIfMentionMe(a.info.mention) {
			return true
		}
		return false
	}
	return false
}

type EmptyAction struct { /*空消息*/
}

func (*EmptyAction) Execute(a *ActionInfo) bool {
	if len(a.info.qParsed) == 0 {
		sendMsg(*a.ctx, "🤖️：你想知道什么呢~", a.info.chatId)
		fmt.Println("msgId", *a.info.msgId,
			"message.text is empty")
		return false
	}
	return true
}

type ClearAction struct { /*清除消息*/
}

func (*ClearAction) Execute(a *ActionInfo) bool {
	if _, foundClear := utils.EitherTrimEqual(a.info.qParsed,
		"/clear", "清除"); foundClear {
		sendClearCacheCheckCard(*a.ctx, a.info.sessionId,
			a.info.msgId)
		return false
	}
	return true
}

type RolePlayAction struct { /*角色扮演*/
}

func (*RolePlayAction) Execute(a *ActionInfo) bool {
	if system, foundSystem := utils.EitherCutPrefix(a.info.qParsed,
		"/system ", "角色扮演 "); foundSystem {
		a.handler.sessionCache.Clear(*a.info.sessionId)
		systemMsg := append([]services.Messages{}, services.Messages{
			Role: "system", Content: system,
		})
		a.handler.sessionCache.SetMsg(*a.info.sessionId, systemMsg)
		sendSystemInstructionCard(*a.ctx, a.info.sessionId,
			a.info.msgId, system)
		return false
	}
	return true
}

type HelpAction struct { /*帮助*/
}

func (*HelpAction) Execute(a *ActionInfo) bool {
	if _, foundHelp := utils.EitherTrimEqual(a.info.qParsed, "/help",
		"帮助"); foundHelp {
		sendHelpCard(*a.ctx, a.info.sessionId, a.info.msgId)
		return false
	}
	return true
}

type PicAction struct { /*图片*/
}

func (*PicAction) Execute(a *ActionInfo) bool {
	// 开启图片创作模式
	if _, foundPic := utils.EitherTrimEqual(a.info.qParsed,
		"/picture", "图片创作"); foundPic {
		a.handler.sessionCache.Clear(*a.info.sessionId)
		a.handler.sessionCache.SetMode(*a.info.sessionId,
			services.ModePicCreate)
		sendPicCreateInstructionCard(*a.ctx, a.info.sessionId,
			a.info.msgId)
		return false
	}

	// 生成图片
	mode := a.handler.sessionCache.GetMode(*a.info.sessionId)
	if mode == services.ModePicCreate {
		bs64, err := a.handler.gpt.GenerateOneImage(a.info.qParsed,
			"256x256")
		if err != nil {
			replyMsg(*a.ctx, fmt.Sprintf(
				"🤖️：图片生成失败，请稍后再试～\n错误信息: %v", err), a.info.msgId)
			return false
		}
		replayImageByBase64(*a.ctx, bs64, a.info.msgId)
		return false
	}

	return true
}

type MessageAction struct { /*消息*/
}

func (*MessageAction) Execute(a *ActionInfo) bool {
	msg := a.handler.sessionCache.GetMsg(*a.info.sessionId)
	msg = append(msg, services.Messages{
		Role: "user", Content: a.info.qParsed,
	})
	completions, err := a.handler.gpt.Completions(msg)
	if err != nil {
		replyMsg(*a.ctx, fmt.Sprintf(
			"🤖️：消息机器人摆烂了，请稍后再试～\n错误信息: %v", err), a.info.msgId)
		return false
	}
	msg = append(msg, completions)
	a.handler.sessionCache.SetMsg(*a.info.sessionId, msg)
	//if new topic
	if len(msg) == 2 {
		fmt.Println("new topic", msg[1].Content)
		sendNewTopicCard(*a.ctx, a.info.sessionId, a.info.msgId,
			completions.Content)
		return false
	}
	err = replyMsg(*a.ctx, completions.Content, a.info.msgId)
	if err != nil {
		replyMsg(*a.ctx, fmt.Sprintf(
			"🤖️：消息机器人摆烂了，请稍后再试～\n错误信息: %v", err), a.info.msgId)
		return false
	}
	return true
}
