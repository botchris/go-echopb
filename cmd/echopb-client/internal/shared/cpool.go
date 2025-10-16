package shared

import (
	"sync"

	"google.golang.org/grpc"
)

// ConnectionPool is a simple round-robin connection pool for gRPC client connections.
type ConnectionPool struct {
	conn []*grpc.ClientConn

	count uint64
	mu    sync.Mutex
}

// NewConnectionPool creates a new ConnectionPool with the given gRPC client connections.
func NewConnectionPool(cc []*grpc.ClientConn) *ConnectionPool {
	return &ConnectionPool{conn: cc}
}

// Next returns the next gRPC client connection in a round-robin fashion.
func (cp *ConnectionPool) Next() *grpc.ClientConn {
	cp.mu.Lock()

	key := cp.count % uint64(len(cp.conn))
	cp.count++

	cp.mu.Unlock()

	return cp.conn[key]
}
