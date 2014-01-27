package qtff

type Atom interface {
}

type FileTypeAtom struct {
	Length uint64
	Type []byte
}