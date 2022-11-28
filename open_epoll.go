package posixmq

import (
	"bytes"
	"context"
	"fmt"
	"os"

	"github.com/go-kit/log"
	"github.com/mailru/easygo/netpoll"
	"golang.org/x/sys/unix"
)

func epollConfig() *netpoll.EpollConfig {
	return &netpoll.EpollConfig{
		OnWaitError: func(err error) {
			fmt.Printf("error on wait: %s", err.Error())
		},
	}
}

func OpenWithEpoll(ctx context.Context, cfg *Config) (m *MessageQueue, err error) {

	w := log.NewSyncWriter(os.Stdout)
	logger := log.NewLogfmtLogger(w)
	logger = log.With(logger, "caller", log.DefaultCaller)

	ep, err := netpoll.EpollCreate(epollConfig())
	if err != nil {
		panic(err)
	}

	m, err = Open(ctx, cfg)

	if err != nil {
		panic(err)
	}
	defer m.Close(ctx)

	var received bytes.Buffer
	done := make(chan struct{})

	handler := func(evt netpoll.EpollEvent) {

		logger.Log("call", "first handler")
		// @nevroz I couln't import _EPOLLCLOSED from netpoll package
		if evt&0x20 != 0 {
			return
		}

		ep.Add(m.fd, netpoll.EPOLLIN|netpoll.EPOLLET|netpoll.EPOLLHUP|netpoll.EPOLLRDHUP, func(evt netpoll.EpollEvent) {
			logger.Log("call", "second handler")
			// If EPOLLRDHUP is supported, it will be triggered after conn
			// close() or shutdown(). In older versions EPOLLHUP is triggered.
			if evt&0x20 != 0 {
				return
			}

			var buf [100]byte
			for {

				n, _ := unix.Read(m.fd, buf[:])
				if n == 0 {
					close(done)
				}
				if n <= 0 {
					break
				}
				received.Write(buf[:n])
			}
		})

	}
	ep.Add(m.fd, netpoll.EPOLLIN, handler)
	<-done

	if err = ep.Close(); err != nil {

		panic(err)
	}

	logger.Log("received", string(received.Bytes()))
	return
}
