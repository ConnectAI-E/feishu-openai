package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"start-feishubot/initialization"
	"start-feishubot/services"
	"start-feishubot/utils"
	"strings"

	larkcard "github.com/larksuite/oapi-sdk-go/v3/card"

	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
)

type GroupMessageHandler struct {
	sessionCache services.SessionServiceCacheInterface
	msgCache     services.MsgCacheInterface
	gpt          services.ChatGPT
	config       initialization.Config
}

func (p GroupMessageHandler) cardHandler(_ context.Context,
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

func (p GroupMessageHandler) msgReceivedHandler(ctx context.Context, event *larkim.P2MessageReceiveV1) error {
	ifMention := p.judgeIfMentionMe(event)
	if !ifMention {
		return nil
	}
	content := event.Event.Message.Content
	msgId := event.Event.Message.MessageId
	rootId := event.Event.Message.RootId
	chatId := event.Event.Message.ChatId
	sessionId := rootId
	if sessionId == nil || *sessionId == "" {
		sessionId = msgId
	}

	if p.msgCache.IfProcessed(*msgId) {
		fmt.Println("msgId", *msgId, "processed")
		return nil
	}
	p.msgCache.TagProcessed(*msgId)
	qParsed := strings.Trim(parseContent(*content), " ")
	if len(qParsed) == 0 {
		sendMsg(ctx, "🤖️：你想知道什么呢~", chatId)
		fmt.Println("msgId", *msgId, "message.text is empty")
		return nil
	}

	if _, foundClear := utils.EitherTrimEqual(qParsed, "/clear", "清除"); foundClear {
		sendClearCacheCheckCard(ctx, sessionId, msgId)
		return nil
	}

	if system, foundSystem := utils.EitherCutPrefix(qParsed, "/system ", "角色扮演 "); foundSystem {
		p.sessionCache.Clear(*sessionId)
		systemMsg := append([]services.Messages{}, services.Messages{
			Role: "system", Content: system,
		})
		p.sessionCache.SetMsg(*sessionId, systemMsg)
		sendSystemInstructionCard(ctx, sessionId, msgId, system)
		return nil
	}

	if _, foundHelp := utils.EitherTrimEqual(qParsed, "/help", "帮助"); foundHelp {
		sendHelpCard(ctx, sessionId, msgId)
		return nil
	}

	if _, foundPicture := utils.EitherTrimEqual(qParsed,
		"/picture", "图片创作"); foundPicture {
		p.sessionCache.Clear(*sessionId)
		p.sessionCache.SetMode(*sessionId, services.ModePicCreate)
		sendPicCreateInstructionCard(ctx, sessionId, msgId)
		return nil
	}
	mode := p.sessionCache.GetMode(*sessionId)
	if mode == services.ModePicCreate {
		bs64, err := p.gpt.GenerateOneImage(qParsed, "256x256")
		if err != nil {
			replyMsg(ctx, fmt.Sprintf("🤖️：图片生成失败，请稍后再试～\n错误信息: %v", err), msgId)
			return nil
		}
		replayImageByBase64(ctx, bs64, msgId)
		return nil
	}

	msg := p.sessionCache.GetMsg(*sessionId)
	msg = append(msg, services.Messages{
		Role: "user", Content: qParsed,
	})
	completions, err := p.gpt.Completions(msg)
	if err != nil {
		replyMsg(ctx, fmt.Sprintf("🤖️：消息机器人摆烂了，请稍后再试～\n错误信息: %v", err), msgId)
		return nil
	}
	msg = append(msg, completions)
	p.sessionCache.SetMsg(*sessionId, msg)
	if len(msg) == 2 {
		sendNewTopicCard(ctx, sessionId, msgId, completions.Content)
		return nil
	}
	err = replyMsg(ctx, completions.Content, msgId)
	if err != nil {
		replyMsg(ctx, fmt.Sprintf("🤖️：消息机器人摆烂了，请稍后再试～\n错误信息: %v", err), msgId)
		return nil
	}
	return nil

}

var _ MessageHandlerInterface = (*GroupMessageHandler)(nil)

func NewGroupMessageHandler(gpt services.ChatGPT, config initialization.Config) MessageHandlerInterface {
	return &GroupMessageHandler{
		sessionCache: services.GetSessionCache(),
		msgCache:     services.GetMsgCache(),
		gpt:          gpt,
		config:       config,
	}
}

func (p GroupMessageHandler) judgeIfMentionMe(event *larkim.P2MessageReceiveV1) bool {
	mention := event.Event.Message.Mentions
	if len(mention) != 1 {
		return false
	}
	return *mention[0].Name == p.config.FeishuBotName
}
