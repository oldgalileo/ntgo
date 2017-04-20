package ntgo

import (
	"errors"
	"github.com/imdario/mergo"
	"io"
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

//
//
//
//
//
//
//
//
//
//


var (
	ErrEntryDataInvalid = errors.New("entry: data invalid")
	ErrArrayIndexOutOfBounds = errors.New("array: index out of bounds")
	ErrArrayOutOfSpace = errors.New("array: no more space")
)

type EntryType byte

type Entry struct {
	Name ValueString
	Type EntryType
	Value EntryValue
}

type EntryFlag byte

type EntryValue interface {
	GetRaw() []byte
	UpdateRaw(io.Reader) error
}

func validateEntryValue() {
	var _ EntryValue = BuildBoolean(true).Value
	var _ EntryValue = BuildString("meme").Value
	var _ EntryValue = BuildRaw([]byte("kek")).Value
}

type EntryValueArray interface {
	Get(int) EntryValue
}

type ValueBoolean struct {
	Value    bool
	RawValue []byte
}

func DecodeBoolean(r io.Reader) (*ValueBoolean, error) {
	val := make([]byte, 1)
	_, readErr := io.ReadFull(r, val)
	if readErr != nil {
		return nil, readErr
	}
	entry := &ValueBoolean{RawValue: val}
	if entry.RawValue[0] == BoolFalse {
		entry.Value = false
		return entry, nil
	} else if entry.RawValue[0] == BoolTrue {
		entry.Value = true
		return entry, nil
	} else {
		return nil, ErrEntryDataInvalid
	}
}

func BuildBoolean(value bool) *ValueBoolean {
	var rawValue []byte
	if value {
		rawValue = []byte{BoolTrue}
	} else {
		rawValue = []byte{BoolFalse}
	}
	return &ValueBoolean{
		Value: value,
		RawValue: rawValue,
	}
}

func (entry *ValueBoolean) UpdateRaw(r io.Reader) error {
	newEntry, entryErr := DecodeBoolean(r)
	if entryErr != nil {
		return entryErr
	}
	return mergo.MergeWithOverwrite(entry, *newEntry)
}

func (entry *ValueBoolean) UpdateValue(value bool) error {
	return mergo.MergeWithOverwrite(entry, *BuildBoolean(value))
}

type ValueString struct {
	Value    string
	RawValue []byte
}

func DecodeString(r io.Reader) (*ValueString, error) {
	uleb, ulebData := DecodeAndSaveULEB128(r)
	data := make([]byte, uleb)
	_, readErr := io.ReadFull(r, data)
	if readErr != nil {
		return nil, readErr
	}
	return &ValueString{
		Value: string(data),
		RawValue: append(ulebData, data...),
	}, nil
}

func BuildString(value string) *ValueString {
	stringBytes := []byte(value)
	rawValue := append(EncodeULEB128(uint32(len(stringBytes))), stringBytes...)
	return &ValueString{
		Value: value,
		RawValue: rawValue,
	}
}

func (entry *ValueString) UpdateRaw(r io.Reader) error {
	newEntry, newErr := DecodeString(r)
	if newErr != nil {
		return newErr
	}
	return mergo.MergeWithOverwrite(entry, *newEntry)
}

func (entry *ValueString) UpdateValue(value string) error {
	return mergo.MergeWithOverwrite(entry, *BuildString(value))
}

type ValueDouble struct {
	Value    float64
	RawValue []byte
}

func DecodeDouble(r io.Reader) (*ValueDouble, error) {
	data := make([]byte, 8)
	_, readErr := io.ReadFull(r, data)
	if readErr != nil {
		return nil, readErr
	}
	return &ValueDouble{
		Value: BytesToFloat64(data),
		RawValue: data,
	}, nil
}

func BuildDouble(value float64) *ValueDouble {
	return &ValueDouble{
		Value: value,
		RawValue: Float64ToBytes(value),
	}
}

func (entry *ValueDouble) UpdateRaw(r io.Reader) error {
	newEntry, newErr := DecodeDouble(r)
	if newErr != nil {
		return newErr
	}
	return mergo.MergeWithOverwrite(entry, *newEntry)
}

func (entry *ValueDouble) UpdateValue(value float64) error {
	return mergo.MergeWithOverwrite(entry, *BuildDouble(value))
}

type ValueRaw struct {
	Value []byte
	RawValue []byte
}

func DecodeRaw(r io.Reader) (*ValueRaw, error) {
	uleb, ulebData := DecodeAndSaveULEB128(r)
	data := make([]byte, uleb)
	_, readErr := io.ReadFull(r, data)
	if readErr != nil {
		return nil, readErr
	}
	return &ValueRaw{
		Value: data,
		RawValue: append(ulebData, data...),
	}, nil
}

func BuildRaw(value []byte) *ValueRaw {
	return &ValueRaw{
		Value: value,
		RawValue: append(EncodeULEB128(uint32(len(value))), value...),
	}
}


func (entry *ValueRaw) UpdateRaw(r io.Reader) error {
	newEntry, newErr := DecodeRaw(r)
	if newErr != nil {
		return newErr
	}
	return mergo.MergeWithOverwrite(entry, *newEntry)
}

func (entry *ValueRaw) UpdateValue(value []byte) error {
	return mergo.MergeWithOverwrite(entry, *BuildRaw(value))
}

type ValueBooleanArray struct {
	index uint8
	elements []*ValueBoolean
}

func DecodeBooleanArray(r io.Reader) (*ValueBooleanArray, error) {
	indexData := make([]byte, 1)
	_, readErr := io.ReadFull(r, indexData)
	if readErr != nil {
		return nil, readErr
	}
	index := uint8(indexData[0])
	elements := make([]*ValueBoolean, index)
	var i uint8 = 0
	for ; i < index; i++ {
		boolean, decodeErr := DecodeBoolean(r)
		if decodeErr != nil {
			return nil, decodeErr
		}
		elements = append(elements, boolean)
	}
	return &ValueBooleanArray{
		index: index,
		elements: elements,
	}, nil
}

func BuildBooleanArray(values []*ValueBoolean) *ValueBooleanArray {
	var index uint8
	if len(values) > 255 {
		index = 255
	} else {
		index = uint8(len(values))
	}
	return &ValueBooleanArray{
		index: index,
		elements: values[:index],
	}
}

func (array *ValueBooleanArray) Get(index uint8) (*ValueBoolean, error) {
	if index > array.index {
		return nil, ErrArrayIndexOutOfBounds
	}
	return array.elements[index], nil
}

func (array *ValueBooleanArray) Update(index uint8, boolean ValueBoolean) error {
	if index > array.index {
		return ErrArrayIndexOutOfBounds
	}
	return mergo.MergeWithOverwrite(array.elements[index], boolean)
}

func (array *ValueBooleanArray) Add(boolean *ValueBoolean) error {
	if array.index == 255 {
		return ErrArrayOutOfSpace
	}
	array.elements = append(array.elements, boolean)
	return nil
}

func (array *ValueBooleanArray) ToBytes() []byte {
	data := []byte(array.index)
	var i uint8 = 0
	for ; i < array.index; i++ {
		data = append(data, array.elements[i].RawValue...)
	}
	return data
}

func DecodeDoubleArray(r io.Reader) (*ValueDoubleArray, error) {
	indexData := make([]byte, 1)
	_, readErr := io.ReadFull(r, indexData)
	if readErr != nil {
		return nil, readErr
	}
	index := uint8(indexData[0])
	elements := make([]*ValueDouble, index)
	var i uint8 = 0
	for ; i < index; i++ {
		double, decodeErr := DecodeDouble(r)
		if decodeErr != nil {
			return nil, decodeErr
		}
		elements = append(elements, double)
	}
	return &ValueDoubleArray{
		index: index,
		elements: elements,
	}, nil
}