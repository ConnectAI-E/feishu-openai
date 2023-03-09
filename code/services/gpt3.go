package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"sync"
	"time"
)

const (
	BASEURL     = "https://api.openai.com/v1/"
	maxTokens   = 2000
	temperature = 0.7
	engine      = "gpt-3.5-turbo"
)

// ChatGPTResponseBody 请求体
type ChatGPTResponseBody struct {
	ID      string                 `json:"id"`
	Object  string                 `json:"object"`
	Created int                    `json:"created"`
	Model   string                 `json:"model"`
	Choices []ChoiceItem           `json:"choices"`
	Usage   map[string]interface{} `json:"usage"`
}
type ChoiceItem struct {
	Message      Messages `json:"message"`
	Index        int      `json:"index"`
	FinishReason string   `json:"finish_reason"`
}

type Messages struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatGPTRequestBody 响应体
type ChatGPTRequestBody struct {
	Model            string     `json:"model"`
	Messages         []Messages `json:"messages"`
	MaxTokens        int        `json:"max_tokens"`
	Temperature      float32    `json:"temperature"`
	TopP             int        `json:"top_p"`
	FrequencyPenalty int        `json:"frequency_penalty"`
	PresencePenalty  int        `json:"presence_penalty"`
}
type ChatGPT struct {
	ApiKeys            []string
	currentApiKeyIndex int
	apiKeyUsage        map[string]int
	apiKeyWeights      []int
	apiKeyUsageMutex   sync.Mutex
}

func (gpt *ChatGPT) Completions(msg []Messages) (resp Messages, err error) {
	// 检查当前apikey的调用次数是否超过限制，如果超过限制则重新计算权重
	currentApiKey := gpt.ApiKeys[gpt.currentApiKeyIndex]
	gpt.apiKeyUsageMutex.Lock()
	if gpt.apiKeyUsage[currentApiKey] >= 3 {
		gpt.calculateApiKeyWeights()
	}
	// 构造请求体
	requestBody := ChatGPTRequestBody{
		Model:            engine,
		Messages:         msg,
		MaxTokens:        maxTokens,
		Temperature:      temperature,
		TopP:             1,
		FrequencyPenalty: 0,
		PresencePenalty:  0,
	}
	requestData, err := json.Marshal(requestBody)
	if err != nil {
		gpt.apiKeyUsageMutex.Unlock()
		return resp, err
	}
	log.Printf("request gtp json string : %v", string(requestData))
	req, err := http.NewRequest("POST", BASEURL+"chat/completions", bytes.NewBuffer(requestData))
	if err != nil {
		gpt.apiKeyUsageMutex.Unlock()
		return resp, err
	}
	req.Header.Set("Content-Type", "application/json")
	// 获取apikey并设置到header中
	apiKey := gpt.getApiKey()
	req.Header.Set("Authorization", "Bearer "+apiKey)
	client := &http.Client{Timeout: 110 * time.Second}
	response, err := client.Do(req)
	if err != nil {
		// 请求出错时，将当前apikey的调用次数减1
		gpt.apiKeyUsage[currentApiKey]--
		gpt.apiKeyUsageMutex.Unlock()
		return resp, err
	}
	defer response.Body.Close()
	if response.StatusCode/2 != 100 {
		// 请求出错时，将当前apikey的调用次数减1
		gpt.apiKeyUsage[currentApiKey]--
		gpt.apiKeyUsageMutex.Unlock()
		return resp, fmt.Errorf("gtp api %s", response.Status)
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		gpt.apiKeyUsageMutex.Unlock()
		return resp, err
	}
	gptResponseBody := &ChatGPTResponseBody{}
	err = json.Unmarshal(body, gptResponseBody)
	if err != nil {
		gpt.apiKeyUsageMutex.Unlock()
		return resp, err
	}
	resp = gptResponseBody.Choices[0].Message
	// 当前apikey的调用次数加1
	gpt.apiKeyUsage[currentApiKey]++
	gpt.apiKeyUsageMutex.Unlock()
	// 更新ApiKeyUsage并写入文件
	err = gpt.updateApiKeyUsage()
	if err != nil {
		log.Printf("update api key usage failed: %v", err)
	}
	return resp, nil
}

func (gpt *ChatGPT) getApiKey() string {
	// 加权随机选择可用apikey
	totalWeight := 0
	for _, weight := range gpt.apiKeyWeights {
		totalWeight += weight
	}
	randNum := rand.Intn(totalWeight)
	for i, weight := range gpt.apiKeyWeights {
		if randNum < weight {
			gpt.currentApiKeyIndex = i
			break
		}
		randNum -= weight
	}
	return gpt.ApiKeys[gpt.currentApiKeyIndex]
}

func (gpt *ChatGPT) calculateApiKeyWeights() {
	// 计算apikey的权重
	totalUsage := 0
	for _, usage := range gpt.apiKeyUsage {
		totalUsage += usage
	}
	gpt.apiKeyWeights = make([]int, len(gpt.ApiKeys))
	for i, apiKey := range gpt.ApiKeys {
		usage, ok := gpt.apiKeyUsage[apiKey]
		if ok {
			if usage >= 0 {
				gpt.apiKeyWeights[i] = totalUsage - usage + 1
			} else {
				gpt.apiKeyWeights[i] = 0
			}
		}
	}
}

func (gpt *ChatGPT) initApiKeyUsage() {
	// 初始化apikey的调用次数和权重
	gpt.apiKeyUsage = make(map[string]int)
	gpt.apiKeyWeights = make([]int, len(gpt.ApiKeys))
	for i := range gpt.ApiKeys {
		gpt.apiKeyWeights[i] = 1
	}
}

