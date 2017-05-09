package ntgo

import (
	"io"
	"errors"
)

const (
	MessageTypeKeepAlive                MessageType = 0x00
	MessageTypeClientHello                          = 0x01
	MessageTypeProtocVersionUnsupported             = 0x02
	MessageTypeServerHelloComplete                  = 0x03
	MessageTypeServerHello                          = 0x04
	MessageTypeClientHelloComplete                  = 0x05
	MessageTypeEntryAssignment                      = 0x10
	MessageTypeEntryUpdate                          = 0x11
	MessageTypeEntryFlagsUpdate                     = 0x12
	MessageTypeEntryDelete                          = 0x13
	MessageTypeClearAll                             = 0x14
	MessageTypeRPCExecute                           = 0x20
	MessageTypeRPCResponse                          = 0x21
	MessageTypeUndef                                = 0xFF
)

const (
	FlagMessageClientNew      MessageFlag = 0x00
	FlagMessageClientSeen                 = 0x01
	FlagMessageClientReserved             = 0xFF
)

var (
	ErrMessageNoSuchType     = errors.New("message: no such message type")
	ErrMessageFlagNoSuchType = errors.New("messageflag: no such flag")
)

var (
	// DangerousMagic is a 4 byte sequence specified by the protocol documentation
	// that should be sent as the message payload for a ClearAll message.
	DangerousMagic = [4]byte{0xD0, 0x6C, 0xB2, 0x7A}
)

type ProtocolRevision [2]byte

type Message struct {
	Type MessageType
	Data MessageData
}

func DecodeMessage(r io.Reader) (*Message, error) {
	messageType, typeErr := DecodeMessageType(r)
	if typeErr != nil {
		return nil, typeErr
	}
	message := &Message{
		Type: messageType,
	}
	var messageData MessageData
	var dataErr error = nil
	switch messageType {
	case MessageTypeKeepAlive:
		messageData = &MessageDataKeepAlive{}
	case MessageTypeClientHello:
		messageData, dataErr = DecodeDataClientHello(r)
	case MessageTypeProtocVersionUnsupported:
		messageData, dataErr = DecodeDataProtocVersionUnsupported(r)
	case MessageTypeServerHelloComplete:
		messageData = &MessageDataServerHelloComplete{}
	case MessageTypeServerHello:
		messageData, dataErr = DecodeDataServerHello(r)
	case MessageTypeClientHelloComplete:
		messageData = &MessageDataClientHelloComplete{}
	case MessageTypeEntryAssignment:
		messageData, dataErr = DecodeDataEntryAssignment(r)
	case MessageTypeEntryUpdate:
		messageData, dataErr = DecodeDataEntryUpdate(r)
	case MessageTypeEntryFlagsUpdate:
		messageData, dataErr = DecodeDataEntryFlagsUpdate(r)
	case MessageTypeClearAll:
		messageData, dataErr = DecodeDataClearAll(r)
	default:
		dataErr = ErrMessageNoSuchType
	}
	message.Data = messageData
	return message, dataErr
}

type MessageFlag byte

func DecodeMessageFlag(r io.Reader) (MessageFlag, error) {
	flagRaw := make([]byte, 1)
	_, flagErr := r.Read(flagRaw)
	if flagErr != nil {
		return FlagMessageClientReserved, flagErr
	}
	switch MessageFlag(flagRaw[0]) {
	case FlagMessageClientNew:
		return FlagMessageClientNew, nil
	case FlagMessageClientSeen:
		return FlagMessageClientSeen, nil
	default:
		return FlagMessageClientReserved, ErrMessageFlagNoSuchType
	}
}

type MessageType byte

func DecodeMessageType(r io.Reader) (MessageType, error) {
	typeRaw := make([]byte, 1)
	_, typeErr := r.Read(typeRaw)
	if typeErr != nil {
		return EntryTypeUndef, typeErr
	}
	return MessageType(typeRaw[0]), nil
}

type MessageData interface{}

type MessageDataKeepAlive struct{}

type MessageDataClientHello struct {
	ProtocVersion ProtocolRevision
	Identity      *ValueString
}

func DecodeDataClientHello(r io.Reader) (*MessageDataClientHello, error) {
	protocRaw := [2]byte{}
	_, protocErr := r.Read(protocRaw[:])
	if protocErr != nil {
		return nil, protocErr
	}
	identity, identityErr := DecodeString(r)
	if identityErr != nil {
		return nil, identityErr
	}
	return &MessageDataClientHello{
		ProtocVersion: ProtocolRevision(protocRaw),
		Identity:      identity,
	}, nil
}

type MessageDataProtocVersionUnsupported struct {
	SupportedProtoc ProtocolRevision
}

func DecodeDataProtocVersionUnsupported(r io.Reader) (*MessageDataProtocVersionUnsupported, error) {
	protocRaw := [2]byte{}
	_, protocErr := r.Read(protocRaw[:])
	if protocErr != nil {
		return nil, protocErr
	}
	return &MessageDataProtocVersionUnsupported{
		SupportedProtoc: protocRaw,
	}, nil
}

type MessageDataServerHelloComplete struct{}

type MessageDataServerHello struct {
	Flags    MessageFlag
	Identity *ValueString
}

func DecodeDataServerHello(r io.Reader) (*MessageDataServerHello, error) {
	flag, flagErr := DecodeMessageFlag(r)
	if flagErr != nil {
		return nil, flagErr
	}
	identity, identityErr := DecodeString(r)
	if identityErr != nil {
		return nil, identityErr
	}
	return &MessageDataServerHello{
		Flags:    flag,
		Identity: identity,
	}, nil
}

type MessageDataClientHelloComplete struct{}

