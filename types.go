package qtff

import (
	"encoding/binary"
	"io"
)

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

// atom type knows how to parse the data portion itself
type dataAtom interface {
	parseRemainingData(rdr io.Reader) error
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
	PreferredVolume   uint16 `qtff:" "` // FIXME: this should be 16-bit fixed point
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

type TrackHeaderAtom struct {
	*BasicAtom
	Version          byte   `qtff:" "`
	Flags            []byte `qtff:"3"`
	CreationTime     uint32 `qtff:" "`
	ModificationTime uint32 `qtff:" "`
	TrackId          uint32 `qtff:" "`
	Reserved         []byte `qtff:"4"`
	Duration         uint32 `qtff:" "`
	Reserved2        []byte `qtff:"8"`
	Layer            uint16 `qtff:" "`
	AlternateGroup   uint16 `qtff:" "`
	Volume           uint16 `qtff:" "` // FIXME: this should be 16-bit fixed point
	Reserved3        []byte `qtff:"2"`
	MatrixStructure  []byte `qtff:"36"`
	TrackWidth       uint32 `qtff:" "` // FIXME: this should be 32-bit fixed point
	TrackHeight      uint32 `qtff:" "` // FIXME: this should be 32-bit fixed point
}

type EditAtom struct {
	*BasicAtom
}

func (a *EditAtom) Leaf() bool {
	return false
}

type EditListAtom struct {
	*BasicAtom
	Version  byte   `qtff:" "`
	Flags    []byte `qtff:"3"`
	NumEdits uint32 `qtff:" "`
	Edits    []EditListEdit
}

type EditListEdit struct {
	TrackDuration uint32
	MediaTime     int32
	MediaRate     uint32 // FIXME: this should be 32-bit fixed point
}

func (e *EditListAtom) parseRemainingData(rdr io.Reader) error {
	var readBuffer = make([]byte, 12)
	e.Edits = make([]EditListEdit, e.NumEdits)
	for i := uint32(0); i < e.NumEdits; i++ {
		if _, err := rdr.Read(readBuffer); err == nil {
			e.Edits[i] = EditListEdit{
				binary.BigEndian.Uint32(readBuffer[0:4]),
				int32(binary.BigEndian.Uint32(readBuffer[4:8])),
				binary.BigEndian.Uint32(readBuffer[8:12]),
			}
		} else {
			return err
		}
	}
	return nil
}

type MediaAtom struct {
	*BasicAtom
}

func (a *MediaAtom) Leaf() bool {
	return false
}

type MediaHeaderAtom struct {
	*BasicAtom
	Version          byte   `qtff:" "`
	Flags            []byte `qtff:"3"`
	CreationTime     uint32 `qtff:" "`
	ModificationTime uint32 `qtff:" "`
	TimeScale        uint32 `qtff:" "`
	Duration         uint32 `qtff:" "`
	Language         uint16 `qtff:" "`
	Quality          uint16 `qtff:" "`
}

type MediaInfoAtom struct {
	*BasicAtom
}

func (a *MediaInfoAtom) Leaf() bool {
	return false
}
