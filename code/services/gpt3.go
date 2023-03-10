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

type Messages struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatGPTResponseBody 请求体
type ChatGPTResponseBody struct {
	ID      string                 `json:"id"`
	Object  string                 `json:"object"`
	Created int                    `json:"created"`
	Model   string                 `json:"model"`
	Choices []ChatGPTChoiceItem    `json:"choices"`
	Usage   map[string]interface{} `json:"usage"`
}
type ChatGPTChoiceItem struct {
	Message      Messages `json:"message"`
	Index        int      `json:"index"`
	FinishReason string   `json:"finish_reason"`
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
	apiKeyUsageMutex   sync.RWMutex
}

func (gpt *ChatGPT) Completions(msg []Messages) (resp Messages, err error) {
	if len(gpt.ApiKeys) == 1 {
		// 如果只有一个apikey，则直接使用该apikey进行请求
		apiKey := gpt.ApiKeys[0]
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
			return resp, err
		}
		req, err := http.NewRequest("POST", BASEURL+"chat/completions", bytes.NewBuffer(requestData))
		if err != nil {
			return resp, err
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+apiKey)
		client := &http.Client{Timeout: 110 * time.Second}
		response, err := client.Do(req)
		if err != nil {
			return resp, err
		}
		defer response.Body.Close()
		if response.StatusCode/2 != 100 {
			return resp, fmt.Errorf("gtp api %s", response.Status)
		}
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return resp, err
		}
		gptResponseBody := &ChatGPTResponseBody{}
		err = json.Unmarshal(body, gptResponseBody)
		if err != nil {
			return resp, err
		}
		resp = gptResponseBody.Choices[0].Message
		return resp, nil
	} else {
		// 如果有多个apikey，则进行负载均衡操作
		currentApiKey := gpt.ApiKeys[gpt.currentApiKeyIndex]
		gpt.apiKeyUsageMutex.Lock()
		if gpt.apiKeyUsage[currentApiKey] >= 3 {
			gpt.calculateApiKeyWeights()
		}
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
		apiKey, err := gpt.getApiKey()
		if err != nil {
			return resp, err
		}
		req.Header.Set("Authorization", "Bearer "+apiKey)
		client := &http.Client{Timeout: 110 * time.Second}
		response, err := client.Do(req)
		if err != nil {
			gpt.apiKeyUsage[currentApiKey]--
			gpt.apiKeyUsageMutex.Unlock()
			return resp, err
		}
		defer response.Body.Close()
		if response.StatusCode/2 != 100 {
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
		gpt.apiKeyUsage[currentApiKey]++
		gpt.apiKeyUsageMutex.Unlock()
		gpt.updateApiKeyUsage(apiKey)
		return resp, nil
	}
}

func (gpt *ChatGPT) getApiKey() (string, error) {
	// 获取可用的apikey
	if len(gpt.ApiKeys) == 0 {
		return "", fmt.Errorf("no api keys provided")
	}
	if len(gpt.ApiKeys) == 1 {
		// 如果只有一个apikey，则直接返回该apikey
		return gpt.ApiKeys[0], nil
	}
	totalWeight := 0
	for _, weight := range gpt.apiKeyWeights {
		if weight <= 0 {
			// 如果权重为0或负数，则永久禁用该apikey
			continue
		}
		totalWeight += weight
	}
	if totalWeight == 0 {
		// 如果所有apikey的权重都是0，则随机返回一个apikey
		return gpt.ApiKeys[rand.Intn(len(gpt.ApiKeys))], nil
	}
	randNum := rand.Intn(totalWeight)
	for i, weight := range gpt.apiKeyWeights {
		if weight <= 0 {
			// 如果权重为0或负数，则跳过该apikey
			continue
		}
		if randNum < weight {
			gpt.currentApiKeyIndex = i
			break
		}
		randNum -= weight
	}
	return gpt.ApiKeys[gpt.currentApiKeyIndex], nil
}

