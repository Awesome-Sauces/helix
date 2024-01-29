package hydrogen

import (
	"log"
	"net"
)

const (
	API_PORT = "7771"
	N2N_PORT = "7770"

	CLOSE_EVENT = 1
	ERROR_EVENT = 2
)

type EndpointReceiver interface {
	Call(conn *Connection)
	Response() []byte
}

type Server struct {
	address   string
	events    chan int
	listener  net.Listener
	endpoints map[string]EndpointReceiver
}

func NewServer(address string) *Server {
	return &Server{
		address:   address,
		endpoints: make(map[string]EndpointReceiver),
	}
}

func (server *Server) a_Listen() error {
	listener, err := net.Listen("tcp", server.address)
	server.listener = listener

	server.events = make(chan int)

	return err
}

func (server *Server) Listen() error {
	err := server.a_Listen()

	if err != nil {
		return err
	}

	defer server.listener.Close()

	event := 0

	event_listener := func() {
		event = <-server.events
	}

	go event_listener()

	for event == 0 {

		conn, err := server.listener.Accept()
		if err != nil {
			log.Println(err.Error())
		}

		endpoint, args := DigestRequest((conn))
		variables := MapVariables(args)

		connect := Connection{
			args: variables,
		}

		server.endpoints[endpoint].Call(&connect)
		conn.Write(server.endpoints[endpoint].Response())
	}

	return err
}

func (server *Server) Close() {
	server.events <- CLOSE_EVENT
}