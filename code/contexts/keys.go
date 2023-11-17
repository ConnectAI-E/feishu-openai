package contexts

import (
	"context"
	"fmt"
	"net/url"
	"reflect"
)

type ChatContext struct {
	Tenant        string
	ChatID        string
	ChatType      string
	MessageID     string
	MessageType   string
	SenderOpenID  string
	SenderType    string
	SenderUnionID string
	SenderUserID  string
	SessionID     string
}

func (cc *ChatContext) Encode() string {
	v := url.Values{}
	v.Set("feishu.tenant", cc.Tenant)
	v.Set("feishu.session_id", cc.SessionID)
	v.Set("feishu.chat_id", cc.ChatID)
	v.Set("feishu.chat_type", cc.ChatType)
	v.Set("feishu.message_id", cc.MessageID)
	v.Set("feishu.message_type", cc.MessageType)
	v.Set("feishu.sender_user_id", cc.SenderUserID)
	v.Set("feishu.sender_union_id", cc.SenderUnionID)
	v.Set("feishu.sender_open_id", cc.SenderOpenID)
	v.Set("feishu.sender_type", cc.SenderType)
	for k, vv := range v {
		if len(vv) == 0 {
			delete(v, k)
		}
	}
	return v.Encode()
}

var ChatContextKey = CreateContextKey[*ChatContext]()

type ContextKey[T any] interface {
	Value(ctx context.Context) (T, bool)
	Get(ctx context.Context) T
	Must(ctx context.Context) T
	WithValue(ctx context.Context, val T) context.Context
}

type key[T any] struct {
	opts CreateContextKeyOptions[T]
}

func (k key[T]) Value(ctx context.Context) (T, bool) {
	o, ok := ctx.Value(k.opts.key).(T)
	return o, ok
}

func (k key[T]) Get(ctx context.Context) T {
	o, _ := ctx.Value(k.opts.key).(T)
	return o
}

func (k key[T]) Must(ctx context.Context) T {
	o, ok := ctx.Value(k.opts.key).(T)
	if !ok {
		panic(fmt.Errorf("%s not found in context", k.String()))
	}
	return o
}

func (k key[T]) WithValue(ctx context.Context, val T) context.Context {
	return context.WithValue(ctx, k.opts.key, val)
}

func (k key[T]) String() string {
	name := k.opts.Name
	if name != "" {
		name = "@" + name
	}
	return fmt.Sprintf("ContextKey(%s%s)", reflect.TypeOf(new(T)).Elem().String(), name)
}

var _ ContextKey[string] = (*key[string])(nil)

type CreateContextKeyOptions[T any] struct {
	Name string
	key  any
}

func CreateContextKey[T any](opts ...CreateContextKeyOptions[T]) ContextKey[T] {
	var opt CreateContextKeyOptions[T]
	if len(opts) > 0 {
		// reduce
		for _, o := range opts {
			opt = o
		}
	}
	if opt.Name != "" {
		opt.key = opt.Name
	} else {
		opt.key = reflect.TypeOf(new(T)).Elem()
	}
	return &key[T]{opts: opt}
}
