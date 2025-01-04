package recognize

import (
	"fmt"
	"io"
	"os"
	"testing"
)

func TestRecognizer_Recognize(t *testing.T) {
	f, err := os.Open("./data/test.wav")
	if err != nil {
		t.Error("failed open test file")
		return
	}
	data, err := io.ReadAll(f)
	if err != nil {
		t.Error("failed read test file")
		return
	}

	r := New("")
	got, err := r.Recognize(data, SetLanguage(LangRu), SetFormat(LPCM))
	if err != nil {
		t.Errorf("Recognize() error = %v", err)
		return
	}
	fmt.Println(got)
	t.Log(got)
}
