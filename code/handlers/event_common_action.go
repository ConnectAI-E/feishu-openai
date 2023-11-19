package handlers

import (
	"context"
	"fmt"

	"start-feishubot/initialization"
	"start-feishubot/services/openai"
	"start-feishubot/utils"

	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
)

type MsgInfo struct {
	handlerType HandlerType
	msgType     string
	msgId       *string
	chatId      *string
	qParsed     string
	fileKey     string
	imageKey    string
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
		systemMsg := append([]openai.Messages{}, openai.Messages{
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

type BalanceAction struct { /*余额*/
}

func (*BalanceAction) Execute(a *ActionInfo) bool {
	if _, foundBalance := utils.EitherTrimEqual(a.info.qParsed,
		"/balance", "余额"); foundBalance {
		balanceResp, err := a.handler.gpt.GetBalance()
		if err != nil {
			replyMsg(*a.ctx, "查询余额失败，请稍后再试", a.info.msgId)
			return false
		}
		sendBalanceCard(*a.ctx, a.info.sessionId, *balanceResp)
		return false
	}
	return true
}

type RoleListAction struct { /*角色列表*/
}

func (*RoleListAction) Execute(a *ActionInfo) bool {
	if _, foundSystem := utils.EitherTrimEqual(a.info.qParsed,
		"/roles", "角色列表"); foundSystem {
		//a.handler.sessionCache.Clear(*a.info.sessionId)
		//systemMsg := append([]openai.Messages{}, openai.Messages{
		//	Role: "system", Content: system,
		//})
		//a.handler.sessionCache.SetMsg(*a.info.sessionId, systemMsg)
		//sendSystemInstructionCard(*a.ctx, a.info.sessionId,
		//	a.info.msgId, system)
		tags := initialization.GetAllUniqueTags()
		SendRoleTagsCard(*a.ctx, a.info.sessionId, a.info.msgId, *tags)
		return false
	}
	return true
}

type AIModeAction struct { /*发散模式*/
}

func (*AIModeAction) Execute(a *ActionInfo) bool {
	if _, foundMode := utils.EitherCutPrefix(a.info.qParsed,
		"/ai_mode", "发散模式"); foundMode {
		SendAIModeListsCard(*a.ctx, a.info.sessionId, a.info.msgId, openai.AIModeStrs)
		return false
	}
	return true
}
