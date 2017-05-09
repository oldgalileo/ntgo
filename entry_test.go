package ntgo

import (
	"testing"
	"bytes"
	"reflect"
)

func TestBuildBoolean(t *testing.T) {
	result := BuildBoolean(true)
	var expected = &ValueBoolean{
		Value:    true,
		RawValue: []byte{BoolTrue},
	}
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("Expected %s but got %s", expected, result)
	}
}

func TestBuildDouble(t *testing.T) {
	result := BuildDouble(0.8)
	var expected = &ValueDouble{
		Value:    0.8,
		RawValue: []byte{0x3f, 0xe9, 0x99, 0x99, 0x99, 0x99, 0x99, 0x9a}, // Value of float64(0.8)
	}
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("Expected %s but got %s", expected, result)
	}
}

func TestBuildString(t *testing.T) {
	result := BuildString("test")
	var expected = &ValueString{
		Value:    "test",
		RawValue: []byte{0x04, 0x74, 0x65, 0x73, 0x74}, // Value of "test"
	}
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("Expected %s but got %s", expected, result)
	}
}

func TestBuildRaw(t *testing.T) {
	result := BuildRaw([]byte{0x50, 0x21})
	var expected = &ValueRaw{
		Value:    []byte{0x50, 0x21},
		RawValue: []byte{0x02, 0x50, 0x21}, // Meaningless Value
	}
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("Expected %s but got %s", expected, result)
	}
}

func TestBuildBooleanArray(t *testing.T) {
	result := BuildBooleanArray([]*ValueBoolean{BuildBoolean(true)})
	var expected = &ValueBooleanArray{
		elements: []*ValueBoolean{BuildBoolean(true)},
	}
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("Expected %s but got %s", expected, result)
	}
}

func TestBuildDoubleArray(t *testing.T) {
	result := BuildDoubleArray([]*ValueDouble{BuildDouble(49.04)})
	var expected = &ValueDoubleArray{
		elements: []*ValueDouble{BuildDouble(49.04)},
	}
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("Expected %s but got %s", expected, result)
	}
}

func TestBuildStringArray(t *testing.T) {
	result := BuildStringArray([]*ValueString{BuildString("str")})
	var expected = &ValueStringArray{
		elements: []*ValueString{BuildString("str")},
	}
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("Expected %s but got %s", expected, result)
	}
}

func TestDecodeEntryValueBoolean(t *testing.T) {
	entryBytes := []byte{byte(EntryTypeBoolean), BoolFalse}
	result, _, err := DecodeEntryValueAndType(bytes.NewBuffer(entryBytes))
	if err != nil {
		t.Fatalf("Unexpected error! %s", err)
	}
	var expected = &ValueBoolean{
		Value:    false,
		RawValue: []byte{BoolFalse},
	}
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("Expected %s but got %s", expected, result)
	}
}

func TestDecodeEntryValueString(t *testing.T) {
	entryBytes := []byte{byte(EntryTypeString), 0x05, 0x6f, 0x74, 0x68, 0x65, 0x72} // Value of other
	result, _, err := DecodeEntryValueAndType(bytes.NewBuffer(entryBytes))
	if err != nil {
		t.Fatalf("Unexpected error! %s", err)
	}
	var expected = &ValueString{
		Value:    "other",
		RawValue: entryBytes[1:],
	}
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("Expected %s but got %s", expected, result)
	}
}

func TestDecodeEntryValueDouble(t *testing.T) {
	entryBytes := []byte{byte(EntryTypeDouble), 0x3f, 0xe0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	result, _, err := DecodeEntryValueAndType(bytes.NewBuffer(entryBytes))
	if err != nil {
		t.Fatalf("Unexpected error! %s", err)
	}
	var expected = &ValueDouble{
		Value:    0.5,
		RawValue: entryBytes[1:],
	}
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("Expected %s but got %s", expected, result)
	}
}

