package qtff

import (
	"bytes"
	"encoding/binary"
	"io"
	"io/ioutil"
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
	Version           byte    `qtff:" "`
	Flags             []byte  `qtff:"3"`
	CreationTime      uint32  `qtff:" "`
	ModificationTime  uint32  `qtff:" "`
	TimeScale         uint32  `qtff:" "`
	Duration          uint32  `qtff:" "`
	PreferredRate     float64 `qtff:"4"`
	PreferredVolume   float64 `qtff:"2"`
	Reserved          []byte  `qtff:"10"`
	PreviewTime       uint32  `qtff:" "`
	PosterTime        uint32  `qtff:" "`
	SelectionTime     uint32  `qtff:" "`
	SelectionDuration uint32  `qtff:" "`
	CurrentTime       uint32  `qtff:" "`
	NextTrackId       uint32  `qtff:" "`
}

type TrackAtom struct {
	*BasicAtom
}

func (a *TrackAtom) Leaf() bool {
	return false
}

type TrackHeaderAtom struct {
	*BasicAtom
	Version          byte    `qtff:" "`
	Flags            []byte  `qtff:"3"`
	CreationTime     uint32  `qtff:" "`
	ModificationTime uint32  `qtff:" "`
	TrackId          uint32  `qtff:" "`
	Reserved         []byte  `qtff:"4"`
	Duration         uint32  `qtff:" "`
	Reserved2        []byte  `qtff:"8"`
	Layer            uint16  `qtff:" "`
	AlternateGroup   uint16  `qtff:" "`
	Volume           float64 `qtff:"2"`
	Reserved3        []byte  `qtff:"2"`
	MatrixStructure  []byte  `qtff:"36"`
	TrackWidth       float64 `qtff:"4"`
	TrackHeight      float64 `qtff:"4"`
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
	MediaRate     float64 // 32-bit fixed point
}

func (e *EditListAtom) parseRemainingData(rdr io.Reader) error {
	var readBuffer = make([]byte, 12)
	e.Edits = make([]EditListEdit, e.NumEdits)
	for i := uint32(0); i < e.NumEdits; i++ {
		if _, err := rdr.Read(readBuffer); err == nil {
			e.Edits[i] = EditListEdit{
				binary.BigEndian.Uint32(readBuffer[0:4]),
				int32(binary.BigEndian.Uint32(readBuffer[4:8])),
				parse32BitFixed(readBuffer[8:12]),
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

type VideoMediaHeaderAtom struct {
	*BasicAtom
	Version      byte   `qtff:" "`
	Flags        []byte `qtff:"3"`
	GraphicsMode uint16 `qtff:" "`
	Opcolor      []byte `qtff:"6"`
}

type SoundMediaHeaderAtom struct {
	*BasicAtom
	Version  byte   `qtff:" "`
	Flags    []byte `qtff:"3"`
	Balance  uint16 `qtff:" "`
	Reserved uint16 `qtff:" "`
}

type SampleTableAtom struct {
	*BasicAtom
}

func (a *SampleTableAtom) Leaf() bool {
	return false
}

// They reall call it that!
type DataInformationAtom struct {
	*BasicAtom
}

func (a *DataInformationAtom) Leaf() bool {
	return false
}

type DataReferenceAtom struct {
	*BasicAtom
	Version    byte   `qtff:" "`
	Flags      []byte `qtff:"3"`
	NumEntries uint32 `qtff:" "`
}

func (a *DataReferenceAtom) Leaf() bool {
	return false
}

type DataReferenceAliasAtom struct {
	*BasicAtom
	Version byte   `qtff:" "`
	Flags   []byte `qtff:"3"`
	Data    []byte
}

func (a *DataReferenceAliasAtom) parseRemainingData(rdr io.Reader) error {
	var err error
	a.Data, err = ioutil.ReadAll(rdr)
	return err
}

type DataReferenceURLAtom struct {
	*BasicAtom
	Version byte   `qtff:" "`
	Flags   []byte `qtff:"3"`
	Data    []byte
	URL     string
}

func (a *DataReferenceURLAtom) parseRemainingData(rdr io.Reader) error {
	var err error
	a.Data, err = ioutil.ReadAll(rdr)
	if err == nil {
		var index = bytes.Index(a.Data, []byte{0})
		if index > 0 {
			a.URL = string(a.Data[0:index])
		} else {
			a.URL = string(a.Data)
		}
	}
	return err
}
