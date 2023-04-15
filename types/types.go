package types

type (
	Client  chan<- string // canal para transmitir usuarios
	Message chan string
)
