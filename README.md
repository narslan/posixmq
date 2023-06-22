### posixmq

This package provides two things: 
A wrapper around posix message queue.
It allows monitoring the queues via epoll.

To listen a queue
```sh
go run cmd/posixmq-listener/main.go sample_queue 
```

Send a message to the queue 
```sh
go run cmd/posixmq-cli/main.go send sample_queue "a message"
```


