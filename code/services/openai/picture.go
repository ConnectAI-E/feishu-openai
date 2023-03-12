package openai

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
)

type ImageGenerationRequestBody struct {
	Prompt         string `json:"prompt"`
	N              int    `json:"n"`
	Size           string `json:"size"`
	ResponseFormat string `json:"response_format"`
}

type ImageResponseBody struct {
	Created int64 `json:"created"`
	Data    []struct {
		Base64Json string `json:"b64_json"`
	} `json:"data"`
}

type ImageVariantRequestBody struct {
	Image          string `json:"image"`
	N              int    `json:"n"`
	Size           string `json:"size"`
	ResponseFormat string `json:"response_format"`
}

func (gpt ChatGPT) GenerateImage(prompt string, size string, n int) ([]string, error) {
	requestBody := ImageGenerationRequestBody{
		Prompt:         prompt,
		N:              n,
		Size:           size,
		ResponseFormat: "b64_json",
	}

	imageResponseBody := &ImageResponseBody{}
	err := gpt.sendRequestWithBodyType(gpt.ApiUrl+"/v1/images/generations",
		"POST", jsonBody, requestBody, imageResponseBody)

	if err != nil {
		return nil, err
	}

	var b64Pool []string
	for _, data := range imageResponseBody.Data {
		b64Pool = append(b64Pool, data.Base64Json)
	}
	return b64Pool, nil
}

func (gpt ChatGPT) GenerateOneImage(prompt string, size string) (string, error) {
	b64s, err := gpt.GenerateImage(prompt, size, 1)
	if err != nil {
		return "", err
	}
	return b64s[0], nil
}

func (gpt ChatGPT) GenerateOneImageWithDefaultSize(prompt string) (string, error) {
	return gpt.GenerateOneImage(prompt, "512x512")
}

func (gpt ChatGPT) GenerateImageVariation(images string, size string, n int) ([]string, error) {
	requestBody := ImageVariantRequestBody{
		Image:          images,
		N:              n,
		Size:           size,
		ResponseFormat: "b64_json",
	}

	imageResponseBody := &ImageResponseBody{}
	err := gpt.sendRequestWithBodyType(gpt.ApiUrl+"/v1/images/variations",
		"POST", formPictureDataBody, requestBody, imageResponseBody)

	if err != nil {
		return nil, err
	}

	var b64Pool []string
	for _, data := range imageResponseBody.Data {
		b64Pool = append(b64Pool, data.Base64Json)
	}
	return b64Pool, nil
}

func (gpt ChatGPT) GenerateOneImageVariation(images string, size string) (string, error) {
	b64s, err := gpt.GenerateImageVariation(images, size, 1)
	if err != nil {
		return "", err
	}
	return b64s[0], nil
}

func pictureMultipartForm(request ImageVariantRequestBody,
	w *multipart.Writer) error {

	f, err := os.Open(request.Image)
	if err != nil {
		return fmt.Errorf("opening audio file: %w", err)
	}
	fw, err := w.CreateFormFile("image", f.Name())
	if err != nil {
		return fmt.Errorf("creating form file: %w", err)
	}
	if _, err = io.Copy(fw, f); err != nil {
		return fmt.Errorf("reading from opened audio file: %w", err)
	}

	err = w.WriteField("size", request.Size)
	if err != nil {
		return fmt.Errorf("writing size: %w", err)
	}

	err = w.WriteField("n", fmt.Sprintf("%d", request.N))
	if err != nil {
		return fmt.Errorf("writing n: %w", err)
	}

	err = w.WriteField("response_format", request.ResponseFormat)
	if err != nil {
		return fmt.Errorf("writing response_format: %w", err)
	}

	//fw, err = w.CreateFormField("model")
	//if err != nil {
	//	return fmt.Errorf("creating form field: %w", err)
	//}
	//modelName := bytes.NewReader([]byte(request.Model))
	//if _, err = io.Copy(fw, modelName); err != nil {
	//	return fmt.Errorf("writing model name: %w", err)
	//}

	//fmt.Printf("w.FormDataContentType(): %s ", w.FormDataContentType())

	w.Close()

	return nil
}
