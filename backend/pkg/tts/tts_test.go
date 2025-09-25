package tts

import (
	"fmt"
	"net/http"
	"os"
	"testing"
)

const (
	text = "Artificial intelligence is the intelligence of machines or software, as opposed to the intelligence of humans or animals. It is also the field of study in computer science that develops and studies intelligent machines."
)

func TestGetCharacters(t *testing.T) {
	characters := GetCharacters()
	fmt.Println("characters: ", characters)
}

func TestTTS(t *testing.T) {
	audio, err := TTS(text, "zh-Xinran_woman", os.Getenv("HF_TOKEN"))
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("audio format: ", http.DetectContentType(audio))
}
