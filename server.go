package main

import (
	"container/list"
	"fmt"
	"log"
	"net"

)

type MessageFormat string 
const (
	TEXT MessageFormat = "TEXT"
	VIDEO MessageFormat = "VIDEO"
	AUDIO MessageFormat = "AUDIO"
)

type MessageType int8
const(
	clientConnected MessageType= iota
	clientDisconnected MessageType = iota
	newMessage MessageType= iota
)

type body struct{
	data []byte
	_type MessageFormat 
	read bool
	received bool
	sent bool
}

type client struct{
	conn net.Conn
	text body 
}

type message struct{
	// from client
	// to client
	// message body
	MessageType MessageType 
	Connection net.Conn
	Text []byte
}

//TODO: DO I need to make a list of the connections in order to ban them somehow in the future????
type server struct{
	active_conn *list.List
	pending_messages []message //maybe if the queue gets too long
	listener net.Listener
	err error
	address string

}

func (s *server) start_server(port string){

	s.address = port
	s.active_conn = list.New()
	s.listener, s.err = net.Listen("tcp", s.address)
	if s.err != nil {
		log.Println("Could not start the server at port", s.address)
		log.Println(s.err)
	}
	println("Started Listening at port ", s.address)
}

func (s *server) read_data(connection net.Conn,channel chan<- message){

	println("Reading Data")
	// temp_message := <-channel
	for{
		buffer := make([]byte, 1024)
		n , err := connection.Read(buffer)
		if err != nil {
			log.Println("Error while Reading the data")
			log.Println(err)
			defer connection.Close()
			return
		}
		dataStr := string(buffer[0:n])
		fmt.Print("Message Received:", dataStr)
		// connection.Write(buffer)

		channel <- message{
			MessageType: newMessage,
			Connection:connection,
			Text: buffer[0:n],
		}
	}
}

func (s *server) write_data(channel <-chan message){
	for{
		temp_message:=  <-channel	// 
		if(temp_message.MessageType == clientConnected){
			s.active_conn.PushBack(temp_message.Connection)
			println("Remote Address of client",temp_message.Connection.RemoteAddr().String())
			println("Local Address of client",temp_message.Connection.LocalAddr().String())
			println("Message Age: ",temp_message.MessageType)
			println("Added a new connection to the list of active connections")
		}else if(temp_message.MessageType == newMessage){
			println("Message Age: ",temp_message.MessageType)
			println("Message Text: ", string(temp_message.Text))
			temp_message.Connection.Write(temp_message.Text)
		}
	}
}

func main(){

	server := new(server)
	server.start_server(":8000")

	data_channel := make(chan message)

	go server.write_data(data_channel)

	for{
		conn, err := server.listener.Accept()
		if err != nil {
			log.Println(err)
			defer conn.Close()
			return
		}
		data_channel <- message{
			MessageType: clientConnected,
			Connection: conn,
		}
		go server.read_data(conn,data_channel)
	}
}