func TestDecodeEntryValueBooleanArray(t *testing.T) {
	entryBytes := []byte{byte(EntryTypeBooleanArr), byte(uint8(1)), 0x01}
	result, _, err := DecodeEntryValueAndType(bytes.NewBuffer(entryBytes))
	if err != nil {
		t.Fatalf("Unexpected error! %s", err)
	}
	var expected = &ValueBooleanArray{
		elements: []*ValueBoolean{
			BuildBoolean(true),
		},
	}
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("Expected %s but got %s", expected, result)
	}
}

func TestEntryValueBooleanArrayGetSafe(t *testing.T) {
	var testEntry = &ValueBoolean{
		Value:    true,
		RawValue: []byte{BoolTrue},
	}
	array := BuildBooleanArray([]*ValueBoolean{
		BuildBoolean(true), BuildBoolean(false),
		BuildBoolean(false), testEntry,
	})
	entry, err := array.Get(3)
	if err != nil {
		t.Fatalf("Unexpected error! %s", err)
	}
	if !reflect.DeepEqual(entry, testEntry) {
		t.Fatalf("Expected %s but got %s", testEntry, entry)
	}
}

func TestEntryValueBooleanArrayGetFail(t *testing.T) {
	array := BuildBooleanArray([]*ValueBoolean{
		BuildBoolean(true), BuildBoolean(false),
		BuildBoolean(false), BuildBoolean(false),
	})
	_, err := array.Get(4)
	if err == nil {
		t.Fatal("Expected error but received nil")
	}
	if err != ErrArrayIndexOutOfBounds {
		t.Fatalf("Expected error \"%s\" but received \"%s\"", ErrArrayIndexOutOfBounds, err)
	}
}

func TestDecodeEntryValueDoubleArray(t *testing.T) {
	entryBytes := []byte{byte(EntryTypeDoubleArr), byte(uint8(1)), 0x3f, 0xf2, 0xe1, 0x47, 0xae, 0x14, 0x7a, 0xe1}
	result, _, err := DecodeEntryValueAndType(bytes.NewBuffer(entryBytes))
	if err != nil {
		t.Fatalf("Unexpected error! %s", err)
	}
	var expected = &ValueDoubleArray{
		elements: []*ValueDouble{
			BuildDouble(1.18),
		},
	}
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("Expected %s but got %s", expected, result)
	}
}

func TestEntryValueDoubleArrayGetSafe(t *testing.T) {
	var testEntry = BuildDouble(0.3)
	array := BuildDoubleArray([]*ValueDouble{
		BuildDouble(0.0), BuildDouble(0.1),
		BuildDouble(0.2), testEntry,
	})
	entry, err := array.Get(3)
	if err != nil {
		t.Fatalf("Unexpected error! %s", err)
	}
	if !reflect.DeepEqual(entry, testEntry) {
		t.Fatalf("Expected %s but got %s", testEntry, entry)
	}
}

func TestEntryValueDoubleArrayGetFail(t *testing.T) {
	array := BuildDoubleArray([]*ValueDouble{
		BuildDouble(0.0), BuildDouble(0.1),
		BuildDouble(0.2), BuildDouble(0.3),
	})
	_, err := array.Get(4)
	if err == nil {
		t.Fatal("Expected error but received nil")
	}
	if err != ErrArrayIndexOutOfBounds {
		t.Fatalf("Expected error \"%s\" but received \"%s\"", ErrArrayIndexOutOfBounds, err)
	}
}

func TestDecodeEntryValueStringArray(t *testing.T) {
	entryBytes := []byte{byte(EntryTypeStringArr), byte(uint8(1)), 0x05, 0x61, 0x72, 0x72, 0x61, 0x79}
	result, _, err := DecodeEntryValueAndType(bytes.NewBuffer(entryBytes))
	if err != nil {
		t.Fatalf("Unexpected error! %s", err)
	}
	var expected = &ValueStringArray{
		elements: []*ValueString{
			BuildString("array"),
		},
	}
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("Expected %s but got %s", expected, result)
	}
}

