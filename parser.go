package qtff

import (
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"reflect"
)

func Parse(rdr io.Reader) ([]Atom, error) {
	atoms := make([]Atom, 0, 5)
	var err error
	var a Atom
	for err == nil {
		if a, err = parseNext(rdr); err == nil {
			fmt.Println(a, a.Length(), string(a.Type()))
			atoms = append(atoms, a)
		}
	}
	return atoms, err
}

func parseNext(rdr io.Reader) (Atom, error) {
	if a, i, err := readAtomHeader(rdr); err == nil {
		// Limit the rest of the reads
		if a.Length() == 0 {
			rdr = &io.LimitedReader{rdr, math.MaxInt64}
		} else {
			rdr = &io.LimitedReader{rdr, int64(a.Length() - uint64(i))}
		}

		// Find out if atom is of a known type
		atom := upgradeType(a)

		// Parse special headers if there are any
		// fmt.Println(rdr.(*io.LimitedReader).N)
		if err = parseSpecialHeaders(rdr, atom); err != nil {
			if err != io.EOF {
				return atom, err
			}
		}

		// Handle the remaining data.  An atom either has
		// child atoms or data, never both
		if !atom.Leaf() {
			if c, err := Parse(rdr); err == nil {
				a.ChildAtoms = c
			}
		}

		// Discard remaining data
		// FIXME: we need this data eventually
		io.Copy(ioutil.Discard, rdr)

		return atom, err
	} else {
		return nil, err
	}
}

func readAtomHeader(rdr io.Reader) (atom *BasicAtom, bytesRead int, err error) {
	atom = &BasicAtom{}
	atom.typ = make([]byte, 4)
	var sizeBlock = make([]byte, 8)
	var readExtendedSize = false
	// Read simple size
	if _, err = rdr.Read(sizeBlock[0:4]); err == nil {
		atom.length = uint64(binary.BigEndian.Uint32(sizeBlock[0:4]))
		if atom.length == 1 {
			readExtendedSize = true
		}
	} else {
		return
	}
	bytesRead += 4

	// Read type header
	if _, err = rdr.Read(atom.typ); err != nil {
		return
	}
	bytesRead += 4

	// Read extended size
	if readExtendedSize {
		if _, err = rdr.Read(sizeBlock); err == nil {
			atom.length = binary.BigEndian.Uint64(sizeBlock)
		} else {
			return
		}
		bytesRead += 8
	}

	return
}

func upgradeType(b *BasicAtom) Atom {
	switch string(b.Type()) {
	case "ftyp":
		return &FileTypeAtom{BasicAtom: b}
	default:
		return b
	}
}

func parseSpecialHeaders(rdr io.Reader, atom Atom) error {
	value := reflect.ValueOf(atom)
	writeValue := value.Elem()
	if value.Kind() == reflect.Ptr {
		value = reflect.Indirect(value)
	}
	if value.Kind() == reflect.Struct {
		var i = 0
		for n := value.NumField(); i < n; i++ {
			field := writeValue.Field(i)
			fieldDesc := value.Type().Field(i)
			if tag := fieldDesc.Tag.Get("qtff"); tag != "" {
				readBlock := make([]byte, 16)
				// fmt.Println(field)
				switch field.Kind() {
				case reflect.Uint32:
					if _, err := rdr.Read(readBlock[0:4]); err == nil {
						field.Set(reflect.ValueOf(binary.BigEndian.Uint32(readBlock[0:4])))
						// fmt.Println(atom, field)
					} else {
						// return err
					}
				}
			}
		}
	}
	return nil
}
