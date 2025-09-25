package api

import (
	"context"
	"errors"
	"time"

	storage "github.com/LTSlw/QiniuIIPorject/backend/pkg/storge"
	"github.com/openai/openai-go/v2"
	"github.com/openai/openai-go/v2/option"
)

var (
	errInvalidRole = errors.New("invalid role in message")
	errNoChoices   = errors.New("no choices returned from llm")
)

type Model struct {
	Addr      string `json:"addr"`
	ModelName string `json:"model_name"`
	ApiKey    string `json:"api_key"`
}

func CallLLM(model Model, messages []storage.Message) (storage.Message, error) {
	var res storage.Message
	client := openai.NewClient(
		option.WithBaseURL(model.Addr),
		option.WithAPIKey(model.ApiKey),
	)

	// build messages
	// this function builds all the messages to request from input
	// assume that messages are selected from outside the function?
	msgs := make([]openai.ChatCompletionMessageParamUnion, 0, len(messages))
	for _, message := range messages {
		switch message.Role {
		case "user":
			msgs = append(msgs, openai.UserMessage(message.Content))
		case "system":
			msgs = append(msgs, openai.SystemMessage(message.Content))
		case "assistant":
			msgs = append(msgs, openai.AssistantMessage(message.Content))
		default:
			return res, errInvalidRole
		}
	}

	params := openai.ChatCompletionNewParams{
		// Model:    openai.(model.ModelName),
		Model:    model.ModelName,
		Messages: msgs,
	}

	chatCompletion, err := client.Chat.Completions.New(context.TODO(), params)
	if err != nil {
		return res, err
	}
	if len(chatCompletion.Choices) == 0 {
		return res, errNoChoices
	}

	res.Role = "assistant"
	res.Content = chatCompletion.Choices[0].Message.Content
	res.CreateAt = time.Unix(chatCompletion.Created, 0)
	return res, nil

	// return chatCompletion.Choices[0].Message.Content, nil

}
