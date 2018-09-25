package serializer

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/void616/gm-sumus-lib/amount"
)

// NewDeserializer instance
func NewDeserializer(data []byte) *Deserializer {
	return &Deserializer{
		buf: bytes.NewBuffer(data),
		err: nil,
	}
}

// Deserializer data
type Deserializer struct {
	buf *bytes.Buffer
	err error
}

// ---

// Error if any occured
func (s *Deserializer) Error() error {
	return s.err
}

// ---

// GetByte ...
func (s *Deserializer) GetByte() byte {
	if s.err == nil {
		v, err := s.buf.ReadByte()
		if err == nil {
			return v
		}
		s.err = err
	}
	return byte(0)
}

// GetBytes ...
func (s *Deserializer) GetBytes(n int) []byte {
	if s.err == nil {
		v := make([]byte, n)
		cnt, err := s.buf.Read(v)
		if err == nil {
			if cnt == n {
				return v
			}
			s.err = fmt.Errorf("Didn't read specified amount of bytes. Got %v, expected %v", cnt, n)
		} else {
			s.err = err
		}
	}
	return nil
}

// GetUint16 ...
func (s *Deserializer) GetUint16() uint16 {
	if s.err == nil {
		b := s.GetBytes(2)
		if b != nil {
			return uint16(s.shiftInt(b).Uint64() & 0xFFFF)
		}
	}
	return uint16(0)
}

// GetUint32 ...
func (s *Deserializer) GetUint32() uint32 {
	if s.err == nil {
		b := s.GetBytes(4)
		if b != nil {
			return uint32(s.shiftInt(b).Uint64() & 0xFFFFFFFF)
		}
	}
	return uint32(0)
}

// GetUint64 ...
func (s *Deserializer) GetUint64() uint64 {
	if s.err == nil {
		b := s.GetBytes(8)
		if b != nil {
			return s.shiftInt(b).Uint64()
		}
	}
	return uint64(0)
}

// GetString64 ...
func (s *Deserializer) GetString64() string {
	const max = 64

	if s.err == nil {
		b := s.GetBytes(max)
		if b != nil {
			to := max
			for i, v := range b {
				if v == 0 {
					to = i
					break
				}
			}
			return string(b[:to])
		}
	}
	return ""
}

// GetAmount ...
func (s *Deserializer) GetAmount() *amount.Amount {

	// must be even
	const imax = 10
	const fmax = 18

	if s.err == nil {

		var err error

		var sign = s.GetByte()
		var fragPart = s.GetBytes(fmax / 2)
		var intPart = s.GetBytes(imax / 2)

		// check sign
		strSign := ""
		if sign > 0 {
			if sign > 1 {
				s.err = fmt.Errorf("Amount sign byte has invalid value: %v", sign)
			}
			strSign = "-"
		}

		// unflip frag part
		strFrag := ""
		if s.err == nil {
			strFrag, err = unflipAmountStringLikeAShit(fragPart)
			if err != nil {
				s.err = err
			}
		}

		// unflip int part
		strInt := ""
		if s.err == nil {
			strInt, err = unflipAmountStringLikeAShit(intPart)
			if err != nil {
				s.err = err
			}
		}

		// try parse amount
		if s.err == nil {
			together := fmt.Sprintf("%v%v.%v", strSign, strInt, strFrag)
			ret := amount.NewFloatString(together)
			if ret != nil {
				return ret
			}
			s.err = fmt.Errorf("Failed to parse amount from: %v", together)
		}
	}

	return nil
}

func (s *Deserializer) shiftInt(b []byte) *big.Int {
	x := big.NewInt(0)
	ret := big.NewInt(0)
	for i, v := range b {
		x = x.SetUint64(uint64(v) & 0xFF)
		x = x.Lsh(x, uint(i)*8)
		ret = ret.Or(ret, x)
	}
	return ret
}

// Convert some kind of a shit into string: [0x78 0x56 .. 0x34 0x12] => "1234...5678"
func unflipAmountStringLikeAShit(b []byte) (string, error) {
	if b == nil || len(b) == 0 {
		return "", fmt.Errorf("Buffer is null or empty")
	}

	// reverse array
	tmp := make([]byte, len(b))
	for i := 0; i < len(b); i++ {
		tmp[i] = b[len(b)-i-1]
	}

	// to the hex
	return hex.EncodeToString(tmp), nil
}