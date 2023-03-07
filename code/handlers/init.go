package handlers

import (
	"context"
	"fmt"
	"start-feishubot/initialization"
	"start-feishubot/services"

	larkcard "github.com/larksuite/oapi-sdk-go/v3/card"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
)

type MessageHandlerInterface interface {
	msgReceivedHandler(ctx context.Context, event *larkim.P2MessageReceiveV1) error
	cardHandler(ctx context.Context, cardAction *larkcard.CardAction) (interface{}, error)
}

type HandlerType string

const (
	GroupHandler = "group"
	UserHandler  = "personal"
)

// handlers 所有消息类型类型的处理器
var handlers map[HandlerType]MessageHandlerInterface

func InitHandlers(gpt services.ChatGPT, config initialization.Config) {
	handlers = make(map[HandlerType]MessageHandlerInterface)
	handlers[GroupHandler] = NewGroupMessageHandler(gpt, config)
	handlers[UserHandler] = NewPersonalMessageHandler(gpt)

}

func Handler(ctx context.Context, event *larkim.P2MessageReceiveV1) error {
	handlerType := judgeChatType(event)
	if handlerType == "otherChat" {
		fmt.Println("unknown chat type")
		return nil
	}
	msgType := judgeMsgType(event)
	if msgType != "text" {
		fmt.Println("unknown msg type")
		return nil
	}
	return handlers[handlerType].msgReceivedHandler(ctx, event)
}

func ReadHandler(ctx context.Context, event *larkim.P2MessageReadV1) error {
	_ = event.Event.Reader.ReaderId.OpenId
	//fmt.Printf("msg is read by : %v \n", *readerId)
	return nil
}

func CardHandler() func(ctx context.Context,
	cardAction *larkcard.CardAction) (interface{}, error) {
	return func(ctx context.Context, cardAction *larkcard.CardAction) (interface{}, error) {
		handlerType := judgeCardType(cardAction)
		return handlers[handlerType].cardHandler(ctx, cardAction)
	}
}

func judgeCardType(cardAction *larkcard.CardAction) HandlerType {
	actionValue := cardAction.Action.Value
	chatType := actionValue["chatType"]
	fmt.Printf("chatType: %v", chatType)
	if chatType == "group" {
		return GroupHandler
	}
	if chatType == "personal" {
		return UserHandler
	}
	return "otherChat"
}

func judgeChatType(event *larkim.P2MessageReceiveV1) HandlerType {
	chatType := event.Event.Message.ChatType
	fmt.Printf("chatType: %v", *chatType)
	if *chatType == "group" {
		return GroupHandler
	}
	if *chatType == "p2p" {
		return UserHandler
	}
	return "otherChat"
}

func judgeMsgType(event *larkim.P2MessageReceiveV1) string {
	msgType := event.Event.Message.MessageType
	if *msgType == "text" {
		return "text"
	}
	return ""
}
