package main

import (
	"net"
	"github.com/libp2p/go-libp2p-core/crypto"
	"context"
	"github.com/libp2p/go-libp2p-noise"
	"io"
	"github.com/libp2p/go-libp2p-core/sec"
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

func NewNoiseTransport(typ, bits int) *noise.Transport {
	priv, _, err := crypto.GenerateKeyPair(typ, bits)
	if err != nil {
		print(err)
	}
	//id, err := peer.IDFromPublicKey(pub)
	//if err != nil {
	//	print(err)
	//}
	transport,err := noise.New(priv)
	return  transport

}

func ReadBuf( serverConn sec.SecureConn) {

	after := make([]byte, 100)
	_, err := io.ReadFull(serverConn, after)
	if err != nil {
		print(err.Error())
	}
	fmt.Printf("%s", after)
}

func main()  {
	server := newServerConn()
	serverTpt := NewNoiseTransport(crypto.Ed25519, 2048)
	serverConn, _ := serverTpt.SecureInbound(context.TODO(), server)
	ReadBuf(serverConn)
}