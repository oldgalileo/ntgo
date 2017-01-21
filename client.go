package ntgo

import "net"

type client struct {
	conn *net.Conn
}

func (cl *client) Initialize(nt NetworkTables) error {
	conn, connErr := net.Dial("tcp", net.JoinHostPort(nt.Address, nt.Port))
	if connErr != nil {
		return connErr
	}
	cl.conn = &conn
	return nil
}

func (cl *client) CreateEntry(message string) error

func (cl *client) UpdateEntry(message string) error

func (cl *client) DeleteEntry(message string) error

func (cl *client) GetEntry(message string) error