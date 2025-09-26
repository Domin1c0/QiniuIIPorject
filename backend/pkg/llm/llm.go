package llm

import (
	storage "github.com/LTSlw/QiniuIIPorject/backend/pkg/storge"
	"github.com/pkoukk/tiktoken-go"
)

// Select latest messages from a session based on the token limit.
// the input session should contain the latest user input, in order to make the calculated tokens contains the latest user input.
func SelectMessage(session SessionsWithMessages, tokenLimit int, model Model) ([]storage.Message, error) {
	var selectedMessages []storage.Message
	// return selectedMessages, nil

	sumToken := 0
	// for message := range session.Messages
	for i := len(session.Messages) - 1; i >= 0; i-- {
		msg := session.Messages[i]
		token, err := GetStringToken(msg.Content, model)
		if err != nil {
			return nil, err
		}
		if sumToken+token >= tokenLimit {
			// if the sum of tokens exceeds the limit, break
			break
		} else {
			selectedMessages = append([]storage.Message{msg}, selectedMessages...) // revert append
			sumToken += token
		}
	}
	return selectedMessages, nil
}

func GetStringToken(input string, model Model) (int, error) {
	// enc, err := tiktoken.EncodingForModel(model.ModelName)
	var enc *tiktoken.Tiktoken
	var err error
	switch model.ModelName {
	// o200k_base model
	case "gpt-4o", "gpt-4o-mini", "gpt-4.1", "gpt-4.5":
		enc, err = tiktoken.GetEncoding("o200k_base")

	// cl100k_base model
	case "gpt-4", "gpt-4-32k", "gpt-3.5-turbo", "gpt-3.5-turbo-16k",
		"text-embedding-ada-002", "text-embedding-3-small", "text-embedding-3-large":
		enc, err = tiktoken.GetEncoding("cl100k_base")

	// p50k_base series（code / davinci-002/003）
	case "text-davinci-002", "text-davinci-003",
		"code-davinci-002", "davinci-codex":
		enc, err = tiktoken.GetEncoding("p50k_base")

	// r50k_base / gpt2 series（GPT-3）
	case "davinci", "curie", "babbage", "ada":
		enc, err = tiktoken.GetEncoding("r50k_base")

	default:
		// default cl100k_base
		enc, err = tiktoken.GetEncoding("cl100k_base")
	}
	if err != nil {
		return 0, err
	}

	tokens := enc.Encode(input, nil, nil)
	if err != nil {
		return 0, err
	}
	return len(tokens), nil
}
