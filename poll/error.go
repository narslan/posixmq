package poll

import "syscall"

func temporaryErr(err error) bool {
	errno, ok := err.(syscall.Errno)
	if !ok {
		return false
	}
	return errno.Temporary()
}

func setNonblock(fd int, nonblocking bool) (err error) {
	return syscall.SetNonblock(fd, nonblocking)
}
