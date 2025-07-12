package testdata

import (
	"context"
	"errors"
	"io"
	"net"
	"time"
)

type TCPServer struct {
	Listener net.Listener
	Handle   func(context.Context, []byte) []byte
}

func (server *TCPServer) Start(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				conn, err := server.Listener.Accept()
				if err != nil {
					if errors.Is(err, net.ErrClosed) {
						return
					}

					panic(err)
				}

				go server.serve(ctx, conn)
			}
		}
	}()
}

func (server *TCPServer) serve(ctx context.Context, conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, 1024)

	for {
		length := server.read(ctx, conn, buf)
		if length == 0 {
			return
		}

		res := server.Handle(ctx, buf[:length])
		if res == nil {
			continue
		}

		server.write(ctx, conn, res)
	}
}

func (server *TCPServer) read(ctx context.Context, conn net.Conn, buf []byte) int {
	var netErr net.Error

	for {
		select {
		case <-ctx.Done():
			return 0
		default:
			conn.SetReadDeadline(time.Now().Add(time.Millisecond))

			length, err := conn.Read(buf)
			if err == nil {
				return length
			}

			if errors.As(err, &netErr) && netErr.Timeout() {
				continue
			}

			if errors.Is(err, net.ErrClosed) {
				return 0
			}

			if errors.Is(err, io.EOF) {
				return 0
			}

			panic(err)
		}
	}
}

func (server *TCPServer) write(ctx context.Context, conn net.Conn, data []byte) {
	var netErr net.Error

	for {
		select {
		case <-ctx.Done():
			return
		default:
			conn.SetWriteDeadline(time.Now().Add(time.Millisecond))

			_, err := conn.Write(data)
			if err == nil {
				return
			}

			if errors.As(err, &netErr) && netErr.Timeout() {
				continue
			}

			if errors.Is(err, net.ErrClosed) {
				return
			}

			panic(err)
		}
	}
}
