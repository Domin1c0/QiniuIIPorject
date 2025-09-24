package stt

import "github.com/LTSlw/QiniuIIPorject/backend/pkg/gradio"

func STT(audio []byte) (string, error) {
	client := gradio.NewClient("https://openai-whisper.hf.space/gradio_api")
	client.UploadFile("audio.mp3", audio, "audio/mpeg")
	client.AppendString("transcribe")
	_, err := client.Predict("predict")
	if err != nil {
		return "", err
	}
	result, err := client.Result()
	if err != nil {
		return "", err
	}
	return result[0].(string), nil
}
