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
	MajorBrand   uint32 `qtff:" "`
	MinorVersion uint32 `qtff:" "`
}

type MovieAtom struct {
	*BasicAtom
}

func (a *MovieAtom) Leaf() bool {
	return false
}

type MovieHeaderAtom struct {
	*BasicAtom
	Version           byte   `qtff:" "`
	Flags             []byte `qtff:"3"`
	CreationTime      uint32 `qtff:" "`
	ModificationTime  uint32 `qtff:" "`
	TimeScale         uint32 `qtff:" "`
	Duration          uint32 `qtff:" "`
	PreferredRate     uint32 `qtff:" "` // FIXME: this should be 32-bit fixed point
	PreferredVolue    uint16 `qtff:" "` // FIXME: this should be 16-bit fixed point
	Reserved          []byte `qtff:"10"`
	PreviewTime       uint32 `qtff:" "`
	PosterTime        uint32 `qtff:" "`
	SelectionTime     uint32 `qtff:" "`
	SelectionDuration uint32 `qtff:" "`
	CurrentTime       uint32 `qtff:" "`
	NextTrackId       uint32 `qtff:" "`
}

type TrackAtom struct {
	*BasicAtom
}

func (a *TrackAtom) Leaf() bool {
	return false
}
