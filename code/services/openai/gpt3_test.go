package openai

import (
	"fmt"
	"testing"

	"start-feishubot/initialization"
)

func TestCompletions(t *testing.T) {
	config := initialization.LoadConfig("../../config.yaml")

	msgs := []Messages{
		{Role: "system", Content: "你是一个专业的翻译官，负责中英文翻译。"},
		{Role: "user", Content: "翻译这段话: The assistant messages help store prior responses. They can also be written by a developer to help give examples of desired behavior."},
	}

	gpt := NewChatGPT(*config)

	resp, err := gpt.Completions(msgs, Balance)
	if err != nil {
		t.Errorf("TestCompletions failed with error: %v", err)
	}

	fmt.Println(resp.Content, resp.Role)
}

func TestGenerateOneImage(t *testing.T) {
	config := initialization.LoadConfig("../../config.yaml")
	gpt := NewChatGPT(*config)
	prompt := "a red apple"
	size := "256x256"
	imageURL, err := gpt.GenerateOneImage(prompt, size)
	if err != nil {
		t.Errorf("TestGenerateOneImage failed with error: %v", err)
	}
	if imageURL == "" {
		t.Errorf("TestGenerateOneImage returned empty imageURL")
	}
}

func TestAudioToText(t *testing.T) {
	config := initialization.LoadConfig("../../config.yaml")
	gpt := NewChatGPT(*config)
	audio := "./test_file/test.wav"
	text, err := gpt.AudioToText(audio)
	if err != nil {
		t.Errorf("TestAudioToText failed with error: %v", err)
	}
	fmt.Printf("TestAudioToText returned text: %s \n", text)
	if text == "" {
		t.Errorf("TestAudioToText returned empty text")
	}

}

func TestVariateOneImage(t *testing.T) {
	config := initialization.LoadConfig("../../config.yaml")
	gpt := NewChatGPT(*config)
	image := "./test_file/img.png"
	size := "256x256"
	//compressionType, err := GetImageCompressionType(image)
	//if err != nil {
	//	return
	//}
	//fmt.Println("compressionType: ", compressionType)
	ConvertToRGBA(image, image)
	err := VerifyPngs([]string{image})
	if err != nil {
		t.Errorf("TestVariateOneImage failed with error: %v", err)
		return
	}

	imageBs64, err := gpt.GenerateOneImageVariation(image, size)
	if err != nil {
		t.Errorf("TestVariateOneImage failed with error: %v", err)
	}
	//fmt.Printf("TestVariateOneImage returned imageBs64: %s \n", imageBs64)
	if imageBs64 == "" {
		t.Errorf("TestVariateOneImage returned empty imageURL")
	}
}

func TestVariateOneImageWithJpg(t *testing.T) {
	config := initialization.LoadConfig("../../config.yaml")
	gpt := NewChatGPT(*config)
	image := "./test_file/test.jpg"
	size := "256x256"
	compressionType, err := GetImageCompressionType(image)
	if err != nil {
		return
	}
	fmt.Println("compressionType: ", compressionType)
	//ConvertJPGtoPNG(image)
	ConvertToRGBA(image, image)
	err = VerifyPngs([]string{image})
	if err != nil {
		t.Errorf("TestVariateOneImage failed with error: %v", err)
		return
	}

	imageBs64, err := gpt.GenerateOneImageVariation(image, size)
	if err != nil {
		t.Errorf("TestVariateOneImage failed with error: %v", err)
	}
	fmt.Printf("TestVariateOneImage returned imageBs64: %s \n", imageBs64)
	if imageBs64 == "" {
		t.Errorf("TestVariateOneImage returned empty imageURL")
	}
}

// 余额接口已经被废弃
func TestChatGPT_GetBalance(t *testing.T) {
	config := initialization.LoadConfig("../../config.yaml")
	gpt := NewChatGPT(*config)
	balance, err := gpt.GetBalance()
	if err != nil {
		t.Errorf("TestChatGPT_GetBalance failed with error: %v", err)
	}
	fmt.Println("balance: ", balance)
}
