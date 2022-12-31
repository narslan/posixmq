package posixmq

import (
	"context"
	"github.com/go-kit/log"

	"os"
)

func OpenWithEpoll(ctx context.Context, cfg *Config) (m *MessageQueue, err error) {

	w := log.NewSyncWriter(os.Stdout)
	logger := log.NewLogfmtLogger(w)
	logger = log.With(logger, "caller", log.DefaultCaller)

	m, err = Open(ctx, cfg)

	if err != nil {
		panic(err)
	}
	defer m.Close(ctx)

	e, err := MkEpoll()
	if err != nil {
		panic(err)
	}
	err = e.Add(m.fd)
	if err != nil {
		panic(err)
	}

	desc, err := e.Wait()

	if err != nil {
		panic(err)
	}

	for key, value := range desc {
		logger.Log("idx", key, "descriptors", value)
	}

	return
}
