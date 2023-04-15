package handler

import (
	"bufio"
	"fmt"
	"net"

	"github.com/sirArthurDayne/chatGo/message"
	"github.com/sirArthurDayne/chatGo/types"
)

type Handler struct {
	msgService      *message.Service
	messagesBuffer  types.Message
	incomingClients chan types.Client
	leavingClients  chan types.Client
}

func NewHandler(service *message.Service, msg types.Message, incoming chan types.Client, leaving chan types.Client) *Handler {
	return &Handler{
		msgService:      service,
		messagesBuffer:  msg,
		incomingClients: incoming,
		leavingClients:  leaving,
	}
}

func (h *Handler) HandleConnection(conn net.Conn) {
	defer conn.Close()
	clientMessages := make(chan string)
	go h.msgService.MessageWriter(conn, clientMessages)
	clientName := conn.RemoteAddr().String() // 192.168.1.11:8080
	// send greeting message to new client only
	clientMessages <- fmt.Sprintf("Welcome to the server, %s\n", clientName)

	// send messages to all clients
	h.messagesBuffer <- clientName + "has joined!!"

	// add all new clients
	h.incomingClients <- clientMessages

	// read all messages from the clients as long as they are connected.
	inputMessage := bufio.NewScanner(conn)
	for inputMessage.Scan() {
		h.messagesBuffer <- clientName + ": " + inputMessage.Text()
	}

	// remve clients after they disconnect
	h.leavingClients <- clientMessages
	h.messagesBuffer <- clientName + " has left the chat!"
}
