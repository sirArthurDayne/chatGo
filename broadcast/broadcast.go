package broadcast

import "github.com/sirArthurDayne/chatGo/types"

type BroadcastService struct {
	messages        types.Message
	incomingClients chan types.Client
	leavingClients  chan types.Client
}

func NewBroadcastService(messages types.Message, incoming chan types.Client, leaving chan types.Client) *BroadcastService {
	return &BroadcastService{
		messages:        messages,
		incomingClients: incoming,
		leavingClients:  leaving,
	}
}

func (s *BroadcastService) Broadcast() {
	ClientConnState := make(map[types.Client]bool)
	// multiplexing of messages
	for {
		// all events in chat
		select {
		case msg := <-s.messages: // when arrrive a new message from any clients
			// send message to all clients
			for currentClient := range ClientConnState {
				currentClient <- msg
			}
		case currentClient := <-s.incomingClients: // When new client connects
			ClientConnState[currentClient] = true
		case currentClient := <-s.leavingClients: // When client disconnect
			delete(ClientConnState, currentClient)
			close(currentClient)
		}
	}
}
