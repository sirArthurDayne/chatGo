package server

import (
	"fmt"
	"log"
	"net"

	"github.com/sirArthurDayne/chatGo/broadcast"
	"github.com/sirArthurDayne/chatGo/handler"
	"github.com/sirArthurDayne/chatGo/message"
	"github.com/sirArthurDayne/chatGo/types"
)

type Server struct {
	msgService       *message.Service
	broadcastService *broadcast.BroadcastService
	handlers         *handler.Handler
	incomingClients  chan types.Client
	leavingClients   chan types.Client
	messagesBuffer   types.Message
}

func NewServer() *Server {
	return &Server{
		incomingClients: make(chan types.Client),
		leavingClients:  make(chan types.Client),
		messagesBuffer:  make(chan string),
	}
}

func (s *Server) LoadServerComponents() {
	msgService := message.NewService()
	handler := handler.NewHandler(msgService, s.messagesBuffer, s.incomingClients, s.leavingClients)
	broadcastService := broadcast.NewBroadcastService(s.messagesBuffer, s.incomingClients, s.leavingClients)
	s.msgService = msgService
	s.broadcastService = broadcastService
	s.handlers = handler
}

func (s *Server) Start(host *string, port *int) {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *host, *port))
	if err != nil {
		log.Fatal("[ERROR]:" + err.Error())
	}
	go s.broadcastService.Broadcast()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("[ERROR]:" + err.Error())
			continue
		}
		go s.handlers.HandleConnection(conn)
	}
}
