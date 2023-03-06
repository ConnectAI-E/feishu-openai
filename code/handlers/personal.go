package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"start-feishubot/services"
	"start-feishubot/utils"
	"strings"

	larkcard "github.com/larksuite/oapi-sdk-go/v3/card"

	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
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

//消息唯一性
type ProcessedAction struct {
}

func (*ProcessedAction) Execute(data *ActionInfo) bool {
	if data.p.msgCache.IfProcessed(*data.msgId) {
		data.p.msgCache.TagProcessed(*data.msgId)
		return false
	}
	return true
}

//空消息
type EmptyAction struct {
}

func (*EmptyAction) Execute(data *ActionInfo) bool {
	if len(data.qParsed) != 0 {
		sendMsg(*data.ctx, "🤖️：你想知道什么呢~", data.chatId)
		fmt.Println("msgId", *data.msgId, "message.text is empty")
		return false
	}
	return true
}

//清除消息
type ClearAction struct {
}

func (*ClearAction) Execute(data *ActionInfo) bool {
	if _, foundClear := utils.EitherTrimEqual(data.qParsed, "/clear", "清除"); foundClear {
		sendClearCacheCheckCard(*data.ctx, data.sessionId, data.msgId)
		return false
	}
	return true
}

//角色扮演
type RolePlayAction struct {
}

func (*RolePlayAction) Execute(data *ActionInfo) bool {
	if system, foundSystem := utils.EitherCutPrefix(data.qParsed, "/system ", "角色扮演 "); foundSystem {
		data.p.sessionCache.Clear(*data.sessionId)
		systemMsg := append([]services.Messages{}, services.Messages{
			Role: "system", Content: system,
		})
		data.p.sessionCache.Set(*data.sessionId, systemMsg)
		sendSystemInstructionCard(*data.ctx, data.sessionId, data.msgId, system)
		return false
	}
	return true
}

//帮助
type HelpAction struct {
}

func (*HelpAction) Execute(data *ActionInfo) bool {
	if _, foundHelp := utils.EitherTrimEqual(data.qParsed, "/help", "帮助"); foundHelp {
		sendHelpCard(*data.ctx, data.sessionId, data.msgId)
		return false
	}
	return true
}

type MessageAction struct {
}

func (*MessageAction) Execute(data *ActionInfo) bool {
	msg := data.p.sessionCache.Get(*data.sessionId)
	msg = append(msg, services.Messages{
		Role: "user", Content: data.qParsed,
	})
	completions, err := data.p.gpt.Completions(msg)
	if err != nil {
		replyMsg(*data.ctx, fmt.Sprintf("🤖️：消息机器人摆烂了，请稍后再试～\n错误信息: %v", err), data.msgId)
		return false
	}
	msg = append(msg, completions)
	p.sessionCache.Set(*data.sessionId, msg)
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

//责任链
func chain(data *ActionInfo, actions ...Action) bool {
	for _, v := range actions {
		if !v.Execute(data) {
			return false
		}
	}
	return true
}

type PersonalMessageHandler struct {
	sessionCache services.SessionServiceCacheInterface
	msgCache     services.MsgCacheInterface
	gpt          services.ChatGPT
}

func (p PersonalMessageHandler) cardHandler(
	_ context.Context,
	cardAction *larkcard.CardAction) (interface{}, error) {
	var cardMsg CardMsg
	actionValue := cardAction.Action.Value
	actionValueJson, _ := json.Marshal(actionValue)
	json.Unmarshal(actionValueJson, &cardMsg)
	if cardMsg.Kind == ClearCardKind {
		newCard, err, done := CommonProcessClearCache(cardMsg, p.sessionCache)
		if done {
			return newCard, err
		}
	}
	return nil, nil
}

func CommonProcessClearCache(cardMsg CardMsg, session services.SessionServiceCacheInterface) (interface{},
	error,
	bool) {
	if cardMsg.Value == "1" {
		newCard, _ := newSendCard(
			withHeader("️🆑 机器人提醒", larkcard.TemplateRed),
			withMainMd("已删除此话题的上下文信息"),
			withNote("我们可以开始一个全新的话题，继续找我聊天吧"),
		)
		session.Clear(cardMsg.SessionId)
		return newCard, nil, true
	}
	if cardMsg.Value == "0" {
		newCard, _ := newSendCard(
			withHeader("️🆑 机器人提醒", larkcard.TemplateGreen),
			withMainMd("依旧保留此话题的上下文信息"),
			withNote("我们可以继续探讨这个话题,期待和您聊天。如果您有其他问题或者想要讨论的话题，请告诉我哦"),
		)
		return newCard, nil, true
	}
	return nil, nil, false
}

func (p PersonalMessageHandler) handle(ctx context.Context, event *larkim.P2MessageReceiveV1) error {
	content := event.Event.Message.Content
	msgId := event.Event.Message.MessageId
	rootId := event.Event.Message.RootId
	chatId := event.Event.Message.ChatId
	sessionId := rootId
	if sessionId == nil || *sessionId == "" {
		sessionId = msgId
	}
	//责任链重构示例
	data := &ActionInfo{
		p:         &p,
		msgId:     msgId,
		qParsed:   strings.Trim(parseContent(*content), " "),
		ctx:       &ctx,
		chatId:    chatId,
		sessionId: sessionId,
	}
	actions := []Action{
		&ProcessedAction{}, //唯一处理
		&EmptyAction{},     //空消息处理
		&ClearAction{},     //清除消息处理
		&RolePlayAction{},  //角色扮演处理
		&MessageAction{},   //消息处理
	}
	chain(data, actions...)
	return nil

}

var _ MessageHandlerInterface = (*PersonalMessageHandler)(nil)

func NewPersonalMessageHandler(gpt services.ChatGPT) MessageHandlerInterface {
	return &PersonalMessageHandler{
		sessionCache: services.GetSessionCache(),
		msgCache:     services.GetMsgCache(),
		gpt:          gpt,
	}
}
