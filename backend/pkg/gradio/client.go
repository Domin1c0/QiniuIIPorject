package gradio

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/tmaxmax/go-sse"
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
	spaceUrl string
	jwtUrl   string
	data     []any
	apiName  string
	eventID  string
	hfToken  string
}

func NewClient(space string) *Client {
	parts := strings.Split(space, "/")
	if len(parts) != 2 {
		return nil
	}
	username, spacename := parts[0], parts[1]
	spaceUrl := fmt.Sprintf("https://%s-%s.hf.space/gradio_api", username, spacename)
	jwtUrl := fmt.Sprintf("https://huggingface.co/api/spaces/%s/%s/jwt", username, spacename)

	return &Client{
		spaceUrl: spaceUrl,
		jwtUrl:   jwtUrl,
		data:     make([]any, 0),
	}
}

func (x *Client) SetHFToken(hfToken string) {
	x.hfToken = hfToken
}

func (x *Client) Predict(apiName string) (string, error) {
	x.apiName = apiName
	reqBody := &bodyPredict{
		Data: x.data,
	}
	bodyJson, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", x.spaceUrl+urlPredict+apiName, bytes.NewBuffer(bodyJson))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	if x.hfToken != "" {
		req.Header.Set("Authorization", "Bearer "+x.hfToken)
		jwt, err := x.getJWT()
		if err != nil {
			return "", err
		}
		req.Header.Set("x-zerogpu-token", jwt)
	}

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

	x.eventID = bodyPredict.EventID
	return bodyPredict.EventID, nil
}

func (x *Client) getJWT() (string, error) {
	req, err := http.NewRequest("GET", x.jwtUrl, nil)
	if err != nil {
		return "", err
	}
	if x.hfToken != "" {
		req.Header.Set("Authorization", "Bearer "+x.hfToken)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	jwt := &respJWT{}
	err = json.Unmarshal(body, &jwt)
	if err != nil {
		return "", err
	}
	return jwt.Token, nil
}

func (x *Client) Result() ([]any, error) {
	req, err := http.NewRequest("GET", x.spaceUrl+urlPredict+x.apiName+"/"+x.eventID, nil)
	if err != nil {
		return nil, err
	}
	if x.hfToken != "" {
		req.Header.Set("Authorization", "Bearer "+x.hfToken)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	for ev, err := range sse.Read(resp.Body, nil) {
		if err != nil {
			return nil, err
		}
		switch ev.Type {
		case "complete":
			res := make([]any, 0)
			if err := json.Unmarshal([]byte(ev.Data), &res); err != nil {
				return nil, err
			}
			return res, nil
		case "error":
			return nil, errors.New(ev.Data)
		case "generating":
		case "heartbeat":
		}
	}
	return nil, errors.New("should not happen")
}

func (x *Client) AppendString(s string) {
	x.data = append(x.data, s)
}

func (x *Client) AppendBool(b bool) {
	x.data = append(x.data, b)
}

func (x *Client) AppendInt(n int) {
	x.data = append(x.data, n)
}

func (x *Client) AppendFloat(n float64) {
	x.data = append(x.data, n)
}

func (x *Client) UploadFile(filename string, file []byte, mimeType string) error {
	path, err := x.uploadFile(filename, file)
	if err != nil {
		return err
	}
	x.data = append(x.data, newDataFile(mimeType, path, filename, int64(len(file)), x.spaceUrl+urlFile+"="+path))
	return nil
}

func (x *Client) UploadFileFromURL(s string) error {
	x.data = append(x.data, &dataFile{
		Path: s,
	})
	return nil
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

func (x *Client) DownloadFile(path string) ([]byte, error) {
	req, err := http.NewRequest("GET", x.spaceUrl+urlFile+"="+path, nil)
	if err != nil {
		return nil, err
	}
	if x.hfToken != "" {
		req.Header.Set("Authorization", "Bearer "+x.hfToken)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
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

	req, err := http.NewRequest("POST", x.spaceUrl+urlUpload, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	if x.hfToken != "" {
		req.Header.Set("Authorization", "Bearer "+x.hfToken)
	}

	// set params
	params := req.URL.Query()
	params.Set("upload_id", id)
	req.URL.RawQuery = params.Encode()

	return req, nil
}
