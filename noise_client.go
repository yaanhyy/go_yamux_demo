package main

import (
	"github.com/libp2p/go-libp2p-noise"
	"github.com/libp2p/go-libp2p-core/crypto"
	"net"
	"context"
	"fmt"
)

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

func newClientConn()(net.Conn) {

	client, clientErr := net.Dial("tcp", "localhost:8981")

	if clientErr != nil {
		print("Failed to accept:", clientErr)
	}



	return client
}

func main() {
	client := newClientConn()
	initTransport := NewNoiseTransport(crypto.Ed25519, 2048)
	secureConn, err := initTransport.SecureOutbound(context.TODO(), client, "12D3KooWNWJ98NujpQb4WwxNmKoqYVuAa5Wo8iZz8X99ArBwcRrY")
	if err == nil {
		before := []byte("hello world")
		secureConn.Write(before)
	} else {
		fmt.Printf("%+v", err)
	}
}
