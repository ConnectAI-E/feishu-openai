package handlers

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// func sendCard
func msgFilter(msg string) string {
	//replace @到下一个非空的字段 为 ''
	regex := regexp.MustCompile(`@[^ ]*`)
	return regex.ReplaceAllString(msg, "")
}

// Parse rich text json to text
func parsePostContent(content string) string {
	var contentMap map[string]interface{}
	err := json.Unmarshal([]byte(content), &contentMap)

	if err != nil {
		fmt.Println(err)
	}

	if contentMap["content"] == nil {
		return ""
	}
	var text string
	// deal with title
	if contentMap["title"] != nil && contentMap["title"] != "" {
		text += contentMap["title"].(string) + "\n"
	}
	// deal with content
	contentList := contentMap["content"].([]interface{})
	for _, v := range contentList {
		for _, v1 := range v.([]interface{}) {
			if v1.(map[string]interface{})["tag"] == "text" {
				text += v1.(map[string]interface{})["text"].(string)
			}
		}
		// add new line
		text += "\n"
	}
	return msgFilter(text)
}

func parsePostImageKeys(content string) []string {
	var contentMap map[string]interface{}
	err := json.Unmarshal([]byte(content), &contentMap)

	if err != nil {
		fmt.Println(err)
		return nil
	}

	var imageKeys []string

	if contentMap["content"] == nil {
		return imageKeys
	}

	contentList := contentMap["content"].([]interface{})
	for _, v := range contentList {
		for _, v1 := range v.([]interface{}) {
			if v1.(map[string]interface{})["tag"] == "img" {
				imageKeys = append(imageKeys, v1.(map[string]interface{})["image_key"].(string))
			}
		}
	}

	return imageKeys
}

func parseContent(content, msgType string) string {
	//"{\"text\":\"@_user_1  hahaha\"}",
	//only get text content hahaha
	if msgType == "post" {
		return parsePostContent(content)
	}

	var contentMap map[string]interface{}
	err := json.Unmarshal([]byte(content), &contentMap)
	if err != nil {
		fmt.Println(err)
	}
	if contentMap["text"] == nil {
		return ""
	}
	text := contentMap["text"].(string)
	return msgFilter(text)
}

func processMessage(msg interface{}) (string, error) {
	msg = strings.TrimSpace(msg.(string))
	msgB, err := json.Marshal(msg)
	if err != nil {
		return "", err
	}

	msgStr := string(msgB)

	if len(msgStr) >= 2 {
		msgStr = msgStr[1 : len(msgStr)-1]
	}
	return msgStr, nil
}

func processNewLine(msg string) string {
	return strings.Replace(msg, "\\n", `
`, -1)
}

func processQuote(msg string) string {
	return strings.Replace(msg, "\\\"", "\"", -1)
}

// 将字符中 \u003c 替换为 <  等等
func processUnicode(msg string) string {
	regex := regexp.MustCompile(`\\u[0-9a-fA-F]{4}`)
	return regex.ReplaceAllStringFunc(msg, func(s string) string {
		r, _ := regexp.Compile(`\\u`)
		s = r.ReplaceAllString(s, "")
		i, _ := strconv.ParseInt(s, 16, 32)
		return string(rune(i))
	})
}

func cleanTextBlock(msg string) string {
	msg = processNewLine(msg)
	msg = processUnicode(msg)
	msg = processQuote(msg)
	return msg
}

func parseFileKey(content string) string {
	var contentMap map[string]interface{}
	err := json.Unmarshal([]byte(content), &contentMap)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	if contentMap["file_key"] == nil {
		return ""
	}
	fileKey := contentMap["file_key"].(string)
	return fileKey
}

func parseImageKey(content string) string {
	var contentMap map[string]interface{}
	err := json.Unmarshal([]byte(content), &contentMap)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	if contentMap["image_key"] == nil {
		return ""
	}
	imageKey := contentMap["image_key"].(string)
	return imageKey
}
