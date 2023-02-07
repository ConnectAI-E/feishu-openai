package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	"regexp"
	"start-feishubot/initialization"
	"strings"
)

func sendMsg(ctx context.Context, msg string, chatId *string) {
	msg = strings.Trim(msg, " ")
	msg = strings.Trim(msg, "\n")
	msg = strings.Trim(msg, "\r")
	msg = strings.Trim(msg, "\t")
	// 去除空行 以及空行前的空格
	regex := regexp.MustCompile(`\n[\s| ]*\r`)
	msg = regex.ReplaceAllString(msg, "\n")
	//换行符转义
	msg = strings.ReplaceAll(msg, "\n", "\\n")
	fmt.Println("sendMsg", msg, chatId)
	client := initialization.GetLarkClient()
	content := larkim.NewTextMsgBuilder().
		TextLine(msg).
		Build()
	fmt.Println("content", content)

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
	}

	// 服务端错误处理
	if !resp.Success() {
		fmt.Println(resp.Code, resp.Msg, resp.RequestId())
	}
}
func msgFilter(msg string) string {
	//replace @到下一个非空的字段 为 ''
	regex := regexp.MustCompile(`@[^ ]*`)
	return regex.ReplaceAllString(msg, "")

}
func parseContent(content string) string {
	//"{\"text\":\"@_user_1  hahaha\"}",
	//only get text content hahaha
	var contentMap map[string]interface{}
	err := json.Unmarshal([]byte(content), &contentMap)
	if err != nil {
		fmt.Println(err)
	}
	text := contentMap["text"].(string)
	return msgFilter(text)
}
