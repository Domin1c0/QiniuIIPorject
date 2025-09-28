package model

type RequestMessage struct {
	Content string `json:"content"`
}

type RequestNewChat struct {
	// TODO more fields e.x. character, session name
	Message RequestMessage `json:"message"`
}

type RequestSend struct {
	// Send message to existing session needs only content
	Message RequestMessage `json:"message"`
}
