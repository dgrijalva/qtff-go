package qtff

type Atom interface {
	// returns the type of atom
	Type() []byte
	// returns true if the atom is a leaf node
	Leaf() bool
	// length of encoded atom
	Length() uint64
	// Child atoms
	Children() []Atom
}

type BasicAtom struct {
	length     uint64
	typ        []byte
	ChildAtoms []Atom
}

func (a *BasicAtom) Type() []byte {
	return a.typ
}

func (a *BasicAtom) Leaf() bool {
	return true
}

func (a *BasicAtom) Length() uint64 {
	return a.length
}

func (a *BasicAtom) Children() []Atom {
	return a.ChildAtoms
}

type FileTypeAtom struct {
	*BasicAtom
	MajorBrand   uint32 `qtff:"0"`
	MinorVersion uint32 `qtff:"1"`
}
