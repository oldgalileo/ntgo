package ntgo

import "errors"

const (
	DefaultAddress string = "0.0.0.0"
	DefaultPort string = "1735"

	ModeClient mode = 0
	ModeServer = 1
)

var (
	ErrUnknownMode error = errors.New("ntgo: unknown or unsupported network table mode. must be client or server")
)

type mode int

type NetworkTables struct {
	Address string
	Port string
	Mode mode
	Operator
}

type Operator interface {
	CreateEntry(message string) error
	DeleteEntry(message string) error
	UpdateEntry(message string) error
	GetEntry(message string)    error
	Initialize(nt NetworkTables) error
}

func DefaultSettings() *NetworkTables {
	return &NetworkTables{
		Address: DefaultAddress,
		Port: DefaultPort,
		Mode: ModeClient,
	}
}

func (nt *NetworkTables) Initialize() error {
	var operator Operator
	if nt.Mode == ModeClient {
		operator = &client{}
	} else if nt.Mode == ModeServer {
		// TODO: Add server
	} else {
		return ErrUnknownMode
	}
	initErr := operator.Initialize(*nt)
	if initErr != nil {
		return initErr
	}
	nt.Operator = operator
	return nil
}