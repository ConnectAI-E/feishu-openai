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
	"start-feishubot/services/loadbalancer"
	"strings"
	"time"
)

type ChatGPT struct {
	Lb        *loadbalancer.LoadBalancer
	ApiKey    []string
	ApiUrl    string
	HttpProxy string
}
type requestBodyType int

const (
	jsonBody requestBodyType = iota
	formVoiceDataBody
	formPictureDataBody

	nilBody
)

func (gpt ChatGPT) doAPIRequestWithRetry(url, method string, bodyType requestBodyType,
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
	req.Header.Set("Authorization", "Bearer "+api.Key)

	var response *http.Response
	var retry int
	for retry = 0; retry <= maxRetries; retry++ {
		response, err = client.Do(req)
		//fmt.Println("--------------------")
		//fmt.Println("req", req.Header)
		//fmt.Printf("response: %v", response)
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

func (gpt ChatGPT) sendRequestWithBodyType(link, method string, bodyType requestBodyType,
	requestBody interface{}, responseBody interface{}) error {
	var err error
	client := &http.Client{Timeout: 110 * time.Second}
	if gpt.HttpProxy == "" {
		err = gpt.doAPIRequestWithRetry(link, method, bodyType,
			requestBody, responseBody, client, 3)
	} else {
		proxyUrl, err := url.Parse(gpt.HttpProxy)
		if err != nil {
			return err
		}
		transport := &http.Transport{
			Proxy: http.ProxyURL(proxyUrl),
		}
		proxyClient := &http.Client{
			Transport: transport,
			Timeout:   110 * time.Second,
		}
		err = gpt.doAPIRequestWithRetry(link, method, bodyType,
			requestBody, responseBody, proxyClient, 3)
	}

	return err
}

func NewChatGPT(config initialization.Config) *ChatGPT {
	apiKeys := config.OpenaiApiKeys
	apiUrl := config.OpenaiApiUrl
	httpProxy := config.HttpProxy
	lb := loadbalancer.NewLoadBalancer(apiKeys)
	return &ChatGPT{
		Lb:        lb,
		ApiKey:    apiKeys,
		ApiUrl:    apiUrl,
		HttpProxy: httpProxy,
	}
}
