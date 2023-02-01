package marc

import (
	"errors"
	"fmt"
	"strconv"
)

// Leader represents the leader of the MARC record.
type Leader struct {
	raw           []byte
	dataOffset    int
	Status        byte // 05 byte position
	Type          byte // 06
	BibLevel      byte // 07
	Control       byte // 08
	EncodingLevel byte // 17
	Form          byte // 18
	Multipart     byte // 19
}

// NewLeader creates a Leader from the data in the MARC record.
func NewLeader(bytes []byte) (Leader, error) {
	if len(bytes) != leaderLength {
		return Leader{}, errors.New("incomplete leader")
	}

	// A typical good leader value is: "01848nam a2200385 i 4500"
	// where as a bad value would be.: "ZZZZZnamZa22ZZZZZzZZ4500"
	offset, err := strconv.Atoi(string(bytes[offsetStart:offsetEnd]))
	if err != nil {
		msg := fmt.Sprintf("Could not determine data offset from leader (%s)", string(bytes))
		err = errors.New(msg)
		offset = -1
	}

	leader := Leader{
		raw:           bytes,
		dataOffset:    offset,
		Status:        bytes[5],
		Type:          bytes[6],
		BibLevel:      bytes[7],
		Control:       bytes[8],
		EncodingLevel: bytes[17],
		Form:          bytes[18],
		Multipart:     bytes[19],
	}
	return leader, err
}

func (l Leader) String() string {
	return fmt.Sprintf("=LDR  %s", string(l.raw))
}

func (l Leader) Raw() string {
	return string(l.raw)
}
