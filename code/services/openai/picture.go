package openai

import (
	"bufio"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"mime/multipart"
	"os"
)

type ImageGenerationRequestBody struct {
	Prompt         string `json:"prompt"`
	N              int    `json:"n"`
	Size           string `json:"size"`
	ResponseFormat string `json:"response_format"`
	Model          string `json:"model,omitempty"`
	Style          string `json:"style,omitempty"`
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

func (gpt *ChatGPT) GenerateImage(prompt string, size string,
	n int, style string) ([]string, error) {
	requestBody := ImageGenerationRequestBody{
		Prompt:         prompt,
		N:              n,
		Size:           size,
		ResponseFormat: "b64_json",
		Model:          "dall-e-3",
		Style:          style,
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

func (gpt *ChatGPT) GenerateOneImage(prompt string,
	size string, style string) (string, error) {
	b64s, err := gpt.GenerateImage(prompt, size, 1, style)
	if err != nil {
		return "", err
	}
	return b64s[0], nil
}

func (gpt *ChatGPT) GenerateOneImageWithDefaultSize(
	prompt string) (string, error) {
	// works for dall-e 2&3
	return gpt.GenerateOneImage(prompt, "1024x1024", "")
}

func (gpt *ChatGPT) GenerateImageVariation(images string,
	size string, n int) ([]string, error) {
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

func (gpt *ChatGPT) GenerateOneImageVariation(images string,
	size string) (string, error) {
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

	//err = w.WriteField("user", "user123456")

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

func VerifyPngs(pngPaths []string) error {
	foundPng := false
	var expectedWidth, expectedHeight int

	for _, pngPath := range pngPaths {
		f, err := os.Open(pngPath)
		if err != nil {
			return fmt.Errorf("os.Open: %v", err)
		}

		fi, err := f.Stat()
		if err != nil {
			return fmt.Errorf("f.Stat: %v", err)
		}
		if fi.Size() > 4*1024*1024 {
			return fmt.Errorf("image size too large, "+
				"must be under %d MB", 4)
		}

		image, err := png.Decode(f)
		if err != nil {
			return fmt.Errorf("image must be valid png, got error: %v", err)
		}
		width := image.Bounds().Dx()
		height := image.Bounds().Dy()
		if width != height {
			return fmt.Errorf("found non-square image with dimensions %dx%d", width, height)
		}

		if !foundPng {
			foundPng = true
			expectedWidth = width
			expectedHeight = height
		} else {
			if width != expectedWidth || height != expectedHeight {
				return fmt.Errorf("dimensions of all images must match, got both (%dx%d) and (%dx%d)", width, height, expectedWidth, expectedHeight)
			}
		}
	}

	return nil
}

func ConvertToRGBA(inputFilePath string, outputFilePath string) error {
	// 打开输入文件
	inputFile, err := os.Open(inputFilePath)
	if err != nil {
		return fmt.Errorf("打开文件时出错：%w", err)
	}
	defer inputFile.Close()

	// 解码图像
	img, _, err := image.Decode(inputFile)
	if err != nil {
		return fmt.Errorf("解码图像时出错：%w", err)
	}

	// 将图像转换为RGBA模式
	rgba := image.NewRGBA(img.Bounds())
	for x := 0; x < img.Bounds().Max.X; x++ {
		for y := 0; y < img.Bounds().Max.Y; y++ {
			rgba.Set(x, y, img.At(x, y))
		}
	}

	// 创建输出文件
	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		return fmt.Errorf("创建输出文件时出错：%w", err)
	}
	defer outputFile.Close()

	// 编码图像为 PNG 格式并写入输出文件
	if err := png.Encode(outputFile, rgba); err != nil {
		return fmt.Errorf("编码图像时出错：%w", err)
	}

	return nil
}

func ConvertJpegToPNG(jpgPath string) error {
	// Open the JPEG file for reading
	f, err := os.Open(jpgPath)
	if err != nil {
		return err
	}
	defer f.Close()

	// Check if the file is a JPEG image
	_, err = jpeg.Decode(f)
	if err != nil {
		// The file is not a JPEG image, no need to convert it
		return fmt.Errorf("file %s is not a JPEG image", jpgPath)
	}

	// Reset the file pointer to the beginning of the file
	_, err = f.Seek(0, 0)
	if err != nil {
		return err
	}

	// Create a new PNG file for writing
	pngPath := jpgPath[:len(jpgPath)-4] + ".png" // replace .jpg extension with .png
	out, err := os.Create(pngPath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Decode the JPEG image and encode it as PNG
	img, err := jpeg.Decode(f)
	if err != nil {
		return err
	}
	err = png.Encode(out, img)
	if err != nil {
		return err
	}

	return nil
}

func GetImageCompressionType(path string) (string, error) {
	// 打开文件
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// 创建 bufio.Reader
	reader := bufio.NewReader(file)

	// 解码图像
	_, format, err := image.DecodeConfig(reader)
	if err != nil {
		fmt.Println("err: ", err)
		return "", err
	}

	fmt.Println("format: ", format)
	// 返回压缩类型
	return format, nil
}
