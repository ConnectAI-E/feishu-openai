package handlers

import (
	"context"
	"fmt"
	larkcard "github.com/larksuite/oapi-sdk-go/v3/card"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
)

type MessageHandlerInterface interface {
	handle(ctx context.Context, event *larkim.P2MessageReceiveV1) error
}

type HandlerType string

const (
	GroupHandler = "group"
	UserHandler  = "personal"
)

// handlers 所有消息类型类型的处理器
var handlers map[HandlerType]MessageHandlerInterface

func init() {
	handlers = make(map[HandlerType]MessageHandlerInterface)
	handlers[GroupHandler] = NewGroupMessageHandler()
	handlers[UserHandler] = NewPersonalMessageHandler()

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
	return handlers[handlerType].handle(ctx, event)
}

func CardHandler() func(ctx context.Context,
	cardAction *larkcard.CardAction) (interface{}, error) {
	return func(ctx context.Context, cardAction *larkcard.CardAction) (interface{}, error) {
		value := cardAction.Action.Value
		//change map to struct CardMsg

		fmt.Println(value)
		return nil, nil
	}
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
