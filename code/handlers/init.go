package handlers

import (
	"context"
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
	//handlers[GroupHandler] = NewGroupMessageHandler()
	handlers[UserHandler] = NewPersonalMessageHandler()
}

func Handler(ctx context.Context, event *larkim.P2MessageReceiveV1) error {
	return handlers[UserHandler].handle(ctx, event)
}
