package llm

type Model struct {
	Addr      string `json:"addr"`
	ModelName string `json:"model_name"`
	ApiKey    string `json:"api_key"`
}
