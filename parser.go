package qtff

import (
	"io"
	"io/ioutil"
	"encoding/binary"
	"fmt"
)

func Parse(rdr io.Reader)([]Atom, error) {
	atoms := make([]Atom, 0, 5)
	var err error
	var a Atom
	for err == nil {
		if a, err = parseNext(rdr); err == nil {
			fmt.Println(a.(*FileTypeAtom).Length, string(a.(*FileTypeAtom).Type))
			atoms = append(atoms, a)
		}
	}
	if err == io.EOF {
		err = nil
	}
	return atoms, err
}

func parseNext(rdr io.Reader)(Atom, error) {
	if size, typ, i, err := readAtomHeader(rdr); err == nil {
		a := &FileTypeAtom{size, typ}
		// the last atom has a size of 0
		if size == 0 {
			io.Copy(ioutil.Discard, rdr)
		} else {
			b := make([]byte, size - uint64(i))
			_, err = rdr.Read(b)
		}
		return a, err
	} else {
		return nil, err
	}
}

func readAtomHeader(rdr io.Reader)(size uint64, typ []byte, bytesRead int, err error){
	var sizeBlock = make([]byte, 8)
	typ = make([]byte, 4)
	var readExtendedSize = false
	// Read simple size
	if _, err = rdr.Read(sizeBlock[0:4]); err == nil {
		size = uint64(binary.BigEndian.Uint32(sizeBlock[0:4]))
		if size == 1 {
			readExtendedSize = true
		}
	} else {
		return
	}
	bytesRead += 4

	// Read type header
	if _, err = rdr.Read(typ); err != nil {
		return
	}
	bytesRead += 4

	// Read extended size
	if readExtendedSize {
		if _, err = rdr.Read(sizeBlock); err == nil {
			size = binary.BigEndian.Uint64(sizeBlock)
		} else {
			return
		}
	}
	bytesRead += 8
	
	return
}