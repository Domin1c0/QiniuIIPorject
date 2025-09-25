package api

import storage "github.com/LTSlw/QiniuIIPorject/backend/pkg/storge"

type Model struct {
	Addr      string `json:"addr"`
	ModelName string `json:"model_name"`
	ApiKey    string `json:"api_key"`
}

func CallLLM(model Model, messages []storage.Message) (string, error) {
	// TODO
	return "", nil
}
