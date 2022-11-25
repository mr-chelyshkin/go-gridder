package fasthttp

import (
	"fmt"
	"net"
	"sync/atomic"
	"time"
)

type gracefulListener struct {
	ln net.Listener

	maxWaitTime time.Duration
	done        chan struct{}
	connsCount  uint64
	shutdown    uint64
}

func newGracefulListener(ln net.Listener, maxWaitTime time.Duration) net.Listener {
	return &gracefulListener{
		ln:          ln,
		maxWaitTime: maxWaitTime,
		done:        make(chan struct{}),
	}
}

// Accept waits for and returns the next connection to the listener.
func (ln *gracefulListener) Accept() (net.Conn, error) {
	c, err := ln.ln.Accept()
	if err != nil {
		return nil, err
	}
	atomic.AddUint64(&ln.connsCount, 1)

	return &gracefulConn{
		Conn: c,
		ln:   ln,
	}, nil
}

// Addr returns the listener's network address.
func (ln *gracefulListener) Addr() net.Addr {
	return ln.ln.Addr()
}

// Close closes the inner listener and waits until all the pending open connections are closed.
func (ln *gracefulListener) Close() error {
	err := ln.ln.Close()
	if err != nil {
		return nil
	}
	return ln.waitForZeroConns()
}

func (ln *gracefulListener) waitForZeroConns() error {
	atomic.AddUint64(&ln.shutdown, 1)
	if atomic.LoadUint64(&ln.connsCount) == 0 {
		close(ln.done)
		return nil
	}

	select {
	case <-ln.done:
		return nil
	case <-time.After(ln.maxWaitTime):
		return fmt.Errorf("cannot complete graceful shutdown in %s", ln.maxWaitTime)
	}
}

func (ln *gracefulListener) closeConn() {
	connsCount := atomic.AddUint64(&ln.connsCount, ^uint64(0))

	if atomic.LoadUint64(&ln.shutdown) != 0 && connsCount == 0 {
		close(ln.done)
	}
}

type gracefulConn struct {
	net.Conn
	ln *gracefulListener
}

func (c *gracefulConn) Close() error {
	err := c.Conn.Close()
	if err != nil {
		return err
	}

	c.ln.closeConn()
	return nil
}
