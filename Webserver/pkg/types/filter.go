package types

import "net"

type Filter interface {
	Init(config map[string]interface{}) error
	Handle(conn net.Conn) error
	SetNext(next Filter)
}
