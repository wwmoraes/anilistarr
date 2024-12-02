package testdata

import (
	"context"
	"net"
)

// Command represents command roundtrip. It contains a string representation of
// a request and its response.
type Command struct {
	Request  string
	Response string
}

// CommandServer is a mock TCP server that matches known commands to responses.
// It does NOT implement any proper tokenization or parsing; instead it does a
// dumb string match.
//
// Commands is a stack, with top on the first element.
type CommandServer struct {
	TCPServer
	Commands []Command
}

// Handle matches an incoming request data with the current command request
// on top of the stack. Returns the associated command response if it matches,
// nil otherwise.
//
// A match also pops that command from the stack.
func (server *CommandServer) Handle(ctx context.Context, req []byte) []byte {
	if len(server.Commands) == 0 {
		return nil
	}

	if string(req) != server.Commands[0].Request {
		return nil
	}

	var res string

	res, server.Commands = server.Commands[0].Response, server.Commands[1:]

	return []byte(res)
}

func NewCommandServer(listener net.Listener, commands ...Command) CommandServer {
	server := CommandServer{
		TCPServer: TCPServer{
			Listener: listener,
		},
		Commands: commands,
	}

	server.TCPServer.Handle = server.Handle

	return server
}
