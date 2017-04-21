package ntgo

import (
	"testing"
	"bytes"
	"reflect"
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
	result := BuildDouble(0.5)
	var expected = &ValueDouble{
		Value: 0.5,
		RawValue: Float64ToBytes(0.5),
	}
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("Expected %s but got %s", expected, result)
	}
}

func TestBuildString(t *testing.T) {
	result := BuildString("test")
	var expected = &ValueString{
		Value: "test",
		RawValue: []byte{0x04, 0x74, 0x65, 0x73, 0x74},
	}
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("Expected %s but got %s", expected, result)
	}
}

func TestBuildRaw(t *testing.T) {
	result := BuildRaw([]byte{0x50, 0x21})
	var expected = &ValueRaw{
		Value: []byte{0x50, 0x21},
		RawValue: []byte{0x02, 0x50, 0x21},
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
	result := BuildDoubleArray([]*ValueDouble{BuildDouble(0.5)})
	var expected = &ValueDoubleArray{
		index: uint8(1),
		elements: []*ValueDouble{BuildDouble(0.5)},
	}
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("Expected %s but got %s", expected, result)
	}
}

func TestBuildStringArray(t *testing.T) {
	result := BuildStringArray([]*ValueString{BuildString("test")})
	var expected = &ValueStringArray{
		index: uint8(1),
		elements: []*ValueString{BuildString("test")},
	}
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("Expected %s but got %s", expected, result)
	}
}

func TestDecodeEntryBoolean(t *testing.T) {
	entryBytes := []byte{
		byte(TypeBoolean),
		BoolTrue,
	}
	result, err := DecodeEntry(bytes.NewBuffer(entryBytes))
	if err != nil {
		t.Fatal("Unexpected error: ", err)
	}
	var expected = &ValueBoolean{
		Value: true,
		RawValue: []byte{BoolTrue},
	}
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("Expected %s but got %s", expected, result)
	}
}

func TestDecodeEntryTypeBoolean(t *testing.T) {
	result, err := DecodeEntryType(bytes.NewBuffer([]byte{byte(TypeBoolean)}))

	if err != nil {
		t.Fatal("Unexpected error: ", err)
	}
	var expected EntryType = TypeBoolean
	if result != expected {
		t.Fatalf("Expected %s but got %s", expected, result)
	}
}

func TestDecodeEntryTypeDouble(t *testing.T) {
	result, err := DecodeEntryType(bytes.NewBuffer([]byte{byte(TypeDouble)}))

	if err != nil {
		t.Fatal("Unexpected error: ", err)
	}
	var expected EntryType = TypeDouble
	if result != expected {
		t.Fatalf("Expected %s but got %s", expected, result)
	}
}

func TestDecodeEntryTypeString(t *testing.T) {
	result, err := DecodeEntryType(bytes.NewBuffer([]byte{byte(TypeString)}))

	if err != nil {
		t.Fatal("Unexpected error: ", err)
	}
	var expected EntryType = TypeString
	if result != expected {
		t.Fatalf("Expected %s but got %s", expected, result)
	}
}

func TestDecodeEntryTypeRawData(t *testing.T) {
	result, err := DecodeEntryType(bytes.NewBuffer([]byte{byte(TypeRawData)}))

	if err != nil {
		t.Fatal("Unexpected error: ", err)
	}
	var expected EntryType = TypeRawData
	if result != expected {
		t.Fatalf("Expected %s but got %s", expected, result)
	}
}

func TestDecodeEntryTypeBooleanArray(t *testing.T) {
	result, err := DecodeEntryType(bytes.NewBuffer([]byte{byte(TypeBooleanArr)}))

	if err != nil {
		t.Fatal("Unexpected error: ", err)
	}
	var expected EntryType = TypeBooleanArr
	if result != expected {
		t.Fatalf("Expected %s but got %s", expected, result)
	}
}

func TestDecodeEntryTypeDoubleArray(t *testing.T) {
	result, err := DecodeEntryType(bytes.NewBuffer([]byte{byte(TypeDoubleArr)}))

	if err != nil {
		t.Fatal("Unexpected error: ", err)
	}
	var expected EntryType = TypeDoubleArr
	if result != expected {
		t.Fatalf("Expected %s but got %s", expected, result)
	}
}

func TestDecodeEntryTypeStringArray(t *testing.T) {
	result, err := DecodeEntryType(bytes.NewBuffer([]byte{byte(TypeStringArr)}))

	if err != nil {
		t.Fatal("Unexpected error: ", err)
	}
	var expected EntryType = TypeStringArr
	if result != expected {
		t.Fatalf("Expected %s but got %s", expected, result)
	}
}
