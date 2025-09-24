package tts

import (
	"errors"

	"github.com/LTSlw/QiniuIIPorject/backend/pkg/gradio"
)

const (
	urlTTS = "https://yasserrmd-vibevoice.hf.space/gradio_api"
)

func TTS(text string, character string) ([]byte, error) {
	client := gradio.NewClient(urlTTS)
	client.AppendString("VibeVoice-Large")
	client.AppendInt(1)
	client.AppendString(text)
	client.AppendString(character)
	client.AppendString(character)
	client.AppendString(character)
	client.AppendString(character)
	client.AppendFloat(1.3)

	_, err := client.Predict("generate_podcast_wrapper")
	if err != nil {
		return nil, err
	}
	result, err := client.Result()
	if err != nil {
		return nil, err
	}

	path, ok := result[0].(string)
	if !ok {
		return nil, errors.New("rate limited")
	}

	body, err := client.DownloadFile(path)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func GetCharacters() []string {
	return []string{
		"ar-Ahmed-man",
		"ar-Hamad-man",
		"ar-Hamdan-man",
		"ar-Yasser-man",
		"en-Alice_woman",
		"en-Alice_woman_bgm",
		"en-Carter_man",
		"en-Frank_man",
		"en-Maya_woman",
		"en-Yasser_man",
		"hi-Asad-man",
		"in-Samuel_man",
		"ta-Hameed-man",
		"ta-Thenkatchi-man",
		"ta-Yasser-man",
		"zh-Anchen_man_bgm",
		"zh-Bowen_man",
		"zh-Xinran_woman",
	}
}
