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