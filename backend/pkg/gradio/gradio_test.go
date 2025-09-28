package gradio

import (
	"fmt"
	"os"
	"testing"
)

func TestClient_Predict(t *testing.T) {
	file, err := os.ReadFile("../samples/wikipedia_ai.mp3")
	if err != nil {
		t.Fatal(err)
	}

	client := NewClient("openai/whisper")
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

func TestClient_GetJWT(t *testing.T) {
	client := NewClient("openai/whisper")
	client.SetHFToken(os.Getenv("HF_TOKEN"))
	jwt, err := client.getJWT()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("jwt: ", jwt)
}
