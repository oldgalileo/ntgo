package ntgo

import (
	"testing"
	"bytes"
	"reflect"
	"fmt"
)

func TestBuildBoolean(t *testing.T) {
	result := BuildBoolean(true)
	var expected = &ValueBoolean{
		Value: true,
		RawValue: []byte{BoolTrue},
	}
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("Expected %s but got %s", expected, result)
	}
}

func TestBuildDouble(t *testing.T) {
	result := BuildDouble(0.8)
	var expected = &ValueDouble{
		Value: 0.8,
		RawValue: []byte{0x3f,0xe9,0x99,0x99,0x99,0x99,0x99,0x9a}, // Value of float64(0.8)
	}
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("Expected %s but got %s", expected, result)
	}
}

func TestBuildString(t *testing.T) {
	result := BuildString("test")
	var expected = &ValueString{
		Value: "test",
		RawValue: []byte{0x04, 0x74, 0x65, 0x73, 0x74}, // Value of "test"
	}
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("Expected %s but got %s", expected, result)
	}
}

func TestBuildRaw(t *testing.T) {
	result := BuildRaw([]byte{0x50, 0x21})
	var expected = &ValueRaw{
		Value: []byte{0x50, 0x21},
		RawValue: []byte{0x02, 0x50, 0x21}, // Meaningless Value
	}
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("Expected %s but got %s", expected, result)
	}
}

func TestBuildBooleanArray(t *testing.T) {
	result := BuildBooleanArray([]*ValueBoolean{BuildBoolean(true)})
	var expected = &ValueBooleanArray{
		index: uint8(1),
		elements: []*ValueBoolean{BuildBoolean(true)},
	}
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("Expected %s but got %s", expected, result)
	}
}

func TestBuildDoubleArray(t *testing.T) {
	result := BuildDoubleArray([]*ValueDouble{BuildDouble(49.04)})
	var expected = &ValueDoubleArray{
		index: uint8(1),
		elements: []*ValueDouble{BuildDouble(49.04)},
	}
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("Expected %s but got %s", expected, result)
	}
}

func TestBuildStringArray(t *testing.T) {
	result := BuildStringArray([]*ValueString{BuildString("str")})
	var expected = &ValueStringArray{
		index: uint8(1),
		elements: []*ValueString{BuildString("str")},
	}
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("Expected %s but got %s", expected, result)
	}
}

func TestDecodeEntryBoolean(t *testing.T) {
	entryBytes := []byte{byte(TypeBoolean),BoolFalse}
	result, err := DecodeEntry(bytes.NewBuffer(entryBytes))
	if err != nil {
		t.Fatalf("Unexpected error! %s", err)
	}
	var expected = &ValueBoolean{
		Value: false,
		RawValue: []byte{BoolFalse},
	}
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("Expected %s but got %s", expected, result)
	}
}

func TestDecodeEntryString(t *testing.T) {
	entryBytes := []byte{byte(TypeString),0x05,0x6f,0x74,0x68,0x65,0x72} // Value of other
	result, err := DecodeEntry(bytes.NewBuffer(entryBytes))
	if err != nil {
		t.Fatalf("Unexpected error! %s", err)
	}
	var expected = &ValueString{
		Value: "other",
		RawValue: entryBytes[1:],
	}
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("Expected %s but got %s", expected, result)
	}
}

func TestDecodeEntryDouble(t *testing.T) {
	entryBytes := []byte{byte(TypeDouble),0x3f,0xe0,0x00,0x00,0x00,0x00,0x00,0x00}
	result, err := DecodeEntry(bytes.NewBuffer(entryBytes))
	if err != nil {
		t.Fatalf("Unexpected error! %s", err)
	}
	var expected = &ValueDouble{
		Value: 0.5,
		RawValue: entryBytes[1:],
	}
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("Expected %s but got %s", expected, result)
	}
}

func TestDecodeEntryBooleanArray(t *testing.T) {
	entryBytes := []byte{byte(TypeBooleanArr),byte(uint8(1)),0x01}
	result, err := DecodeEntry(bytes.NewBuffer(entryBytes))
	if err != nil {
		t.Fatalf("Unexpected error! %s", err)
	}
	var expected = &ValueBooleanArray{
		index: uint8(1),
		elements: []*ValueBoolean{
			BuildBoolean(true),
		},
	}
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("Expected %s but got %s", expected, result)
	}
}

