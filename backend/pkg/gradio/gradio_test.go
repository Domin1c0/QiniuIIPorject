package gradio

import (
	"fmt"
	"os"
	"testing"
)

func TestClient_Predict(t *testing.T) {
	file, err := os.ReadFile("wikipedia_ai.mp3")
	if err != nil {
		t.Fatal(err)
	}

	client := NewClient("https://openai-whisper.hf.space/gradio_api")
	client.UploadFile("wikipedia_ai.mp3", file, "audio/mpeg")
	client.AppendString("transcribe")
	eventID, err := client.Predict("predict")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("eventID: ", eventID)

	result, err := client.Result()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("result: ", result)
}