func TestEntryValueStringArrayGetSafe(t *testing.T) {
	var testEntry = BuildString("test3")
	array := BuildStringArray([]*ValueString{
		BuildString("test0"), BuildString("test1"),
		BuildString("test2"), testEntry,
	})
	entry, err := array.Get(3)
	if err != nil {
		t.Fatalf("Unexpected error! %s", err)
	}
	if !reflect.DeepEqual(entry, testEntry) {
		t.Fatalf("Expected %s but got %s", testEntry, entry)
	}
}

func TestEntryValueStringArrayGetFail(t *testing.T) {
	array := BuildStringArray([]*ValueString{
		BuildString("test0"), BuildString("test1"),
		BuildString("test2"), BuildString("test3"),
	})
	_, err := array.Get(4)
	if err == nil {
		t.Fatal("Expected error but received nil")
	}
	if err != ErrArrayIndexOutOfBounds {
		t.Fatalf("Expected error \"%s\" but received \"%s\"", ErrArrayIndexOutOfBounds, err)
	}
}

func TestDecodeEntryTypeBoolean(t *testing.T) {
	result, err := DecodeEntryType(bytes.NewBuffer([]byte{byte(EntryTypeBoolean)}))
	if err != nil {
		t.Fatalf("Unexpected error! %s", err)
	}
	var expected EntryType = EntryTypeBoolean
	if result != expected {
		t.Fatalf("Expected %s but got %s", expected, result)
	}
}

func TestDecodeEntryTypeDouble(t *testing.T) {
	result, err := DecodeEntryType(bytes.NewBuffer([]byte{byte(EntryTypeDouble)}))

	if err != nil {
		t.Fatalf("Unexpected error! %s", err)
	}
	var expected EntryType = EntryTypeDouble
	if result != expected {
		t.Fatalf("Expected %s but got %s", expected, result)
	}
}

func TestDecodeEntryTypeString(t *testing.T) {
	result, err := DecodeEntryType(bytes.NewBuffer([]byte{byte(EntryTypeString)}))

	if err != nil {
		t.Fatalf("Unexpected error! %s", err)
	}
	var expected EntryType = EntryTypeString
	if result != expected {
		t.Fatalf("Expected %s but got %s", expected, result)
	}
}

func TestDecodeEntryTypeRawData(t *testing.T) {
	result, err := DecodeEntryType(bytes.NewBuffer([]byte{byte(EntryTypeRawData)}))
	if err != nil {
		t.Fatalf("Unexpected error! %s", err)
	}
	var expected EntryType = EntryTypeRawData
	if result != expected {
		t.Fatalf("Expected %s but got %s", expected, result)
	}
}

func TestDecodeEntryTypeBooleanArray(t *testing.T) {
	result, err := DecodeEntryType(bytes.NewBuffer([]byte{byte(EntryTypeBooleanArr)}))
	if err != nil {
		t.Fatalf("Unexpected error! %s", err)
	}
	var expected EntryType = EntryTypeBooleanArr
	if result != expected {
		t.Fatalf("Expected %s but got %s", expected, result)
	}
}

func TestDecodeEntryTypeDoubleArray(t *testing.T) {
	result, err := DecodeEntryType(bytes.NewBuffer([]byte{byte(EntryTypeDoubleArr)}))

	if err != nil {
		t.Fatalf("Unexpected error! %s", err)
	}
	var expected EntryType = EntryTypeDoubleArr
	if result != expected {
		t.Fatalf("Expected %s but got %s", expected, result)
	}
}

func TestDecodeEntryTypeStringArray(t *testing.T) {
	result, err := DecodeEntryType(bytes.NewBuffer([]byte{byte(EntryTypeStringArr)}))

	if err != nil {
		t.Fatalf("Unexpected error! %s", err)
	}
	var expected EntryType = EntryTypeStringArr
	if result != expected {
		t.Fatalf("Expected %s but got %s", expected, result)
	}
}
