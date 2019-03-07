package s4

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"strconv"
)

// Share holds all information a share has
type Share struct {
	// X is the x value of the share
	X byte
	// RequiredShares holds the threshold this share was generated for
	RequiredShares byte
	// Index is the index of this share
	Index byte
	// Values contains a list of y values, that were computed using (a different) f(x) = y for each secret byte
	Values []byte
}

// DecodeBase64Share will try to decode a Share from a Base64 encoded string
func DecodeBase64Share(b64 string) (Share, error) {

	if b64 == "" {
		return Share{}, errors.New("no share supplied")
	}

	bytestream, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return Share{}, errors.New("invalid Base64 provided")
	}

	vals := make([]byte, 0)

	for i := 3; i < len(bytestream); i++ {
		vals = append(vals, bytestream[i])
	}

	idx := bytestream[0]
	reqShares := bytestream[1]
	x := bytestream[2]

	return Share{
		Index:          idx,
		RequiredShares: reqShares,
		X:              x,
		Values:         vals,
	}, nil
}

// Print will print this share in a human readable form to the standard output
func (s Share) Print() {
	fmt.Printf("Index: %d\nRequired: %d\nX: %d\nValues: %v\n\n", s.Index, s.RequiredShares, s.X, s.Values)
}

// GetBase64 will return this share as a Base64-encoded string
func (s Share) GetBase64() string {
	return base64.StdEncoding.EncodeToString(s.asByteArray())
}

func (s Share) asByteArray() []byte {
	buf := new(bytes.Buffer)

	err := binary.Write(buf, binary.BigEndian, s.Index)
	if err != nil {
		panic("error writing to buffer (0)")
	}

	err = binary.Write(buf, binary.BigEndian, s.RequiredShares)
	if err != nil {
		panic("error writing to buffer (1)")
	}

	err = binary.Write(buf, binary.BigEndian, s.X)
	if err != nil {
		panic("error writing to buffer (2)")
	}

	for i, val := range s.Values {
		err = binary.Write(buf, binary.BigEndian, val)
		if err != nil {
			panic("error writing to buffer (val" + strconv.Itoa(i) + ")")
		}
	}

	return buf.Bytes()
}
