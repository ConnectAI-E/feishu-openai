package handlers

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/larksuite/oapi-sdk-go/v3/card"
	"github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	"start-feishubot/initialization"
)

type CardKind string
type CardChatType string

var (
	ClearCardKind = CardKind("clear")
)

var (
	GroupChatType = CardChatType("group")
	UserChatType  = CardChatType("personal")
)

type CardMsg struct {
	Kind      CardKind
	ChatType  CardChatType
	Value     interface{}
	SessionId string
}

func replyCard(ctx context.Context,
	msgId *string,
	cardContent string,
) error {
	client := initialization.GetLarkClient()
	resp, err := client.Im.Message.Reply(ctx, larkim.NewReplyMessageReqBuilder().
		MessageId(*msgId).
		Body(larkim.NewReplyMessageReqBodyBuilder().
			MsgType(larkim.MsgTypeInteractive).
			Uuid(uuid.New().String()).
			Content(cardContent).
			Build()).
		Build())

	// 处理错误
	if err != nil {
		fmt.Println(err)
		return err
	}

	// 服务端错误处理
	if !resp.Success() {
		fmt.Println(resp.Code, resp.Msg, resp.RequestId())
		return err
	}
	return nil
}

func newSendCard(
	header *larkcard.MessageCardHeader,
	elements ...larkcard.MessageCardElement) (string,
	error) {
	config := larkcard.NewMessageCardConfig().
		WideScreenMode(false).
		EnableForward(true).
		UpdateMulti(false).
		Build()
	var aElementPool []larkcard.MessageCardElement
	for _, element := range elements {
		aElementPool = append(aElementPool, element)
	}
	// 卡片消息体
	cardContent, err := larkcard.NewMessageCard().
		Config(config).
		Header(header).
		Elements(
			aElementPool,
		).
		String()
	return cardContent, err
}

func withHeader(title string, color string) *larkcard.
	MessageCardHeader {
	if title == "" {
		title = "🤖️机器人提醒"
	}
	header := larkcard.NewMessageCardHeader().
		Template(color).
		Title(larkcard.NewMessageCardPlainText().
			Content(title).
			Build()).
		Build()
	return header
}
func withNote(note string) larkcard.MessageCardElement {
	noteElement := larkcard.NewMessageCardNote().
		Elements([]larkcard.MessageCardNoteElement{larkcard.NewMessageCardPlainText().
			Content(note).
			Build()}).
		Build()
	return noteElement
}

func withMainMd(msg string) larkcard.MessageCardElement {
	msg, i := processMessage(msg)
	msg = processNewLine(msg)
	if i != nil {
		return nil
	}
	mainElement := larkcard.NewMessageCardDiv().
		Fields([]*larkcard.MessageCardField{larkcard.NewMessageCardField().
			Text(larkcard.NewMessageCardLarkMd().
				Content(msg).
				Build()).
			IsShort(true).
			Build()}).
		Build()
	return mainElement
}
func withMainText(msg string) larkcard.MessageCardElement {
	msg, i := processMessage(msg)
	msg = cleanTextBlock(msg)
	if i != nil {
		return nil
	}
	mainElement := larkcard.NewMessageCardDiv().
		Fields([]*larkcard.MessageCardField{larkcard.NewMessageCardField().
			Text(larkcard.NewMessageCardPlainText().
				Content(msg).
				Build()).
			IsShort(false).
			Build()}).
		Build()
	return mainElement
}

