package main

import (
	"context"
ci "github.com/libp2p/go-libp2p-core/crypto"
	"net"
	"github.com/libp2p/go-libp2p-secio"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/sec"

	"fmt"
)


func newClientConn()(net.Conn) {

	client, clientErr := net.Dial("tcp", "localhost:5679")

	if clientErr != nil {
		print("Failed to accept:", clientErr)
	}



	return client
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


func WriteBuf( clientConn sec.SecureConn) {
	before := []byte("hello world")
	_, err := clientConn.Write(before)
	if err != nil {
		print(err)
	}

	fmt.Printf("%s", before)
}

func main()  {
	client := newClientConn()
	clientTpt := NewTransport(ci.Secp256k1, 2048)
	clientConn, err :=clientTpt.SecureOutbound(context.TODO(), client, "")
	if err == nil {
		WriteBuf(clientConn)
	} else {
		fmt.Printf("%+v", err)
	}
}
