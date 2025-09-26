package llm

import storage "github.com/LTSlw/QiniuIIPorject/backend/pkg/storge"

type Model struct {
	Addr      string `json:"addr"`
	ModelName string `json:"model_name"`
	ApiKey    string `json:"api_key"`
}

type SessionsWithMessages struct {
	Session  storage.Session   `json:"session"`
	Messages []storage.Message `json:"messages"`
}
