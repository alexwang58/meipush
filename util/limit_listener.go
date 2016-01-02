package util

import (
	"fmt"
	"net"
	"sync"
)

type LimitListener struct {
	*net.TCPListener
	sem chan struct{}
}

func NewLimitListener(l *net.TCPListener, n int) *LimitListener {
	return &LimitListener{l, make(chan struct{}, n)}
}

func (l *LimitListener) acquire() { l.sem <- struct{}{} }
func (l *LimitListener) release() { <-l.sem }

func (l *LimitListener) Release() {
	l.release()
}

func (l *LimitListener) Accept() (*limitListenerConn, error) {
	l.acquire()
	fmt.Println("---AcceptTCP")
	c, err := l.AcceptTCP()
	if err != nil {
		l.release()
		return nil, err
	}
	return &limitListenerConn{TCPConn: c, release: l.release}, nil
}

type limitListenerConn struct {
	*net.TCPConn
	releaseOnce sync.Once
	release     func()
}

func (l *limitListenerConn) Close() error {
	err := l.TCPConn.Close()
	l.releaseOnce.Do(l.release)
	return err
}
