package openai

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"start-feishubot/initialization"
	"start-feishubot/logger"
	"start-feishubot/services/loadbalancer"
	"strings"
	"time"
)

type PlatForm string

const (
	AzureApiUrlV1 = "openai.azure.com/openai/deployments/"
)
const (
	OpenAI PlatForm = "openai"
	Azure  PlatForm = "azure"
)

type AzureConfig struct {
	BaseURL        string
	ResourceName   string
	DeploymentName string
	ApiVersion     string
	ApiToken       string
}

type ChatGPT struct {
	Lb          *loadbalancer.LoadBalancer
	ApiKey      []string
	ApiUrl      string
	HttpProxy   string
	Model       string
	MaxTokens   int
	Platform    PlatForm
	AzureConfig AzureConfig
}
type requestBodyType int

const (
	jsonBody requestBodyType = iota
	formVoiceDataBody
	formPictureDataBody

	nilBody
)

func (gpt *ChatGPT) doAPIRequestWithRetry(url, method string,
	bodyType requestBodyType,
	requestBody interface{}, responseBody interface{}, client *http.Client, maxRetries int) error {
	var api *loadbalancer.API
	var requestBodyData []byte
	var err error
	var writer *multipart.Writer
	api = gpt.Lb.GetAPI()

	switch bodyType {
	case jsonBody:
		requestBodyData, err = json.Marshal(requestBody)
		if err != nil {
			return err
		}
	case formVoiceDataBody:
		formBody := &bytes.Buffer{}
		writer = multipart.NewWriter(formBody)
		err = audioMultipartForm(requestBody.(AudioToTextRequestBody), writer)
		if err != nil {
			return err
		}
		err = writer.Close()
		if err != nil {
			return err
		}
		requestBodyData = formBody.Bytes()
	case formPictureDataBody:
		formBody := &bytes.Buffer{}
		writer = multipart.NewWriter(formBody)
		err = pictureMultipartForm(requestBody.(ImageVariantRequestBody), writer)
		if err != nil {
			return err
		}
		err = writer.Close()
		if err != nil {
			return err
		}
		requestBodyData = formBody.Bytes()
	case nilBody:
		requestBodyData = nil

	default:
		return errors.New("unknown request body type")
	}

	if api == nil {
		return errors.New("no available API")
	}

	req, err := http.NewRequest(method, url, bytes.NewReader(requestBodyData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	if bodyType == formVoiceDataBody || bodyType == formPictureDataBody {
		req.Header.Set("Content-Type", writer.FormDataContentType())
	}
	if gpt.Platform == OpenAI {
		req.Header.Set("Authorization", "Bearer "+api.Key)
	} else {
		req.Header.Set("api-key", gpt.AzureConfig.ApiToken)
	}

	var response *http.Response
	var retry int
	for retry = 0; retry <= maxRetries; retry++ {
		// set body
		if retry > 0 {
			req.Body = ioutil.NopCloser(bytes.NewReader(requestBodyData))
		}
		response, err = client.Do(req)
		//fmt.Println("--------------------")
		//fmt.Println("req", req.Header)
		//fmt.Printf("response: %v", response)
		logger.Debug("req", req.Header)

		logger.Debugf("response %v", response)
		// read body
		if err != nil || response.StatusCode < 200 || response.StatusCode >= 300 {

			body, _ := ioutil.ReadAll(response.Body)
			fmt.Println("body", string(body))

			gpt.Lb.SetAvailability(api.Key, false)
			if retry == maxRetries {
				break
			}
			time.Sleep(time.Duration(retry+1) * time.Second)
		} else {
			break
		}
	}
	if response != nil {
		defer response.Body.Close()
	}

	if response == nil || response.StatusCode < 200 || response.StatusCode >= 300 {
		return fmt.Errorf("%s api failed after %d retries", strings.ToUpper(method), retry)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, responseBody)
	if err != nil {
		return err
	}

	gpt.Lb.SetAvailability(api.Key, true)
	return nil
}

func (gpt *ChatGPT) sendRequestWithBodyType(link, method string,
	bodyType requestBodyType,
	requestBody interface{}, responseBody interface{}) error {
	var err error
	proxyString := gpt.HttpProxy

	client, parseProxyError := GetProxyClient(proxyString)
	if parseProxyError != nil {
		return parseProxyError
	}

	err = gpt.doAPIRequestWithRetry(link, method, bodyType,
		requestBody, responseBody, client, 3)

	return err
}

func NewChatGPT(config initialization.Config) *ChatGPT {
	var lb *loadbalancer.LoadBalancer
	if config.AzureOn {
		keys := []string{config.AzureOpenaiToken}
		lb = loadbalancer.NewLoadBalancer(keys)
	} else {
		lb = loadbalancer.NewLoadBalancer(config.OpenaiApiKeys)
	}
	platform := OpenAI

	if config.AzureOn {
		platform = Azure
	}

	return &ChatGPT{
		Lb:        lb,
		ApiKey:    config.OpenaiApiKeys,
		ApiUrl:    config.OpenaiApiUrl,
		HttpProxy: config.HttpProxy,
		Model:     config.OpenaiModel,
		MaxTokens: config.OpenaiMaxTokens,
		Platform:  platform,
		AzureConfig: AzureConfig{
			BaseURL:        AzureApiUrlV1,
			ResourceName:   config.AzureResourceName,
			DeploymentName: config.AzureDeploymentName,
			ApiVersion:     config.AzureApiVersion,
			ApiToken:       config.AzureOpenaiToken,
		},
	}
}

func (gpt *ChatGPT) FullUrl(suffix string) string {
	var url string
	switch gpt.Platform {
	case Azure:
		url = fmt.Sprintf("https://%s.%s%s/%s?api-version=%s",
			gpt.AzureConfig.ResourceName, gpt.AzureConfig.BaseURL,
			gpt.AzureConfig.DeploymentName, suffix, gpt.AzureConfig.ApiVersion)
	case OpenAI:
		url = fmt.Sprintf("%s/v1/%s", gpt.ApiUrl, suffix)
	}
	return url
}

func GetProxyClient(proxyString string) (*http.Client, error) {
	var client *http.Client
	timeOutDuration := time.Duration(initialization.GetConfig().OpenAIHttpClientTimeOut) * time.Second
	if proxyString == "" {
		client = &http.Client{Timeout: timeOutDuration}
	} else {
		proxyUrl, err := url.Parse(proxyString)
		if err != nil {
			return nil, err
		}
		transport := &http.Transport{
			Proxy: http.ProxyURL(proxyUrl),
		}
		client = &http.Client{
			Transport: transport,
			Timeout:   timeOutDuration,
		}
	}
	return client, nil
}
