package stt

import (
	"fmt"
	"os"
	"testing"
)

func TestSTT(t *testing.T) {
	audio, err := os.ReadFile("../samples/wikipedia_ai.mp3")
	if err != nil {
		t.Fatal(err)
	}
	result, err := STT(audio, os.Getenv("HF_TOKEN"))
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(result)
}
