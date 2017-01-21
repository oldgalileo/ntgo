package ntgo

import (
	"errors"
	"github.com/imdario/mergo"
)

const (
	TypeBoolean EntryType = 0x00
	TypeDouble = 0x01
	TypeString = 0x02
	TypeRawData  = 0x03
	TypeBooleanArr = 0x10
	TypeDoubleArr = 0x11
	TypeStringArr = 0x12
	TypeRPCDef = 0x20

	FlagTemporary EntryFlag = 0x00
	FlagPersistent = 0x01
	FlagReserved = 0xFE

	BoolFalse byte = 0x00
	BoolTrue = 0x01
)

var (
	ErrEntryDataInvalid = errors.New("entry: data invalid")
)

type EntryType byte

type Entry struct {
	Name string
	Type EntryType
	Value EntryValue
}

type EntryFlag byte

type EntryValue interface {
	GetRaw() []byte
	GetValue() interface{}
	UpdateRaw([]byte) error
	UpdateValue(interface{}) error
}

type EntryBoolean struct {
	Value bool
	rawValue byte
}

func DecodeBoolean(data []byte) (*EntryBoolean, error) {
	if len(data) > 1 || len(data) < 1 {
		return nil, ErrEntryDataInvalid
	}
	entry := &EntryBoolean{rawValue: data[0]}
	if entry.rawValue == BoolFalse {
		entry.Value = false
		return entry, nil
	} else if entry.rawValue == BoolTrue {
		entry.Value = true
		return entry, nil
	} else {
		return nil, ErrEntryDataInvalid
	}
}

func BuildBoolean(value bool) *EntryBoolean {
	var rawValue byte
	if value {
		rawValue = BoolTrue
	} else {
		rawValue = BoolFalse
	}
	return &EntryBoolean{
		Value: value,
		rawValue: rawValue,
	}
}

func (entry *EntryBoolean) UpdateRaw(raw []byte) error {
	newEntry, entryErr := DecodeBoolean(raw)
	if entryErr != nil {
		return entryErr
	}
	return mergo.MergeWithOverwrite(entry, newEntry)
}

func (entry *EntryBoolean) UpdateValue(value bool) error {
	newEntry := BuildBoolean(value)
	return mergo.MergeWithOverwrite(entry, newEntry)
}

func (entry *EntryBoolean) GetRaw() []byte {
	return []byte{entry.rawValue}
}

func (entry *EntryBoolean) GetValue() bool {
	return entry.Value
}