package qtff

import (
	"fmt"
	"os"
	"testing"
)

func TestParser(t *testing.T) {
	file, err := os.Open("./test.mp4")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	if atoms, err := Parse(file); err == nil {
		printAtoms(atoms)
		t.Errorf("%v", atoms)
	} else {
		printAtoms(atoms)
		t.Errorf("Error parsing: %v", err)
	}
}

func printAtoms(as []Atom) {
	for _, a := range as {
		printAtom(a, 0)
	}
}

func printAtom(a Atom, depth int) {
	padding := ""
	for i := 0; i < depth; i++ {
		padding = padding + "    "
	}
	if _, ok := a.(*BasicAtom); ok {
		fmt.Println(padding, string(a.Type()), a.Length(), "<Basic>")
	} else {
		fmt.Println(padding, string(a.Type()), a.Length(), a)
	}
	children := a.Children()
	if children != nil {
		for _, c := range children {
			printAtom(c, depth+1)
		}
	}
}
