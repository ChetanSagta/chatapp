package main

import (
	"bufio"
	"fmt"
	"os"
	"net"
)

var SERVER = "127.0.0.1"
var PORT = ":8000"

func start_client() net.Conn {
	conn, err := net.Dial("tcp", SERVER+PORT)
	if err != nil {
		println("Error Creating the client", err.Error())
		return nil
	}
	return conn
}

func intro() {
	println("Hello.. Welcome to the \"CHATAPP\" client")
}

func get_input() []byte {
	println("Enter the text you want to send:")
	input := make([]byte,1024)
	reader := bufio.NewReader(os.Stdin)
	input , err := reader.ReadBytes('\n')
	if err != nil {
		println("Error while reading input from the keyboard: ", err.Error())
		return nil 
	}

	return input
}

type client struct {
	connection net.Conn
}

func main() {

	client := client{}
	client.connection = start_client()
	if client.connection == nil {
		return
	}
	intro()
	for {
		if client.connection == nil {
			fmt.Println("Client connection closed, exiting")
			return
		}
		text := get_input()
		if text != nil{
			go client.connection.Write(text)
			_, err :=client.connection.Read(text)
			if(err != nil){
				fmt.Println("Server Socket closed. Closing client connection", err.Error())
				client.connection.Close()
				break
			}
		}
	}
}
