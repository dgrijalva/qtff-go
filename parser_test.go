package qtff

import (
	"os"
	"testing"
	// "fmt"
)

func TestParser(t *testing.T) {
	file, err := os.Open("./test.mp4")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	if atoms, err := Parse(file); err == nil {
		t.Errorf("%v", atoms)
	} else {
		t.Errorf("Error parsing: %v", err)
	}
}
