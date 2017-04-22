package ntgo

import (
	"bytes"
	"io"
	"encoding/binary"
	"math"
	"errors"
)

var (
	ErrULEBDecodeDataInvalid = errors.New("uleb: could not decode invalid uleb encoding")
)

// DecodeULEB128 returns the ULEB128-encoded Value.
func DecodeULEB128(r io.Reader) (uint32, error) {
	var result uint32 = 0
	var shift uint32 = 0
	var currByte = [1]byte{0x80}
	for currByte[0]&0x80 == 0x80 {
		_, readErr := io.ReadFull(r, currByte[:])
		if readErr != nil {
			return 0, ErrULEBDecodeDataInvalid
		}
		result |= uint32(currByte[0]&0x7f) << shift
		shift+=7
	}
	return result, nil
}

// DecodeSaveULEB128 returns the ULEB128-encoded Value and the ULEB128 data.
func DecodeAndSaveULEB128(r io.Reader) (uint32, []byte, error) {
	data := []byte{}
	var result uint32 = 0
	var shift uint32 = 0
	currByte := [1]byte{0x80}
	for currByte[0]&0x80 == 0x80 {
		_, readErr := io.ReadFull(r, currByte[:])
		if readErr != nil {
			return 0, data, ErrULEBDecodeDataInvalid
		}
		result |= uint32(currByte[0]&0x7f) << shift
		shift+=7
		data = append(data, currByte[0])
	}
	return result, data, nil
}

// Encodes data into ULEB128 Value as a byte array
func EncodeULEB128(data uint32) []byte {
	remaining := data >> 7
	var buf = new(bytes.Buffer)
	for remaining != 0 {
		buf.WriteByte(byte(data&0x7f | 0x80))
		data = remaining
		remaining >>= 7
	}
	buf.WriteByte(byte(data & 0x7f))
	return buf.Bytes()
}

// BytesToFloat64 converts bytes to Float64
func BytesToFloat64(bytes []byte) float64 {
	bits := binary.BigEndian.Uint64(bytes)
	float := math.Float64frombits(bits)
	return float
}

func Float64ToBytes(value float64) []byte {
	bits := math.Float64bits(value)
	floatBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(floatBytes, bits)
	return floatBytes
}