func withDoubleCheckBtn(sessionId *string) larkcard.
	MessageCardElement {
	actions := larkcard.NewMessageCardAction().
		Actions([]larkcard.MessageCardActionElement{
			larkcard.NewMessageCardEmbedButton().
				Type(larkcard.MessageCardButtonTypeDanger).
				Value(map[string]interface{}{
					"value":     "1", // 1 代表确认清除
					"kind":      ClearCardKind,
					"chatType":  UserChatType,
					"sessionId": *sessionId,
				}).
				Text(larkcard.NewMessageCardPlainText().
					Content("确认清除").
					Build()),
			larkcard.NewMessageCardEmbedButton().
				Type(larkcard.MessageCardButtonTypePrimary).
				Value(map[string]interface{}{
					"value":     "0", // 0 代表取消清除
					"kind":      ClearCardKind,
					"sessionId": *sessionId,
					"chatType":  UserChatType,
				}).
				Text(larkcard.NewMessageCardPlainText().
					Content("我再想想").
					Build()),
		}).Layout(larkcard.MessageCardActionLayoutBisected.Ptr()).
		Build()
	return actions
}
func replyMsg(ctx context.Context, msg string, msgId *string) error {
	fmt.Println("sendMsg", msg, msgId)
	msg, i := processMessage(msg)
	if i != nil {
		return i
	}
	client := initialization.GetLarkClient()
	content := larkim.NewTextMsgBuilder().
		Text(msg).
		Build()

	resp, err := client.Im.Message.Reply(ctx, larkim.NewReplyMessageReqBuilder().
		MessageId(*msgId).
		Body(larkim.NewReplyMessageReqBodyBuilder().
			MsgType(larkim.MsgTypeText).
			Uuid(uuid.New().String()).
			Content(content).
			Build()).
		Build())

	// 处理错误
	if err != nil {
		fmt.Println(err)
		return err
	}

	// 服务端错误处理
	if !resp.Success() {
		fmt.Println(resp.Code, resp.Msg, resp.RequestId())
		return err
	}
	return nil
}
func sendMsg(ctx context.Context, msg string, chatId *string) error {
	//fmt.Println("sendMsg", msg, chatId)
	msg, i := processMessage(msg)
	if i != nil {
		return i
	}
	client := initialization.GetLarkClient()
	content := larkim.NewTextMsgBuilder().
		Text(msg).
		Build()

	//fmt.Println("content", content)

	resp, err := client.Im.Message.Create(ctx, larkim.NewCreateMessageReqBuilder().
		ReceiveIdType(larkim.ReceiveIdTypeChatId).
		Body(larkim.NewCreateMessageReqBodyBuilder().
			MsgType(larkim.MsgTypeText).
			ReceiveId(*chatId).
			Content(content).
			Build()).
		Build())

	// 处理错误
	if err != nil {
		fmt.Println(err)
		return err
	}

	// 服务端错误处理
	if !resp.Success() {
		fmt.Println(resp.Code, resp.Msg, resp.RequestId())
		return err
	}
	return nil
}

func sendClearCacheCheckCard(ctx context.Context,
	sessionId *string, msgId *string) {
	newCard, _ := newSendCard(
		withHeader("👻️ 机器人提醒", larkcard.TemplateBlue),
		withMainMd("您确定要清除对话上下文吗？"),
		withNote("请注意，这将开始一个全新的对话，您将无法利用之前话题的历史信息"),
		withDoubleCheckBtn(sessionId))
	replyCard(
		ctx,
		msgId,
		newCard,
	)
}

func sendSystemInstructionCard(ctx context.Context,
	sessionId *string, msgId *string, content string) {
	newCard, _ := newSendCard(
		withHeader("🥷 已进入角色扮演模式", larkcard.TemplateBlue),
		withMainText(content),
		withNote("请注意，这将开始一个全新的对话，您将无法利用之前话题的历史信息"))
	replyCard(
		ctx,
		msgId,
		newCard,
	)
}

func sendNewTopicCard(ctx context.Context,
	sessionId *string, msgId *string, content string) {
	newCard, _ := newSendCard(
		withHeader("👻️ 已开启新的话题", larkcard.TemplateBlue),
		withMainText(content),
		withNote("提醒：点击对话框参与回复，可保持话题连贯"))
	replyCard(
		ctx,
		msgId,
		newCard,
	)
}
