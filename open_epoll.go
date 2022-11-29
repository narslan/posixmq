package posixmq

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"time"

	"github.com/go-kit/log"
	"github.com/mailru/easygo/netpoll"
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

	handler := func(evt netpoll.EpollEvent) {

		// @nevroz I couln't import _EPOLLCLOSED from netpoll package
		if evt&0x20 != 0 {
			return
		}
		logger.Log("call", "first handler")
		logger.Log("event", evt)

	}
	ep.Add(m.fd, netpoll.EPOLLIN, handler)

	data := []byte("hello there!")
	m.Send(ctx, data, 10)
	time.Sleep(1 * time.Millisecond)

	logger.Log("call", "after add")
	if err = ep.Close(); err != nil {

		panic(err)
	}

	logger.Log("received", string(received.Bytes()))
	return
}
