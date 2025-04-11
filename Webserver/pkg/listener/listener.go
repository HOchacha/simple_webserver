package listener

import (
	"fmt"
	"log"
	"net"
	"webserver/Webserver/pkg/types"
)

type ListenerFilter struct {
	lSocket net.Listener
	port    string
	next    types.Filter
}

func (f *ListenerFilter) Init(config map[string]interface{}) error {
	f.port = config["port"].(string)
	return nil
}

func (f *ListenerFilter) Handle(_ net.Conn) error {
	listener, err := net.Listen("tcp", f.port)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}
	defer listener.Close()
	log.Printf("Listening on %s", f.port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}
		go func(c net.Conn) {
			defer c.Close()
			if f.next != nil {
				if err := f.next.Handle(c); err != nil {
					log.Printf("Error in next filter: %v", err)
				}
			}
		}(conn)
	}
}

func (f *ListenerFilter) SetNext(next types.Filter) {
	f.next = next
}
