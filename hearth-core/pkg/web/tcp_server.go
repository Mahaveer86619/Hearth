package web

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/Mahaveer86619/Hearth/pkg/logger"
	"github.com/Mahaveer86619/Hearth/pkg/services/ingestion_pipeline"
)

type TCPServer struct {
	addr     string
	listener net.Listener
	quit     chan struct{}
	wg       sync.WaitGroup
	pipeline *ingestion_pipeline.IngestionPipelineService
}

func NewTCPServer(addr string) *TCPServer {
	return &TCPServer{
		addr:     addr,
		quit:     make(chan struct{}),
		pipeline: ingestion_pipeline.GetInstance(),
	}
}

func (s *TCPServer) Start() error {
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", s.addr, err)
	}
	s.listener = listener
	// defer s.listener.Close() // Close is handled in Stop

	// logger.Info("TCP", "Ingestion Server listening on %s", s.addr)

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
				logger.Error("TCP", "Accept error: %v", err)
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
	// logger.Info("TCP", "New connection from: %s", conn.RemoteAddr())

	scanner := bufio.NewScanner(conn)

	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)

	for scanner.Scan() {
		text := scanner.Bytes()
		data := make([]byte, len(text))
		copy(data, text)

		s.pipeline.Ingest(data)
	}

	if err := scanner.Err(); err != nil {
		// Only log if not closing
		select {
		case <-s.quit:
			return
		default:
			if !strings.Contains(err.Error(), "use of closed network connection") {
				logger.Error("TCP", "connection error from %s: %v", conn.RemoteAddr(), err)
			}
		}
	}
}