func TestDecodeEntryDoubleArray(t *testing.T) {
	entryBytes := []byte{byte(TypeDoubleArr),byte(uint8(1)),0x3f,0xf2,0xe1,0x47,0xae,0x14,0x7a,0xe1}
	result, err := DecodeEntry(bytes.NewBuffer(entryBytes))
	if err != nil {
		t.Fatalf("Unexpected error! %s", err)
	}
	var expected = &ValueDoubleArray{
		index: uint8(1),
		elements: []*ValueDouble{
			BuildDouble(1.18),
		},
	}
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("Expected %s but got %s", expected, result)
	}
}

func TestDecodeEntryStringArray(t *testing.T) {
	entryBytes := []byte{byte(TypeStringArr), byte(uint8(1)),0x05,0x61,0x72,0x72,0x61,0x79}
	result, err := DecodeEntry(bytes.NewBuffer(entryBytes))
	if err != nil {
		t.Fatalf("Unexpected error! %s", err)
	}
	var expected = &ValueStringArray{
		index: uint8(1),
		elements: []*ValueString{
			BuildString("array"),
		},
	}
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("Expected %s but got %s", expected, result)
	}
}

func TestDecodeEntryTypeBoolean(t *testing.T) {
	result, err := DecodeEntryType(bytes.NewBuffer([]byte{byte(TypeBoolean)}))
	if err != nil {
		t.Fatalf("Unexpected error! %s", err)
	}
	var expected EntryType = TypeBoolean
	if result != expected {
		t.Fatalf("Expected %s but got %s", expected, result)
	}
}

func TestDecodeEntryTypeDouble(t *testing.T) {
	result, err := DecodeEntryType(bytes.NewBuffer([]byte{byte(TypeDouble)}))

	if err != nil {
		t.Fatalf("Unexpected error! %s", err)
	}
	var expected EntryType = TypeDouble
	if result != expected {
		t.Fatalf("Expected %s but got %s", expected, result)
	}
}

func TestDecodeEntryTypeString(t *testing.T) {
	result, err := DecodeEntryType(bytes.NewBuffer([]byte{byte(TypeString)}))

	if err != nil {
		t.Fatalf("Unexpected error! %s", err)
	}
	var expected EntryType = TypeString
	if result != expected {
		t.Fatalf("Expected %s but got %s", expected, result)
	}
}

func TestDecodeEntryTypeRawData(t *testing.T) {
	result, err := DecodeEntryType(bytes.NewBuffer([]byte{byte(TypeRawData)}))
	if err != nil {
		t.Fatalf("Unexpected error! %s", err)
	}
	var expected EntryType = TypeRawData
	if result != expected {
		t.Fatalf("Expected %s but got %s", expected, result)
	}
}

func TestDecodeEntryTypeBooleanArray(t *testing.T) {
	result, err := DecodeEntryType(bytes.NewBuffer([]byte{byte(TypeBooleanArr)}))
	if err != nil {
		t.Fatalf("Unexpected error! %s", err)
	}
	var expected EntryType = TypeBooleanArr
	if result != expected {
		t.Fatalf("Expected %s but got %s", expected, result)
	}
}

func TestDecodeEntryTypeDoubleArray(t *testing.T) {
	result, err := DecodeEntryType(bytes.NewBuffer([]byte{byte(TypeDoubleArr)}))

	if err != nil {
		t.Fatalf("Unexpected error! %s", err)
	}
	var expected EntryType = TypeDoubleArr
	if result != expected {
		t.Fatalf("Expected %s but got %s", expected, result)
	}
}

func TestDecodeEntryTypeStringArray(t *testing.T) {
	result, err := DecodeEntryType(bytes.NewBuffer([]byte{byte(TypeStringArr)}))

	if err != nil {
		t.Fatalf("Unexpected error! %s", err)
	}
	var expected EntryType = TypeStringArr
	if result != expected {
		t.Fatalf("Expected %s but got %s", expected, result)
	}
}
