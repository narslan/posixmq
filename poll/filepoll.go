// Copyright (c) 2017 Mail.Ru Group All rights reserved.
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

//go:build linux
// +build linux

package poll

import (
	"os"
)

// New creates new epoll-based Poller instance with given config.
func New(c *Config) (Poller, error) {
	cfg := c.withDefaults()

	epoll, err := EpollCreate(&EpollConfig{
		OnWaitError: cfg.OnWaitError,
	})
	if err != nil {
		return nil, err
	}

	return poller{epoll}, nil
}

// poller implements Poller interface.
type poller struct {
	*Epoll
}

// Start implements Poller.Start() method.
func (ep poller) Start(desc *Desc, cb CallbackFn) error {
	err := ep.Add(desc.fd(), toEpollEvent(desc.event),
		func(ep EpollEvent) {
			var event Event

			if ep&EPOLLHUP != 0 {
				event |= EventHup
			}
			if ep&EPOLLRDHUP != 0 {
				event |= EventReadHup
			}
			if ep&EPOLLIN != 0 {
				event |= EventRead
			}
			if ep&EPOLLOUT != 0 {
				event |= EventWrite
			}
			if ep&EPOLLERR != 0 {
				event |= EventErr
			}
			if ep&_EPOLLCLOSED != 0 {
				event |= EventPollerClosed
			}

			cb(event)
		},
	)
	if err == nil {
		if err = setNonblock(desc.fd(), true); err != nil {
			return os.NewSyscallError("setnonblock", err)
		}
	}
	return err
}

// Stop implements Poller.Stop() method.
func (ep poller) Stop(desc *Desc) error {
	return ep.Del(desc.fd())
}

// Resume implements Poller.Resume() method.
func (ep poller) Resume(desc *Desc) error {
	return ep.Mod(desc.fd(), toEpollEvent(desc.event))
}

func toEpollEvent(event Event) (ep EpollEvent) {
	if event&EventRead != 0 {
		ep |= EPOLLIN | EPOLLRDHUP
	}
	if event&EventWrite != 0 {
		ep |= EPOLLOUT
	}
	if event&EventOneShot != 0 {
		ep |= EPOLLONESHOT
	}
	if event&EventEdgeTriggered != 0 {
		ep |= EPOLLET
	}
	return ep
}
