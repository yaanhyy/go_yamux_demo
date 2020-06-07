package main

import (
	"net"
	mss "github.com/multiformats/go-multistream"
	"fmt"
)

func newServerConn()(net.Conn) {
	lstnr, err := net.Listen("tcp", "localhost:8981")
	if err != nil {

		return nil
	}

	server, err := lstnr.Accept()

	lstnr.Close()

	if err != nil {
		print("Failed to accept:", err)
	}



	return server
}

func main() {
	// 建立底层复用连接
	server := newServerConn()
	mux := mss.NewMultistreamMuxer()
	mux.AddHandler("/a", nil)
	selected, _, err := mux.Negotiate(server)
	if err == nil {
		print(selected);
	} else {
		fmt.Printf("%+v",err);
	}
}
