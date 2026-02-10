package web

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

type TCPServer struct {
	addr     string
	listener net.Listener
	quit     chan struct{}
	wg       sync.WaitGroup
}

func NewTCPServer(addr string) *TCPServer {
	return &TCPServer{
		addr: addr,
		quit: make(chan struct{}),
	}
}

func (s *TCPServer) Start() error {
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", s.addr, err)
	}
	s.listener = listener
	// defer s.listener.Close() // Close is handled in Stop

	log.Printf("TCP Ingestion Server listening on %s", s.addr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			select {
			case <-s.quit:
				return nil
			default:
				// Check if error is due to listener closed
				if strings.Contains(err.Error(), "use of closed network connection") {
					return nil
				}
				log.Printf("TCP Accept error: %v", err)
				continue
			}
		}
		s.wg.Add(1)
		go s.handleConnection(conn)
	}
}

func (s *TCPServer) Stop(ctx context.Context) error {
	close(s.quit)
	if s.listener != nil {
		s.listener.Close()
	}
	
	done := make(chan struct{})
	go func() {
		s.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(5 * time.Second): // Safety timeout if ctx is not set or long
		return fmt.Errorf("tcp server shutdown timed out")
	}
}

func (s *TCPServer) Addr() string {
	return s.addr
}

func (s *TCPServer) Name() string {
	return "TCP Ingestion Server"
}

func (s *TCPServer) handleConnection(conn net.Conn) {
	defer s.wg.Done()
	defer conn.Close()
	// log.Printf("New TCP connection from: %s", conn.RemoteAddr())

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		line := scanner.Text()
		// TODO: Push to WAL and WebSocket Hub
		_ = line 
	}

	if err := scanner.Err(); err != nil {
		// Only log if not closing
		select {
		case <-s.quit:
			return
		default:
			if !strings.Contains(err.Error(), "use of closed network connection") {
				log.Printf("TCP connection error from %s: %v", conn.RemoteAddr(), err)
			}
		}
	}
}
