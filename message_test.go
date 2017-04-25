package ntgo

import (
	"testing"
	"bytes"
	"reflect"
)

func TestDecodeDataClientHello(t *testing.T) {
	messageBytes := []byte{0x03, 0x00, 0x04, 0x6e, 0x74, 0x67, 0x6F}
	result, err := DecodeDataClientHello(bytes.NewBuffer(messageBytes))
	if err != nil {
		t.Fatalf("Unexpected error! %s", err)
	}
	var expected = &MessageDataClientHello{
		ProtocVersion: [2]byte{0x03, 0x00},
		Identity: &ValueString{
			Value: "ntgo",
			RawValue: []byte{0x04, 0x6e, 0x74, 0x67, 0x6F},
		},
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Expected %s but got %s", expected, result)
	}
}

func TestDecodeDataProtocVersionUnsupported(t *testing.T) {
	messageBytes := []byte{0x03, 0x00}
	result, err := DecodeDataProtocVersionUnsupported(bytes.NewBuffer(messageBytes))
	if err != nil {
		t.Fatalf("Unexpected error! %s", err)
	}
	var expected = &MessageDataProtocVersionUnsupported{
		SupportedProtoc: [2]byte{0x03, 0x00},
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Expected %s but got %s", expected, result)
	}
}

func TestDecodeDataServerHello(t *testing.T) {
	messageBytes := []byte{0x00, 0x04, 0x6e, 0x74, 0x67, 0x6F}
	result, err := DecodeDataServerHello(bytes.NewBuffer(messageBytes))
	if err != nil {
		t.Fatalf("Unexpected error! %s", err)
	}
	var expected = &MessageDataServerHello{
		Flags: FlagMessageClientNew,
		Identity: &ValueString{
			Value: "ntgo",
			RawValue: []byte{0x04, 0x6e, 0x74, 0x67, 0x6F},
		},
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Expected %s but got %s", expected, result)
	}
}

func TestDecodeDataEntryAssignment(t *testing.T) {
	messageBytes := []byte{0x05, 0x65,0x6e,0x74,0x72,0x79, // Entry Name ("entry")
		byte(TypeBoolean), // Entry Type
		0x50, 0x21, // Unique ID
		0x00, 0x01, // Sequential ID
		byte(FlagEntryTemporary), // Flags
		byte(BoolTrue), // Value
	}
	result, err := DecodeDataEntryAssignment(bytes.NewBuffer(messageBytes))
	if err != nil {
		t.Fatalf("Unexpected error! %s", err)
	}
	var expected = &MessageDataEntryAssignment{
		Entry: &Entry{
			Name: &ValueString{
				Value: "entry",
				RawValue:  []byte{0x05, 0x65,0x6e,0x74,0x72,0x79},
			},
			Type: TypeBoolean,
			ID: [2]byte{0x50, 0x21},
			Sequence: [2]byte{0x00, 0x01},
			Flags: FlagEntryTemporary,
			Value: &ValueBoolean{
				Value: true,
				RawValue: []byte{byte(BoolTrue)},
			},
		},
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Expected %s but got %s", expected, result)
	}
}

func TestDecodeDataEntryUpdate(t *testing.T) {
	messageBytes := []byte{
		0x50, 0x21, // ID
		0x00, 0x01, // Sequential ID
		byte(TypeBoolean), // Entry Type
		byte(BoolFalse), // Value
	}
	result, err := DecodeDataEntryUpdate(bytes.NewBuffer(messageBytes))
	if err != nil {
		t.Fatalf("Unexpected error! %s", err)
	}
	var expected = &MessageDataEntryUpdate{
		Entry: &Entry{
			ID: [2]byte{0x50, 0x21},
			Sequence: [2]byte{0x00, 0x01},
			Type: TypeBoolean,
			Value: &ValueBoolean{
				Value: false,
				RawValue: []byte{byte(BoolFalse)},
			},
		},
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Expected %s but got %s", expected, result)
	}
}

func TestDecodeDataEntryFlagsUpdate(t *testing.T) {
	messageBytes := []byte{
		0x50, 0x21,
		byte(FlagEntryPersistent),
	}
	result, err := DecodeDataEntryFlagsUpdate(bytes.NewBuffer(messageBytes))
	if err != nil {
		t.Fatalf("Unexpected error! %s", err)
	}
	var expected = &MessageDataEntryFlagsUpdate{
		Entry: &Entry{
			ID: [2]byte{0x50, 0x21},
			Flags: FlagEntryPersistent,
		},
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Expected %s but got %s", expected, result)
	}
}

func TestDecodeDataClearAll(t *testing.T) {
	messageBytes := [4]byte{
		0x00, 0x00, 0x00, 0x00,
	}
	result, err := DecodeDataClearAll(bytes.NewBuffer(messageBytes[:]))
	if err != nil {
		t.Fatalf("Unexpected error! %s", err)
	}
	var expected = &MessageDataClearAll{
		PotentialMagic: messageBytes,
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Expected %s but got %s", expected, result)
	}
}