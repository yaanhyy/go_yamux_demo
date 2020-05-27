package main

import (
	"net"
	"context"
	ci "github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-secio"
	"github.com/libp2p/go-libp2p-core/peer"
		"github.com/libp2p/go-libp2p-core/sec"
	"io"
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

func NewTransport(typ, bits int) *secio.Transport {
	priv, pub, err := ci.GenerateKeyPair(typ, bits)
	if err != nil {
		print(err)
	}
	id, err := peer.IDFromPublicKey(pub)
	if err != nil {
		print(err)
	}
	return &secio.Transport{
		PrivateKey: priv,
		LocalID:    id,
	}
}

func ReadBuf( serverConn sec.SecureConn) {

	after := make([]byte, 100)
	_, err := io.ReadFull(serverConn, after)
	if err != nil {
		print(err)
	}
	fmt.Printf("%s", after)
}


func main()  {
	server := newServerConn()
	serverTpt := NewTransport(ci.Ed25519, 2048)
	serverConn, _ := serverTpt.SecureInbound(context.TODO(), server)
	ReadBuf(serverConn)
}
