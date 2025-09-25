package llm

import "testing"

func TestGetStringToken(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		model    Model
		minCount int
		maxCount int
	}{
		{
			name:     "gpt-4o simple string",
			input:    "hello world",
			model:    Model{ModelName: "gpt-4o"},
			minCount: 2,   // token range
			maxCount: 500, // Accuracy is not a part of testing
		},
		{
			name:     "gpt-3.5 turbo chinese",
			input:    "我是月社妃的狗，月社妃你带我走吧，我想和月社妃一起被大运创死",
			model:    Model{ModelName: "gpt-3.5-turbo"},
			minCount: 2,
			maxCount: 500,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetStringToken(tt.input, tt.model)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got < tt.minCount || got > tt.maxCount {
				t.Errorf("GetStringToken() = %v, want between %v and %v", got, tt.minCount, tt.maxCount)
			} else {
				t.Logf("GetStringToken(%q, %s) = %d tokens", tt.input, tt.model.ModelName, got)
			}
		})
	}
}
