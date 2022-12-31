package posixmq

import (
	"golang.org/x/sys/unix"
	"sync"
	"syscall"
)

type epoll struct {
	fd          int
	lock        *sync.RWMutex
	connections map[int]int
}

func MkEpoll() (*epoll, error) {
	fd, err := unix.EpollCreate1(0)
	if err != nil {
		return nil, err
	}
	return &epoll{
		fd:   fd,
		lock: &sync.RWMutex{},
	}, nil
}

func (e *epoll) Add(fd int) error {
	// Extract file descriptor associated with the connection
	err := unix.EpollCtl(e.fd, syscall.EPOLL_CTL_ADD, fd, &unix.EpollEvent{Events: unix.POLLIN | unix.POLLHUP, Fd: int32(fd)})
	if err != nil {
		return err
	}

	return nil
}

func (e *epoll) Remove(fd int) error {

	err := unix.EpollCtl(e.fd, syscall.EPOLL_CTL_DEL, fd, nil)
	if err != nil {
		return err
	}

	return nil
}

func (e *epoll) Wait() (connections []int, err error) {
	events := make([]unix.EpollEvent, 100)
	n, err := unix.EpollWait(e.fd, events, 100)
	if err != nil {
		return nil, err
	}
	e.lock.RLock()
	defer e.lock.RUnlock()
	connections = make([]int, 0)
	for i := 0; i < n; i++ {
		con := e.connections[int(events[i].Fd)]
		connections = append(connections, con)
	}
	return connections, nil
}