func (gpt *ChatGPT) calculateApiKeyWeights() {
	// 计算apikey的权重
	totalUsage := 0
	for _, usage := range gpt.apiKeyUsage {
		totalUsage += usage
	}
	gpt.apiKeyWeights = make([]int, len(gpt.ApiKeys))
	for i, apiKey := range gpt.ApiKeys {
		usage := gpt.apiKeyUsage[apiKey]
		if usage >= 0 {
			gpt.apiKeyWeights[i] = totalUsage - usage + 1
		} else {
			gpt.apiKeyWeights[i] = 0
		}
	}
}

func (gpt *ChatGPT) initApiKeyUsage() {
	// 初始化apikey的调用次数和权重
	gpt.apiKeyUsage = make(map[string]int)
	for _, apiKey := range gpt.ApiKeys {
		gpt.apiKeyUsage[apiKey] = 0
	}
}

func (gpt *ChatGPT) checkApiKeyAvailability(apiKey string) bool {
	// 检查apikey的可用性
	req, err := http.NewRequest("GET", BASEURL+"engines", nil)
	if err != nil {
		log.Printf("check api key %s failed: %v", apiKey, err)
		return false
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)
	client := &http.Client{Timeout: 10 * time.Second}
	response, err := client.Do(req)
	if err != nil {
		log.Printf("check api key %s failed: %v", apiKey, err)
		return false
	}
	if response.StatusCode/2 != 100 {
		log.Printf("api key %s is not available: %s", apiKey, response.Status)
		if response.StatusCode == 429 {
			gpt.apiKeyUsage[apiKey] = -1
			log.Printf("api key %s has been permanently disabled due to too many requests", apiKey)
		}
		return false
	}
	return true
}

func (gpt *ChatGPT) setApiKeyAvailability() {
	// 检查apikey的可用性
	for apiKey := range gpt.apiKeyUsage {
		if usage := gpt.apiKeyUsage[apiKey]; usage < 0 {
			// 如果apikey的调用次数小于0，则禁用它
			continue
		}
		if !gpt.checkApiKeyAvailability(apiKey) {
			// 如果apikey不可用，则将其调用次数设为-1
			gpt.apiKeyUsage[apiKey] = -1
			gpt.calculateApiKeyWeights()
		}
	}
	fmt.Printf("api key usage: %v", gpt.apiKeyUsage)
}

func (gpt *ChatGPT) StartApiKeyAvailabilityCheck() {
	// 初始化apikey的调用次数和权重，检查apikey的可用性，并定时检查apikey的可用性
	if len(gpt.ApiKeys) == 0 {
		log.Fatalf("no api keys provided")
	}
	gpt.initApiKeyUsage()
	gpt.setApiKeyAvailability()
	ticker := time.NewTicker(30 * time.Minute)
	go func() {
		for {
			select {
			case <-ticker.C:
				gpt.setApiKeyAvailability()
			}
		}
	}()
}

func (gpt *ChatGPT) updateApiKeyUsage(apiKey string) {
	// 更新ApiKeyUsage并写入文件
	gpt.apiKeyUsage[apiKey]++
	gpt.calculateApiKeyWeights()
}

func (gpt *ChatGPT) loadApiKeyUsage() error {
	// 从本地文件中读取apikey的调用次数和权重信息
	file, err := os.OpenFile("apikey_usage.json", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}
	if fileInfo.Size() == 0 {
		gpt.initApiKeyUsage()
		gpt.calculateApiKeyWeights()
		err = gpt.saveApiKeyUsage(file)
		if err != nil {
			log.Printf("save api key usage failed: %v", err)
		}
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

func (gpt *ChatGPT) saveApiKeyUsage(file *os.File) error {
	// 将apikey的调用次数和权重信息输出到本地文件
	file.Truncate(0)
	file.Seek(0, 0)
	encoder := json.NewEncoder(file)
	err := encoder.Encode(gpt.apiKeyUsage)
	if err != nil {
		return err
	}
	return nil
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
	apiKey, err := gpt.getApiKey()
	if err != nil {
		return nil, err
	}
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

func NewChatGPT(apiKeys []string) *ChatGPT {
	return &ChatGPT{
		ApiKeys: apiKeys,
	}
}
