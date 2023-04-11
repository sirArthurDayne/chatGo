package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
)

type Client chan<- string // canal para transmitir mensajes

const (
	DEFAULT_PORT = 8080
	DEFAULT_HOST = "localhost"
)

var (
	incomingClients = make(chan Client)
	leavingClients  = make(chan Client)
	messagesBuffer  = make(chan string)
	host            = flag.String("h", DEFAULT_HOST, "Host to be connected.(default=localhost)")
	port            = flag.Int("p", DEFAULT_PORT, "Port to connect.(default=8080)")
)

func HandleConnection(conn net.Conn) {
	defer conn.Close()
	clientMessages := make(chan string)
	go MessageWrite(conn, clientMessages)
	clientName := conn.RemoteAddr().String() // 192.168.1.11:8080
	// send greeting message to new client only
	clientMessages <- fmt.Sprintf("Welcome to the server, %s\n", clientName)

	// send messages to all clients
	messagesBuffer <- clientName + "has joined!!"

	// add all new clients
	incomingClients <- clientMessages

	// read all messages from the clients as long as they are connected.
	inputMessage := bufio.NewScanner(conn)
	for inputMessage.Scan() {
		messagesBuffer <- clientName + ": " + inputMessage.Text()
	}

	// remve clients after they disconnect
	leavingClients <- clientMessages
	messagesBuffer <- clientName + " has left the chat!"
}

func MessageWrite(conn net.Conn, messagesBuffer <-chan string) {
	for message := range messagesBuffer {
		fmt.Fprintln(conn, message)
	}
}

func Broadcast() {
	ClientConnState := make(map[Client]bool)
	// multiplexing of messages
	for {
		// all events in chat
		select {
		case msg := <-messagesBuffer: // when arrrive a new message from any clients
			// send message to all clients
			for currentClient := range ClientConnState {
				currentClient <- msg
			}
		case currentClient := <-incomingClients: // When new client connects
			ClientConnState[currentClient] = true
		case currentClient := <-leavingClients: // When client disconnect
			delete(ClientConnState, currentClient)
			close(currentClient)
		}
	}
}

func main() {
	flag.Parse()
	// create a new server and listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *host, *port))
	if err != nil {
		log.Fatal("[ERROR]:" + err.Error())
	}
	go Broadcast()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("[ERROR]:" + err.Error())
			continue
		}
		go HandleConnection(conn)
	}
}
