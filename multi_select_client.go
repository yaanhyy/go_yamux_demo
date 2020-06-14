package main

import (
	"net"
	mss "github.com/multiformats/go-multistream"
	"fmt"
)

func newClientConn()(net.Conn) {

	client, clientErr := net.Dial("tcp", "localhost:5679")

	if clientErr != nil {
		print("Failed to accept:", clientErr)
	}



	return client
}

func main() {

	client := newClientConn()

	err := mss.SelectProtoOrFail("/b", client)
	if err == nil {
		print("/a");
	} else {
		fmt.Printf("%+v",err);
	}
}
