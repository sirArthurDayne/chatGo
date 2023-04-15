package message

import (
	"fmt"
	"net"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (Service) MessageWriter(conn net.Conn, messagesBuffer <-chan string) {
	for message := range messagesBuffer {
		fmt.Fprintln(conn, message)
	}
}
