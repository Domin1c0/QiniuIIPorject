package gradio

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"math/rand"
	"mime/multipart"
	"net/http"
	"strings"
)

const (
	urlUpload  = "/upload"
	urlFile    = "/file"
	urlPredict = "/call/"
)

func generateUploadId() string {
	charset := "0123456789abcdefghijklmnopqrstuvwxyz"

	result := make([]byte, 11)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}

	return string(result)
}

type Client struct {
	BaseUrl string
	data    []any
}

func NewClient(baseUrl string) *Client {
	return &Client{
		BaseUrl: baseUrl,
		data:    make([]any, 0),
	}
}

func (x *Client) Predict(apiName string) (string, error) {
	reqBody := &bodyPredict{
		Data: x.data,
	}
	bodyJson, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", x.BaseUrl+urlPredict+apiName, bytes.NewBuffer(bodyJson))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	bodyPredict := &respPredict{}
	if err != nil {
		return "", err
	}

	err = json.Unmarshal(body, &bodyPredict)
	if err != nil {
		return "", err
	}

	return bodyPredict.EventID, nil
}

func (x *Client) AppendString(s string) {
	x.data = append(x.data, s)
}

func (x *Client) UploadFile(filename string, file []byte, mimeType string) error {
	path, err := x.uploadFile(filename, file)
	if err != nil {
		return err
	}
	x.data = append(x.data, newDataFile(mimeType, path, filename, int64(len(file)), x.BaseUrl+urlFile+"="+path))
	return nil
}

func (x *Client) UploadFileFromURL(s string) bool {
	return strings.HasPrefix(s, "http://") || strings.HasPrefix(s, "https://")
}

func (x *Client) uploadFile(filename string, file []byte) (string, error) {
	req, err := x.createUploadRequest(filename, file)
	if err != nil {
		return "", err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	paths := make([]string, 0)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	err = json.Unmarshal(body, &paths)
	if err != nil {
		return "", err
	}
	if len(paths) == 0 {
		return "", errors.New("no paths found")
	}

	return paths[0], nil
}

func (x *Client) createUploadRequest(filename string, file []byte) (*http.Request, error) {
	id := generateUploadId()

	// create multipart/form-data
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	field, err := writer.CreateFormFile("files", filename)
	if err != nil {
		return nil, err
	}
	_, err = field.Write(file)
	if err != nil {
		return nil, err
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", x.BaseUrl+urlUpload, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// set params
	params := req.URL.Query()
	params.Set("upload_id", id)
	req.URL.RawQuery = params.Encode()

	return req, nil
}
