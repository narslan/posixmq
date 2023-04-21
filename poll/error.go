// Copyright (c) 2017 Mail.Ru Group All rights reserved.
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

package poll

import (
	"fmt"
	"syscall"
)

var (
	// ErrNotFiler is returned by Handle* functions to indicate that given
	// net.Conn does not provide access to its file descriptor.
	ErrNotFiler = fmt.Errorf("could not get file descriptor")

	// ErrClosed is returned by Poller methods to indicate that instance is
	// closed and operation could not be processed.
	ErrClosed = fmt.Errorf("poller instance is closed")

	// ErrRegistered is returned by Poller Start() method to indicate that
	// connection with the same underlying file descriptor was already
	// registered within the poller instance.
	ErrRegistered = fmt.Errorf("file descriptor is already registered in poller instance")

	// ErrNotRegistered is returned by Poller Stop() and Resume() methods to
	// indicate that connection with the same underlying file descriptor was
	// not registered before within the poller instance.
	ErrNotRegistered = fmt.Errorf("file descriptor was not registered before in poller instance")
)

func temporaryErr(err error) bool {
	errno, ok := err.(syscall.Errno)
	if !ok {
		return false
	}
	return errno.Temporary()
}