type MessageDataEntryAssignment struct {
	Entry *Entry
}

func DecodeDataEntryAssignment(r io.Reader) (*MessageDataEntryAssignment, error) {
	name, nameErr := DecodeString(r)
	if nameErr != nil {
		return nil, nameErr
	}
	entryType, entryErr := DecodeEntryType(r)
	if entryErr != nil {
		return nil, entryErr
	}
	idRaw := [2]byte{}
	_, idErr := r.Read(idRaw[:])
	if idErr != nil {
		return nil, idErr
	}
	seqRaw := [2]byte{}
	_, seqErr := r.Read(seqRaw[:])
	if seqErr != nil {
		return nil, idErr
	}
	flag, flagErr := DecodeEntryFlag(r)
	if flagErr != nil {
		return nil, flagErr
	}
	value, valueErr := DecodeEntryValue(r, entryType)
	if valueErr != nil {
		return nil, valueErr
	}
	return &MessageDataEntryAssignment{
		Entry: &Entry{
			Name:     name,
			Type:     entryType,
			ID:       idRaw,
			Sequence: seqRaw,
			Flags:    flag,
			Value:    value,
		},
	}, nil
}

type MessageDataEntryUpdate struct {
	Entry *Entry
}

func DecodeDataEntryUpdate(r io.Reader) (*MessageDataEntryUpdate, error) {
	idRaw := [2]byte{}
	_, idErr := r.Read(idRaw[:])
	if idErr != nil {
		return nil, idErr
	}
	seqRaw := [2]byte{}
	_, seqErr := r.Read(seqRaw[:])
	if seqErr != nil {
		return nil, idErr
	}
	value, entryType, valueErr := DecodeEntryValueAndType(r)
	if valueErr != nil {
		return nil, valueErr
	}
	return &MessageDataEntryUpdate{
		Entry: &Entry{
			ID:       idRaw,
			Sequence: seqRaw,
			Type:     entryType,
			Value:    value,
		},
	}, nil
}

type MessageDataEntryFlagsUpdate struct {
	Entry *Entry
}

func DecodeDataEntryFlagsUpdate(r io.Reader) (*MessageDataEntryFlagsUpdate, error) {
	idRaw := [2]byte{}
	_, idErr := r.Read(idRaw[:])
	if idErr != nil {
		return nil, idErr
	}
	flag, flagErr := DecodeEntryFlag(r)
	if flagErr != nil {
		return nil, flagErr
	}
	return &MessageDataEntryFlagsUpdate{
		Entry: &Entry{
			ID:    idRaw,
			Flags: flag,
		},
	}, nil
}

type MessageDataEntryDelete struct {
	Entry *Entry
}

func DecodeDataEntryDelete(r io.Reader) (*MessageDataEntryDelete, error) {
	idRaw := [2]byte{}
	_, idErr := r.Read(idRaw[:])
	if idErr != nil {
		return nil, idErr
	}
	return &MessageDataEntryDelete{
		Entry: &Entry{
			ID: idRaw,
		},
	}, nil
}

type MessageDataClearAll struct {
	PotentialMagic [4]byte
}

func DecodeDataClearAll(r io.Reader) (*MessageDataClearAll, error) {
	magicRaw := [4]byte{}
	_, magicErr := r.Read(magicRaw[:])
	if magicErr != nil {
		return nil, magicErr
	}
	return &MessageDataClearAll{
		PotentialMagic: magicRaw,
	}, nil
}

type MessageDataRPCExecute struct {
	EntryID     [2]byte
	UniqueID    [2]byte
	ParamLength uint32
	Params      []byte
}

func DecodeDataRPCExecute(r io.Reader) (*MessageDataRPCExecute, error) {
	entryIDRaw := [2]byte{}
	_, entryIDErr := r.Read(entryIDRaw[:])
	if entryIDErr != nil {
		return nil, entryIDErr
	}
	uniqueIDRaw := [2]byte{}
	_, uniqueIDErr := r.Read(uniqueIDRaw[:])
	if uniqueIDErr != nil {
		return nil, uniqueIDErr
	}
	paramsSize, ulebErr := DecodeULEB128(r)
	if ulebErr != nil {
		return nil, ulebErr
	}
	params := make([]byte, paramsSize)
	_, paramsErr := r.Read(params)
	if paramsErr != nil {
		return nil, paramsErr
	}
	return &MessageDataRPCExecute{
		EntryID:     entryIDRaw,
		UniqueID:    uniqueIDRaw,
		ParamLength: paramsSize,
		Params:      params,
	}, nil
}

type MessageDataRPCResponse struct {
	EntryID      [2]byte
	UniqueID     [2]byte
	ResultLength uint32
	Results      []byte
}

func DecodeDataRPCReseponse(r io.Reader) (*MessageDataRPCResponse, error) {
	entryIDRaw := [2]byte{}
	_, entryIDErr := r.Read(entryIDRaw[:])
	if entryIDErr != nil {
		return nil, entryIDErr
	}
	uniqueIDRaw := [2]byte{}
	_, uniqueIDErr := r.Read(uniqueIDRaw[:])
	if uniqueIDErr != nil {
		return nil, uniqueIDErr
	}
	resultsSize, ulebErr := DecodeULEB128(r)
	if ulebErr != nil {
		return nil, ulebErr
	}
	results := make([]byte, resultsSize)
	_, resultsErr := r.Read(results)
	if resultsErr != nil {
		return nil, resultsErr
	}
	return &MessageDataRPCResponse{
		EntryID:      entryIDRaw,
		UniqueID:     uniqueIDRaw,
		ResultLength: resultsSize,
		Results:      results,
	}, nil
}