func (gpt *ChatGPT) checkApiKeyAvailability() {
	for _, apiKey := range gpt.ApiKeys {
		// 检查apikey的可用性
		req, err := http.NewRequest("GET", BASEURL+"engines", nil)
		if err != nil {
			log.Printf("check api key %s failed: %v", apiKey, err)
			continue
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+apiKey)
		client := &http.Client{Timeout: 10 * time.Second}
		response, err := client.Do(req)
		if err != nil {
			log.Printf("check api key %s failed: %v", apiKey, err)
			continue
		}
		// 如果apikey不可用，则将其调用次数设为-1
		if response.StatusCode/2 != 100 {
			log.Printf("api key %s is not available: %s", apiKey, response.Status)
			gpt.apiKeyUsageMutex.Lock()
			gpt.apiKeyUsage[apiKey] = -1
			gpt.apiKeyUsageMutex.Unlock()
		} else {
			// 如果apikey可用，则将其调用次数设为0
			gpt.apiKeyUsageMutex.Lock()
			gpt.apiKeyUsage[apiKey] = 0
			gpt.apiKeyUsageMutex.Unlock()
		}
	}
}

func (gpt *ChatGPT) loadApiKeyUsage() error {
	// 从本地文件中读取apikey的调用次数和权重信息
	file, err := os.Open("apikey_usage.json")
	if err != nil {
		if os.IsNotExist(err) {
			gpt.initApiKeyUsage()
			return nil
		}
		return err
	}
	defer file.Close()
	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}
	if fileInfo.Size() == 0 {
		gpt.initApiKeyUsage()
		return nil
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&gpt.apiKeyUsage)
	if err != nil {
		return err
	}
	gpt.calculateApiKeyWeights()
	return nil
}

func (gpt *ChatGPT) saveApiKeyUsage() error {
	// 将apikey的调用次数和权重信息输出到本地文件
	file, err := os.Create("apikey_usage.json")
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	err = encoder.Encode(gpt.apiKeyUsage)
	if err != nil {
		return err
	}
	return nil
}

func (gpt *ChatGPT) updateApiKeyUsage() error {
	// 更新ApiKeyUsage并写入文件
	gpt.apiKeyUsageMutex.Lock()
	defer gpt.apiKeyUsageMutex.Unlock()
	gpt.apiKeyUsage[gpt.ApiKeys[gpt.currentApiKeyIndex]]++
	err := gpt.saveApiKeyUsage()
	if err != nil {
		return err
	}
	return nil
}

func (gpt *ChatGPT) setApiKeyAvailability() {
	// 检查apikey的可用性
	gpt.checkApiKeyAvailability()
	for apiKey, usage := range gpt.apiKeyUsage {
		// 如果apikey的调用次数小于0，则禁用它，并输出信息
		if usage < 0 {
			log.Printf("api key %s is disabled", apiKey)
			gpt.currentApiKeyIndex = (gpt.currentApiKeyIndex + 1) % len(gpt.ApiKeys)
		}
	}
}

func (gpt *ChatGPT) StartApiKeyAvailabilityCheck() {
	// 初始化apikey的调用次数和权重，检查apikey的可用性，并定时检查apikey的可用性
	err := gpt.loadApiKeyUsage()
	if err != nil {
		log.Printf("load api key usage failed: %v", err)
	}
	gpt.setApiKeyAvailability()
	ticker := time.NewTicker(30 * time.Minute)
	go func() {
		for {
			select {
			case <-ticker.C:
				gpt.setApiKeyAvailability()
				err := gpt.saveApiKeyUsage()
				if err != nil {
					log.Printf("save api key usage failed: %v", err)
				}
			}
		}
	}()
}

func FormatQuestion(question string) string {
	return "Answer:" + question
}

type ImageGenerationRequestBody struct {
	Prompt         string `json:"prompt"`
	N              int    `json:"n"`
	Size           string `json:"size"`
	ResponseFormat string `json:"response_format"`
}

type ImageGenerationResponseBody struct {
	Created int64 `json:"created"`
	Data    []struct {
		Base64Json string `json:"b64_json"`
	} `json:"data"`
}

func (gpt *ChatGPT) GenerateImage(prompt string, size string,
	n int) ([]string, error) {
	requestBody := ImageGenerationRequestBody{
		Prompt:         prompt,
		N:              n,
		Size:           size,
		ResponseFormat: "b64_json",
	}
	requestData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", BASEURL+"images/generations", bytes.NewBuffer(requestData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	// 获取apikey并设置到header中
	apiKey := gpt.getApiKey()
	req.Header.Set("Authorization", "Bearer "+apiKey)
	client := &http.Client{Timeout: 110 * time.Second}
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode/2 != 100 {
		return nil, fmt.Errorf("image generation api %s",
			response.Status)
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	imageResponseBody := &ImageGenerationResponseBody{}
	err = json.Unmarshal(body, imageResponseBody)
	if err != nil {
		return nil, err
	}

	var b64Pool []string
	for _, data := range imageResponseBody.Data {
		b64Pool = append(b64Pool, data.Base64Json)
	}
	return b64Pool, nil

}

func (gpt *ChatGPT) GenerateOneImage(prompt string, size string) (string, error) {
	b64s, err := gpt.GenerateImage(prompt, size, 1)
	if err != nil {
		return "", err
	}
	return b64s[0], nil
}
