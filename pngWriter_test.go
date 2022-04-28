package weather2png_server

import (
	"testing"
	"os"
)

func TestPngWriter(t *testing.T) {
	f, err := os.Create("output.png")
	if err != nil {
		t.Error(err)
		return
	}
	defer f.Close()
	wr := NewPngWriter(600, 800, os.Getenv("TEST_PNGWRITE_TTF_PATH"))
	wr.Text("Hello,你好啊！", 80, 0, 40.0)
	wr.Reset(f)
}
