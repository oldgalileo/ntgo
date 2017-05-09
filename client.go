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

func (cl *client) CreateEntry(entry Entry) error { return nil }

func (cl *client) UpdateEntry(entry Entry) error { return nil }

func (cl *client) DeleteEntry(entry Entry) error { return nil }

func (cl *client) GetEntry(id [2]byte) error { return nil }